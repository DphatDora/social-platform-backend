package service

import (
	"context"
	"fmt"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/logger"
	"time"
)

type MessageService struct {
	conversationRepo      repository.ConversationRepository
	messageRepo           repository.MessageRepository
	messageAttachmentRepo repository.MessageAttachmentRepository
	userRepo              repository.UserRepository
	sseService            *SSEService
}

func NewMessageService(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	messageAttachmentRepo repository.MessageAttachmentRepository,
	userRepo repository.UserRepository,
	sseService *SSEService,
) *MessageService {
	return &MessageService{
		conversationRepo:      conversationRepo,
		messageRepo:           messageRepo,
		messageAttachmentRepo: messageAttachmentRepo,
		userRepo:              userRepo,
		sseService:            sseService,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, senderID uint64, req *request.SendMessageRequest) (*response.MessageResponse, error) {
	_, err := s.userRepo.GetUserByID(req.RecipientID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Recipient not found: %v", err)
		return nil, fmt.Errorf("recipient not found")
	}

	conversation, err := s.conversationRepo.CreateOrGetConversation(senderID, req.RecipientID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error creating/getting conversation: %v", err)
		return nil, fmt.Errorf("failed to create conversation")
	}

	var metaData *model.MetaData
	if req.MetaData != nil {
		metaData = &model.MetaData{
			ID:        req.MetaData.ID,
			Title:     req.MetaData.Title,
			Tags:      req.MetaData.Tags,
			Content:   req.MetaData.Content,
			MediaURLs: req.MetaData.MediaURLs,
		}
	}

	message := &model.Message{
		ConversationID: conversation.ID,
		SenderID:       senderID,
		Content:        req.Content,
		IsRead:         false,
		MetaData:       metaData,
	}

	if err := s.messageRepo.CreateMessage(message); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error creating message: %v", err)
		return nil, fmt.Errorf("failed to send message")
	}

	// Save attachments if provided
	if len(req.Attachments) > 0 {
		attachments := make([]model.MessageAttachment, len(req.Attachments))
		for i, attachment := range req.Attachments {
			attachments[i] = model.MessageAttachment{
				MessageID: message.ID,
				FileURL:   attachment.FileURL,
				FileType:  attachment.FileType,
				CreatedAt: time.Now(),
			}
		}
		if err := s.messageAttachmentRepo.CreateMessageAttachments(attachments); err != nil {
			logger.WarnfWithCtx(ctx, "[Warn] Failed to save message attachments: %v", err)
		}
	}

	if err := s.conversationRepo.UpdateLastMessage(conversation.ID, message.ID); err != nil {
		logger.WarnfWithCtx(ctx, "[Warn] Failed to update last message: %v", err)
	}

	fullMessage, err := s.messageRepo.GetMessageByID(message.ID)
	if err != nil {
		logger.WarnfWithCtx(ctx, "[Warn] Failed to get full message: %v", err)
		fullMessage = message
	}

	messageResp := response.NewMessageResponse(fullMessage)

	// Broadcast new message to recipient via SSE
	go s.broadcastNewMessage(ctx, req.RecipientID, conversation.ID, messageResp)

	// Broadcast conversation update to both users
	go s.broadcastConversationUpdate(ctx, senderID, conversation.ID)
	go s.broadcastConversationUpdate(ctx, req.RecipientID, conversation.ID)

	return messageResp, nil
}

// sends new message event via SSE
func (s *MessageService) broadcastNewMessage(ctx context.Context, userID uint64, conversationID uint64, message *response.MessageResponse) {
	event := &response.SSEEvent{
		Event: "new_message",
		Data: response.NewMessageEvent{
			ConversationID: conversationID,
			Message:        *message,
		},
	}
	s.sseService.BroadcastToUser(ctx, userID, event)
}

// sends conversation update event via SSE
func (s *MessageService) broadcastConversationUpdate(ctx context.Context, userID uint64, conversationID uint64) {
	conversation, err := s.conversationRepo.GetConversationByID(conversationID)
	if err != nil {
		logger.WarnfWithCtx(ctx, "[Warn] Failed to get conversation for broadcast: %v", err)
		return
	}

	unreadCount, _ := s.messageRepo.GetUnreadCount(conversationID, userID)
	conversationResp := response.NewConversationListResponse(conversation, userID, unreadCount)

	event := &response.SSEEvent{
		Event: "conversation_updated",
		Data: response.ConversationUpdatedEvent{
			Conversation: *conversationResp,
		},
	}
	s.sseService.BroadcastToUser(ctx, userID, event)
}

