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
	UserName      string
	PostID        uint64
	VoteType      string
	CommentID     uint64
	CommunityID   uint64
	CommunityName string
}

type NotificationService struct {
	notificationRepo        repository.NotificationRepository
	notificationSettingRepo repository.NotificationSettingRepository
	botTaskRepo             repository.BotTaskRepository
	userRepo                repository.UserRepository
	sseService              *SSEService
	botTaskService          *BotTaskService
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	notificationSettingRepo repository.NotificationSettingRepository,
	botTaskRepo repository.BotTaskRepository,
	userRepo repository.UserRepository,
	sseService *SSEService,
	botTaskService *BotTaskService,
) *NotificationService {
	return &NotificationService{
		notificationRepo:        notificationRepo,
		notificationSettingRepo: notificationSettingRepo,
		botTaskRepo:             botTaskRepo,
		userRepo:                userRepo,
		sseService:              sseService,
		botTaskService:          botTaskService,
	}
}

func (s *NotificationService) CreateNotification(userID uint64, action string, notifPayload interface{}) error {
	// Check notification settings for this action
	setting, err := s.notificationSettingRepo.GetUserNotificationSetting(userID, action)
	if err != nil {
		// If no setting found, use default settings
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
	case constant.NOTIFICATION_ACTION_POST_APPROVED:
		return basePath + "post_approved.txt"
	case constant.NOTIFICATION_ACTION_POST_REJECTED:
		return basePath + "post_rejected.txt"
	case constant.NOTIFICATION_ACTION_POST_DELETED:
		return basePath + "post_deleted.txt"
	case constant.NOTIFICATION_ACTION_POST_REPORTED:
		return basePath + "post_reported.txt"
	case constant.NOTIFICATION_ACTION_SUBSCRIPTION_APPROVED:
		return basePath + "subscription_approved.txt"
	case constant.NOTIFICATION_ACTION_SUBSCRIPTION_REJECTED:
		return basePath + "subscription_rejected.txt"
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
	case constant.NOTIFICATION_ACTION_POST_APPROVED:
		return basePath + "post_approved_email.html"
	case constant.NOTIFICATION_ACTION_POST_REJECTED:
		return basePath + "post_rejected_email.html"
	case constant.NOTIFICATION_ACTION_POST_DELETED:
		return basePath + "post_deleted_email.html"
	case constant.NOTIFICATION_ACTION_POST_REPORTED:
		return basePath + "post_reported_email.html"
	case constant.NOTIFICATION_ACTION_SUBSCRIPTION_APPROVED:
		return basePath + "subscription_approved_email.html"
	case constant.NOTIFICATION_ACTION_SUBSCRIPTION_REJECTED:
		return basePath + "subscription_rejected_email.html"
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
	case constant.NOTIFICATION_ACTION_POST_REPORTED:
		if p, ok := notifPayload.(payload.PostReportNotificationPayload); ok {
			data.UserName = p.UserName
			data.PostID = p.PostID
		}
	case constant.NOTIFICATION_ACTION_POST_APPROVED,
		constant.NOTIFICATION_ACTION_POST_REJECTED,
		constant.NOTIFICATION_ACTION_POST_DELETED:
		if p, ok := notifPayload.(map[string]interface{}); ok {
			if postID, ok := p["postId"].(uint64); ok {
				data.PostID = postID
			}
		}
	case constant.NOTIFICATION_ACTION_SUBSCRIPTION_APPROVED,
		constant.NOTIFICATION_ACTION_SUBSCRIPTION_REJECTED:
		if p, ok := notifPayload.(payload.SubscriptionNotificationPayload); ok {
			data.CommunityID = p.CommunityID
			data.CommunityName = p.CommunityName
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
	if s.botTaskService != nil {
		fullSubject := fmt.Sprintf("Notification: %s", subject)
		if err := s.botTaskService.CreateEmailTask(email, fullSubject, body); err != nil {
			log.Printf("[Err] Failed to create email bot task: %v", err)
		}
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

func (s *NotificationService) GetUserNotificationSettings(userID uint64) ([]*response.NotificationSettingResponse, error) {
	settings, err := s.notificationSettingRepo.GetUserNotificationSettings(userID)
	if err != nil {
		log.Printf("[Err] Failed to get user notification settings: %v", err)
		return nil, fmt.Errorf("failed to get notification settings")
	}

	settingResponses := make([]*response.NotificationSettingResponse, len(settings))
	for i, setting := range settings {
		settingResponses[i] = response.NewNotificationSettingResponse(setting)
	}

	return settingResponses, nil
}

func (s *NotificationService) UpdateNotificationSetting(userID uint64, action string, isPush, isSendMail *bool) error {
	// Get existing setting
	setting, err := s.notificationSettingRepo.GetUserNotificationSetting(userID, action)
	if err != nil {
		// If setting doesn't exist, create a new one with default values
		setting = &model.NotificationSetting{
			UserID:     userID,
			Action:     action,
			IsPush:     true,
			IsSendMail: false,
		}
	}

	if isPush != nil {
		setting.IsPush = *isPush
	}
	if isSendMail != nil {
		setting.IsSendMail = *isSendMail
	}

	// Upsert the setting
	if err := s.notificationSettingRepo.UpsertNotificationSetting(setting); err != nil {
		log.Printf("[Err] Failed to update notification setting: %v", err)
		return fmt.Errorf("failed to update notification setting")
	}

	return nil
}

// func (s *NotificationService) CreateDefaultNotificationSettings(userID uint64) error {
// 	actions := []string{
// 		constant.NOTIFICATION_ACTION_GET_POST_VOTE,
// 		constant.NOTIFICATION_ACTION_GET_POST_NEW_COMMENT,
// 		constant.NOTIFICATION_ACTION_GET_COMMENT_VOTE,
// 		constant.NOTIFICATION_ACTION_GET_COMMENT_REPLY,
// 	}

// 	now := time.Now()
// 	settings := make([]*model.NotificationSetting, len(actions))
// 	for i, action := range actions {
// 		settings[i] = &model.NotificationSetting{
// 			UserID:     userID,
// 			Action:     action,
// 			IsPush:     true,
// 			IsSendMail: false,
// 			CreatedAt:  now,
// 		}
// 	}

// 	if err := s.notificationSettingRepo.CreateNotificationSettings(settings); err != nil {
// 		log.Printf("[Err] Failed to create default notification settings: %v", err)
// 		return fmt.Errorf("failed to create default notification settings")
// 	}

// 	return nil
// }
