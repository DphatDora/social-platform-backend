package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
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

func (s *MessageService) SendMessage(senderID uint64, req *request.SendMessageRequest) (*response.MessageResponse, error) {
	_, err := s.userRepo.GetUserByID(req.RecipientID)
	if err != nil {
		log.Printf("[Err] Recipient not found: %v", err)
		return nil, fmt.Errorf("recipient not found")
	}

	conversation, err := s.conversationRepo.CreateOrGetConversation(senderID, req.RecipientID)
	if err != nil {
		log.Printf("[Err] Error creating/getting conversation: %v", err)
		return nil, fmt.Errorf("failed to create conversation")
	}

	message := &model.Message{
		ConversationID: conversation.ID,
		SenderID:       senderID,
		Content:        req.Content,
		IsRead:         false,
	}

	if err := s.messageRepo.CreateMessage(message); err != nil {
		log.Printf("[Err] Error creating message: %v", err)
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
			log.Printf("[Warn] Failed to save message attachments: %v", err)
		}
	}

	if err := s.conversationRepo.UpdateLastMessage(conversation.ID, message.ID); err != nil {
		log.Printf("[Warn] Failed to update last message: %v", err)
	}

	fullMessage, err := s.messageRepo.GetMessageByID(message.ID)
	if err != nil {
		log.Printf("[Warn] Failed to get full message: %v", err)
		fullMessage = message
	}

	messageResp := response.NewMessageResponse(fullMessage)

	// Broadcast new message to recipient via SSE
	go s.broadcastNewMessage(req.RecipientID, conversation.ID, messageResp)

	// Broadcast conversation update to both users
	go s.broadcastConversationUpdate(senderID, conversation.ID)
	go s.broadcastConversationUpdate(req.RecipientID, conversation.ID)

	return messageResp, nil
}

// sends new message event via SSE
func (s *MessageService) broadcastNewMessage(userID uint64, conversationID uint64, message *response.MessageResponse) {
	event := &response.SSEEvent{
		Event: "new_message",
		Data: response.NewMessageEvent{
			ConversationID: conversationID,
			Message:        *message,
		},
	}
	s.sseService.BroadcastToUser(userID, event)
}

// sends conversation update event via SSE
func (s *MessageService) broadcastConversationUpdate(userID uint64, conversationID uint64) {
	conversation, err := s.conversationRepo.GetConversationByID(conversationID)
	if err != nil {
		log.Printf("[Warn] Failed to get conversation for broadcast: %v", err)
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
	s.sseService.BroadcastToUser(userID, event)
}

func (s *MessageService) MarkMessageAsRead(userID, messageID uint64) error {
	message, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
		log.Printf("[Err] Message not found: %v", err)
		return fmt.Errorf("message not found")
	}

	isInConversation, err := s.conversationRepo.CheckUserInConversation(message.ConversationID, userID)
	if err != nil || !isInConversation {
		log.Printf("[Err] User not in conversation")
		return fmt.Errorf("unauthorized")
	}

	if err := s.messageRepo.MarkMessageAsRead(messageID, userID); err != nil {
		log.Printf("[Err] Error marking message as read: %v", err)
		return fmt.Errorf("failed to mark message as read")
	}

	// Broadcast read status to sender
	go s.broadcastConversationUpdate(message.SenderID, message.ConversationID)
	go s.broadcastConversationUpdate(userID, message.ConversationID)

	return nil
}

func (s *MessageService) GetConversations(userID uint64, page, limit int) ([]*response.ConversationListResponse, *response.Pagination, error) {
	conversations, total, err := s.conversationRepo.GetUserConversations(userID, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting conversations: %v", err)
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

func (s *MessageService) GetMessages(userID, conversationID uint64, page, limit int) ([]*response.MessageResponse, *response.Pagination, error) {
	isInConversation, err := s.conversationRepo.CheckUserInConversation(conversationID, userID)
	if err != nil || !isInConversation {
		log.Printf("[Err] User not in conversation")
		return nil, nil, fmt.Errorf("unauthorized")
	}

	messages, total, err := s.messageRepo.GetConversationMessages(conversationID, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting messages: %v", err)
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

func (s *MessageService) MarkConversationAsRead(userID, conversationID uint64) error {
	isInConversation, err := s.conversationRepo.CheckUserInConversation(conversationID, userID)
	if err != nil || !isInConversation {
		log.Printf("[Err] User not in conversation")
		return fmt.Errorf("unauthorized")
	}

	if err := s.messageRepo.MarkConversationMessagesAsRead(conversationID, userID); err != nil {
		log.Printf("[Err] Error marking conversation as read: %v", err)
		return fmt.Errorf("failed to mark conversation as read")
	}

	// Broadcast update to both users
	conversation, _ := s.conversationRepo.GetConversationByID(conversationID)
	if conversation != nil {
		go s.broadcastConversationUpdate(conversation.User1ID, conversationID)
		go s.broadcastConversationUpdate(conversation.User2ID, conversationID)
	}

	return nil
}

func (s *MessageService) DeleteMessage(userID, messageID uint64) error {
	message, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
		log.Printf("[Err] Message not found: %v", err)
		return fmt.Errorf("message not found")
	}

	// Check if user is the sender
	if message.SenderID != userID {
		log.Printf("[Err] User is not the sender of this message")
		return fmt.Errorf("unauthorized: only sender can delete message")
	}

	// Check if message is within 10 minutes
	timeSinceCreation := time.Since(message.CreatedAt)
	if timeSinceCreation > 10*time.Minute {
		log.Printf("[Err] Message is older than 10 minutes, cannot delete")
		return fmt.Errorf("message can only be deleted within 10 minutes after sending")
	}

	// Soft delete the message
	if err := s.messageRepo.DeleteMessage(messageID); err != nil {
		log.Printf("[Err] Error deleting message: %v", err)
		return fmt.Errorf("failed to delete message")
	}

	// Broadcast message deletion to both users
	go s.broadcastConversationUpdate(message.SenderID, message.ConversationID)

	// Get the other user in conversation
	conversation, _ := s.conversationRepo.GetConversationByID(message.ConversationID)
	if conversation != nil {
		otherUserID := conversation.User1ID
		if conversation.User1ID == userID {
			otherUserID = conversation.User2ID
		}
		go s.broadcastConversationUpdate(otherUserID, message.ConversationID)
	}

	return nil
}