func (s *MessageService) MarkMessageAsRead(ctx context.Context, userID, messageID uint64) error {
	message, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Message not found: %v", err)
		return fmt.Errorf("message not found")
	}

	isInConversation, err := s.conversationRepo.CheckUserInConversation(message.ConversationID, userID)
	if err != nil || !isInConversation {
		logger.ErrorfWithCtx(ctx, "[Err] User not in conversation")
		return fmt.Errorf("unauthorized")
	}

	if err := s.messageRepo.MarkMessageAsRead(messageID, userID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error marking message as read: %v", err)
		return fmt.Errorf("failed to mark message as read")
	}

	// Broadcast read status to sender
	go s.broadcastConversationUpdate(ctx, message.SenderID, message.ConversationID)
	go s.broadcastConversationUpdate(ctx, userID, message.ConversationID)

	return nil
}

func (s *MessageService) GetConversations(ctx context.Context, userID uint64, page, limit int) ([]*response.ConversationListResponse, *response.Pagination, error) {
	conversations, total, err := s.conversationRepo.GetUserConversations(userID, page, limit)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting conversations: %v", err)
		return nil, nil, fmt.Errorf("failed to get conversations")
	}

	conversationResponses := make([]*response.ConversationListResponse, len(conversations))
	for i, conv := range conversations {
		unreadCount, _ := s.messageRepo.GetUnreadCount(conv.ID, userID)
		conversationResponses[i] = response.NewConversationListResponse(conv, userID, unreadCount)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/messages/conversations?page=%d&limit=%d", page+1, limit)
	}

	return conversationResponses, pagination, nil
}

func (s *MessageService) GetMessages(ctx context.Context, userID, conversationID uint64, page, limit int) ([]*response.MessageResponse, *response.Pagination, error) {
	isInConversation, err := s.conversationRepo.CheckUserInConversation(conversationID, userID)
	if err != nil || !isInConversation {
		logger.ErrorfWithCtx(ctx, "[Err] User not in conversation")
		return nil, nil, fmt.Errorf("unauthorized")
	}

	messages, total, err := s.messageRepo.GetConversationMessages(conversationID, page, limit)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting messages: %v", err)
		return nil, nil, fmt.Errorf("failed to get messages")
	}

	messageResponses := make([]*response.MessageResponse, len(messages))
	for i, msg := range messages {
		messageResponses[i] = response.NewMessageResponse(msg)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/messages/conversations/%d/messages?page=%d&limit=%d", conversationID, page+1, limit)
	}

	return messageResponses, pagination, nil
}

func (s *MessageService) MarkConversationAsRead(ctx context.Context, userID, conversationID uint64) error {
	isInConversation, err := s.conversationRepo.CheckUserInConversation(conversationID, userID)
	if err != nil || !isInConversation {
		logger.ErrorfWithCtx(ctx, "[Err] User not in conversation")
		return fmt.Errorf("unauthorized")
	}

	if err := s.messageRepo.MarkConversationMessagesAsRead(conversationID, userID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error marking conversation as read: %v", err)
		return fmt.Errorf("failed to mark conversation as read")
	}

	// Broadcast update to both users
	conversation, _ := s.conversationRepo.GetConversationByID(conversationID)
	if conversation != nil {
		go s.broadcastConversationUpdate(ctx, conversation.User1ID, conversationID)
		go s.broadcastConversationUpdate(ctx, conversation.User2ID, conversationID)
	}

	return nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, userID, messageID uint64) error {
	message, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Message not found: %v", err)
		return fmt.Errorf("message not found")
	}

	// Check if user is the sender
	if message.SenderID != userID {
		logger.ErrorfWithCtx(ctx, "[Err] User is not the sender of this message")
		return fmt.Errorf("unauthorized: only sender can delete message")
	}

	// Check if message is within 10 minutes
	timeSinceCreation := time.Since(message.CreatedAt)
	if timeSinceCreation > 10*time.Minute {
		logger.ErrorfWithCtx(ctx, "[Err] Message is older than 10 minutes, cannot delete")
		return fmt.Errorf("message can only be deleted within 10 minutes after sending")
	}

	// Soft delete the message
	if err := s.messageRepo.DeleteMessage(messageID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error deleting message: %v", err)
		return fmt.Errorf("failed to delete message")
	}

	// Broadcast message deletion to both users
	go s.broadcastConversationUpdate(ctx, message.SenderID, message.ConversationID)

	// Get the other user in conversation
	conversation, _ := s.conversationRepo.GetConversationByID(message.ConversationID)
	if conversation != nil {
		otherUserID := conversation.User1ID
		if conversation.User1ID == userID {
			otherUserID = conversation.User2ID
		}
		go s.broadcastConversationUpdate(ctx, otherUserID, message.ConversationID)
	}

	return nil
}
