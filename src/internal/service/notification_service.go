package service

import (
	"encoding/json"
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/template/payload"
	"social-platform-backend/package/util"
	"time"
)

type NotificationTemplateData struct {
	UserName  string
	PostID    uint64
	VoteType  string
	CommentID uint64
}

type NotificationService struct {
	notificationRepo        repository.NotificationRepository
	notificationSettingRepo repository.NotificationSettingRepository
	botTaskRepo             repository.BotTaskRepository
	userRepo                repository.UserRepository
	sseService              *SSEService
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	notificationSettingRepo repository.NotificationSettingRepository,
	botTaskRepo repository.BotTaskRepository,
	userRepo repository.UserRepository,
	sseService *SSEService,
) *NotificationService {
	return &NotificationService{
		notificationRepo:        notificationRepo,
		notificationSettingRepo: notificationSettingRepo,
		botTaskRepo:             botTaskRepo,
		userRepo:                userRepo,
		sseService:              sseService,
	}
}

func (s *NotificationService) CreateNotification(userID uint64, action string, notifPayload interface{}) error {
	// Check notification settings for this action
	setting, err := s.notificationSettingRepo.GetUserNotificationSetting(userID, action)
	if err != nil {
		// If no setting found, create default settings (both enabled)
		log.Printf("[Info] No notification setting found for user %d and action %s, using defaults", userID, action)
		setting = &model.NotificationSetting{
			UserID:     userID,
			Action:     action,
			IsPush:     true,
			IsSendMail: true,
		}
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("[Err] Failed to get user for notification: %v", err)
		return fmt.Errorf("failed to get user")
	}

	templateData := s.prepareTemplateData(action, notifPayload)

	if setting.IsPush {
		templatePath := s.getNotificationTemplatePath(action)
		body, err := util.RenderTemplate(templatePath, templateData)
		if err != nil {
			log.Printf("[Err] Failed to render notification body: %v", err)
			return err
		}

		payloadBytes, err := json.Marshal(notifPayload)
		if err != nil {
			log.Printf("[Err] Failed to marshal notification payload: %v", err)
			return err
		}

		rawPayload := json.RawMessage(payloadBytes)
		notification := &model.Notification{
			UserID:    userID,
			Body:      body,
			Action:    action,
			Payload:   &rawPayload,
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		if err := s.notificationRepo.CreateNotification(notification); err != nil {
			log.Printf("[Err] Failed to create notification: %v", err)
			return err
		}

		go s.broadcastNewNotification(userID, notification)
	}

	if setting.IsSendMail {
		emailTemplatePath := s.getNotificationEmailTemplatePath(action)
		emailBody, err := util.RenderTemplate(emailTemplatePath, templateData)
		if err != nil {
			log.Printf("[Err] Failed to render notification email body: %v", err)
		} else {
			go s.sendEmailNotification(user.Email, action, emailBody)
		}
	}

	return nil
}

func (s *NotificationService) getNotificationTemplatePath(action string) string {
	basePath := "package/template/notification/"
	switch action {
	case constant.NOTIFICATION_ACTION_GET_POST_VOTE:
		return basePath + "post_vote.txt"
	case constant.NOTIFICATION_ACTION_GET_POST_NEW_COMMENT:
		return basePath + "post_comment.txt"
	case constant.NOTIFICATION_ACTION_GET_COMMENT_VOTE:
		return basePath + "comment_vote.txt"
	case constant.NOTIFICATION_ACTION_GET_COMMENT_REPLY:
		return basePath + "comment_reply.txt"
	default:
		return ""
	}
}

func (s *NotificationService) getNotificationEmailTemplatePath(action string) string {
	basePath := "package/template/notification/"
	switch action {
	case constant.NOTIFICATION_ACTION_GET_POST_VOTE:
		return basePath + "post_vote_email.html"
	case constant.NOTIFICATION_ACTION_GET_POST_NEW_COMMENT:
		return basePath + "post_comment_email.html"
	case constant.NOTIFICATION_ACTION_GET_COMMENT_VOTE:
		return basePath + "comment_vote_email.html"
	case constant.NOTIFICATION_ACTION_GET_COMMENT_REPLY:
		return basePath + "comment_reply_email.html"
	default:
		return ""
	}
}

func (s *NotificationService) prepareTemplateData(action string, notifPayload interface{}) NotificationTemplateData {
	data := NotificationTemplateData{}

	switch action {
	case constant.NOTIFICATION_ACTION_GET_POST_VOTE:
		if p, ok := notifPayload.(payload.PostVoteNotificationPayload); ok {
			data.UserName = p.UserName
			data.PostID = p.PostID
			if p.VoteType {
				data.VoteType = "upvoted"
			} else {
				data.VoteType = "downvoted"
			}
		}
	case constant.NOTIFICATION_ACTION_GET_POST_NEW_COMMENT:
		if p, ok := notifPayload.(payload.PostCommentNotificationPayload); ok {
			data.UserName = p.UserName
			data.PostID = p.PostID
		}
	case constant.NOTIFICATION_ACTION_GET_COMMENT_VOTE:
		if p, ok := notifPayload.(payload.CommentVoteNotificationPayload); ok {
			data.UserName = p.UserName
			data.CommentID = p.CommentID
			if p.VoteType {
				data.VoteType = "upvoted"
			} else {
				data.VoteType = "downvoted"
			}
		}
	case constant.NOTIFICATION_ACTION_GET_COMMENT_REPLY:
		if p, ok := notifPayload.(payload.CommentReplyNotificationPayload); ok {
			data.UserName = p.UserName
			data.CommentID = p.CommentID
		}
	}

	return data
}

func (s *NotificationService) broadcastNewNotification(userID uint64, notification *model.Notification) {
	unreadCount, _ := s.notificationRepo.GetUnreadCount(userID)

	event := &response.SSEEvent{
		Event: "new_notification",
		Data: response.NewNotificationEvent{
			Notification: *response.NewNotificationResponse(notification),
			UnreadCount:  unreadCount,
		},
	}
	s.sseService.BroadcastToUser(userID, event)
}

func (s *NotificationService) sendEmailNotification(email, subject, body string) {
	emailPayload := payload.EmailPayload{
		To:      email,
		Subject: fmt.Sprintf("Notification: %s", subject),
		Body:    body,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		log.Printf("[Err] Failed to marshal email payload: %v", err)
		return
	}

	rawPayload := json.RawMessage(payloadBytes)
	now := time.Now()
	botTask := &model.BotTask{
		Action:     constant.BOT_TASK_ACTION_SEND_EMAIL,
		Payload:    &rawPayload,
		CreatedAt:  now,
		ExecutedAt: &now,
	}

	if err := s.botTaskRepo.CreateBotTask(botTask); err != nil {
		log.Printf("[Err] Failed to create email bot task: %v", err)
	}
}

func (s *NotificationService) GetUserNotifications(userID uint64, page, limit int) ([]*response.NotificationResponse, *response.Pagination, error) {
	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	offset := (page - 1) * limit

	notifications, total, err := s.notificationRepo.GetUserNotifications(userID, limit, offset)
	if err != nil {
		log.Printf("[Err] Failed to get user notifications: %v", err)
		return nil, nil, fmt.Errorf("failed to get notifications")
	}

	notifResponses := make([]*response.NotificationResponse, len(notifications))
	for i, notif := range notifications {
		notifResponses[i] = response.NewNotificationResponse(notif)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/notifications?page=%d&limit=%d", page+1, limit)
	}

	return notifResponses, pagination, nil
}

func (s *NotificationService) MarkAsRead(userID, notificationID uint64) error {
	notification, err := s.notificationRepo.GetNotificationByID(notificationID)
	if err != nil {
		log.Printf("[Err] Notification not found: %v", err)
		return fmt.Errorf("notification not found")
	}

	if notification.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *NotificationService) MarkAllAsRead(userID uint64) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *NotificationService) DeleteNotification(userID, notificationID uint64) error {
	notification, err := s.notificationRepo.GetNotificationByID(notificationID)
	if err != nil {
		log.Printf("[Err] Notification not found: %v", err)
		return fmt.Errorf("notification not found")
	}

	if notification.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return s.notificationRepo.DeleteNotification(notificationID)
}

func (s *NotificationService) GetUnreadCount(userID uint64) (int64, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}
