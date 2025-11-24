package service

import (
	"errors"
	"testing"

	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func uint64Ptr(v uint64) *uint64 {
	return &v
}

func TestMessageService_SendMessage_Success(t *testing.T) {
	mockConversationRepo := new(MockConversationRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockUserRepo := new(MockUserRepository)
	sseService := NewSSEService() // Create real SSE service to avoid nil pointer

	messageService := NewMessageService(
		mockConversationRepo,
		mockMessageRepo,
		nil,
		mockUserRepo,
		sseService,
	)

	senderID := uint64(123)
	req := &request.SendMessageRequest{
		RecipientID: 456,
		Content:     "Hello, world!",
	}

	recipient := &model.User{
		ID:       req.RecipientID,
		Username: "recipient",
	}

	conversation := &model.Conversation{
		ID:      789,
		User1ID: senderID,
		User2ID: req.RecipientID,
	}

	message := &model.Message{
		ID:             1,
		ConversationID: conversation.ID,
		SenderID:       senderID,
		Content:        req.Content,
		IsRead:         false,
	}

	mockUserRepo.On("GetUserByID", req.RecipientID).Return(recipient, nil)
	mockConversationRepo.On("CreateOrGetConversation", senderID, req.RecipientID).Return(conversation, nil)
	mockMessageRepo.On("CreateMessage", mock.AnythingOfType("*model.Message")).Return(nil)
	mockConversationRepo.On("UpdateLastMessage", conversation.ID, mock.AnythingOfType("uint64")).Return(nil)
	mockMessageRepo.On("GetMessageByID", mock.AnythingOfType("uint64")).Return(message, nil)
	// Background goroutine calls
	mockConversationRepo.On("GetConversationByID", conversation.ID).Return(conversation, nil).Maybe()
	mockMessageRepo.On("GetUnreadCount", conversation.ID, senderID).Return(int64(0), nil).Maybe()
	mockMessageRepo.On("GetUnreadCount", conversation.ID, req.RecipientID).Return(int64(1), nil).Maybe()

	result, err := messageService.SendMessage(senderID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockUserRepo.AssertExpectations(t)
	mockConversationRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

func TestMessageService_SendMessage_RecipientNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	messageService := NewMessageService(
		nil,
		nil,
		nil,
		mockUserRepo,
		nil,
	)

	req := &request.SendMessageRequest{
		RecipientID: 999,
		Content:     "Hello",
	}

	mockUserRepo.On("GetUserByID", req.RecipientID).Return(nil, errors.New("not found"))

	result, err := messageService.SendMessage(123, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "recipient not found")
	mockUserRepo.AssertExpectations(t)
}

func TestMessageService_SendMessage_ConversationError(t *testing.T) {
	mockConversationRepo := new(MockConversationRepository)
	mockUserRepo := new(MockUserRepository)

	messageService := NewMessageService(
		mockConversationRepo,
		nil,
		nil,
		mockUserRepo,
		nil,
	)

	senderID := uint64(123)
	req := &request.SendMessageRequest{
		RecipientID: 456,
		Content:     "Hello",
	}

	recipient := &model.User{ID: req.RecipientID}

	mockUserRepo.On("GetUserByID", req.RecipientID).Return(recipient, nil)
	mockConversationRepo.On("CreateOrGetConversation", senderID, req.RecipientID).Return(nil, errors.New("db error"))

	result, err := messageService.SendMessage(senderID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create conversation")
	mockUserRepo.AssertExpectations(t)
	mockConversationRepo.AssertExpectations(t)
}

func TestMessageService_MarkMessageAsRead_Success(t *testing.T) {
	mockConversationRepo := new(MockConversationRepository)
	mockMessageRepo := new(MockMessageRepository)
	sseService := NewSSEService()

	messageService := NewMessageService(
		mockConversationRepo,
		mockMessageRepo,
		nil,
		nil,
		sseService,
	)

	userID := uint64(123)
	messageID := uint64(456)

	message := &model.Message{
		ID:             messageID,
		ConversationID: 789,
		SenderID:       999,
		Content:        "Test message",
		IsRead:         false,
	}

	conversation := &model.Conversation{
		ID:      789,
		User1ID: userID,
		User2ID: 999,
	}

	mockMessageRepo.On("GetMessageByID", messageID).Return(message, nil)
	mockConversationRepo.On("CheckUserInConversation", message.ConversationID, userID).Return(true, nil)
	mockMessageRepo.On("MarkMessageAsRead", messageID, userID).Return(nil)
	// Background goroutine calls
	mockConversationRepo.On("GetConversationByID", message.ConversationID).Return(conversation, nil).Maybe()
	mockMessageRepo.On("GetUnreadCount", message.ConversationID, message.SenderID).Return(int64(0), nil).Maybe()
	mockMessageRepo.On("GetUnreadCount", message.ConversationID, userID).Return(int64(0), nil).Maybe()

	err := messageService.MarkMessageAsRead(userID, messageID)

	assert.NoError(t, err)
	mockMessageRepo.AssertExpectations(t)
	mockConversationRepo.AssertExpectations(t)
}

func TestMessageService_MarkMessageAsRead_MessageNotFound(t *testing.T) {
	mockMessageRepo := new(MockMessageRepository)

	messageService := NewMessageService(
		nil,
		mockMessageRepo,
		nil,
		nil,
		nil,
	)

	messageID := uint64(999)
	mockMessageRepo.On("GetMessageByID", messageID).Return(nil, errors.New("not found"))

	err := messageService.MarkMessageAsRead(123, messageID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "message not found")
	mockMessageRepo.AssertExpectations(t)
}

func TestMessageService_MarkMessageAsRead_Unauthorized(t *testing.T) {
	mockConversationRepo := new(MockConversationRepository)
	mockMessageRepo := new(MockMessageRepository)
	sseService := NewSSEService()

	messageService := NewMessageService(
		mockConversationRepo,
		mockMessageRepo,
		nil,
		nil,
		sseService,
	)

	userID := uint64(123)
	messageID := uint64(456)

	message := &model.Message{
		ID:             messageID,
		ConversationID: 789,
		SenderID:       999,
		Content:        "Test message",
	}

	mockMessageRepo.On("GetMessageByID", messageID).Return(message, nil)
	mockConversationRepo.On("CheckUserInConversation", message.ConversationID, userID).Return(false, nil)

	err := messageService.MarkMessageAsRead(userID, messageID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	mockMessageRepo.AssertExpectations(t)
	mockConversationRepo.AssertExpectations(t)
}

func TestMessageService_GetConversations_Success(t *testing.T) {
	mockConversationRepo := new(MockConversationRepository)
	mockMessageRepo := new(MockMessageRepository)

	messageService := NewMessageService(
		mockConversationRepo,
		mockMessageRepo,
		nil,
		nil,
		nil,
	)

	userID := uint64(123)
	page := 1
	limit := 10

	conversations := []*model.Conversation{
		{
			ID:            1,
			User1ID:       userID,
			User2ID:       456,
			LastMessageID: uint64Ptr(10),
		},
		{
			ID:            2,
			User1ID:       userID,
			User2ID:       789,
			LastMessageID: uint64Ptr(20),
		},
	}

	mockConversationRepo.On("GetUserConversations", userID, page, limit).Return(conversations, int64(2), nil)
	mockMessageRepo.On("GetUnreadCount", uint64(1), userID).Return(int64(3), nil)
	mockMessageRepo.On("GetUnreadCount", uint64(2), userID).Return(int64(0), nil)

	result, pagination, err := messageService.GetConversations(userID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, pagination)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), pagination.Total)
	mockConversationRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

func TestMessageService_GetConversations_EmptyResult(t *testing.T) {
	mockConversationRepo := new(MockConversationRepository)

	messageService := NewMessageService(
		mockConversationRepo,
		nil,
		nil,
		nil,
		nil,
	)

	userID := uint64(123)
	page := 1
	limit := 10

	mockConversationRepo.On("GetUserConversations", userID, page, limit).Return([]*model.Conversation{}, int64(0), nil)

	result, pagination, err := messageService.GetConversations(userID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, pagination)
	assert.Len(t, result, 0)
	mockConversationRepo.AssertExpectations(t)
}

func TestMessageService_GetMessages_Success(t *testing.T) {
	mockConversationRepo := new(MockConversationRepository)
	mockMessageRepo := new(MockMessageRepository)

	messageService := NewMessageService(
		mockConversationRepo,
		mockMessageRepo,
		nil,
		nil,
		nil,
	)

	userID := uint64(123)
	conversationID := uint64(789)
	page := 1
	limit := 20

	messages := []*model.Message{
		{
			ID:             1,
			ConversationID: conversationID,
			SenderID:       userID,
			Content:        "Hello",
			IsRead:         true,
		},
		{
			ID:             2,
			ConversationID: conversationID,
			SenderID:       456,
			Content:        "Hi there!",
			IsRead:         false,
		},
	}

	mockConversationRepo.On("CheckUserInConversation", conversationID, userID).Return(true, nil)
	mockMessageRepo.On("GetConversationMessages", conversationID, page, limit).Return(messages, int64(2), nil)

	result, pagination, err := messageService.GetMessages(userID, conversationID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, pagination)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), pagination.Total)
	mockConversationRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

func TestMessageService_GetMessages_Unauthorized(t *testing.T) {
	mockConversationRepo := new(MockConversationRepository)

	messageService := NewMessageService(
		mockConversationRepo,
		nil,
		nil,
		nil,
		nil,
	)

	userID := uint64(123)
	conversationID := uint64(789)
	page := 1
	limit := 20

	mockConversationRepo.On("CheckUserInConversation", conversationID, userID).Return(false, nil)

	result, pagination, err := messageService.GetMessages(userID, conversationID, page, limit)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, pagination)
	assert.Contains(t, err.Error(), "unauthorized")
	mockConversationRepo.AssertExpectations(t)
}

func TestMessageService_GetMessages_ConversationCheckError(t *testing.T) {
	mockConversationRepo := new(MockConversationRepository)

	messageService := NewMessageService(
		mockConversationRepo,
		nil,
		nil,
		nil,
		nil,
	)

	userID := uint64(123)
	conversationID := uint64(789)
	page := 1
	limit := 20

	mockConversationRepo.On("CheckUserInConversation", conversationID, userID).Return(false, errors.New("db error"))

	result, pagination, err := messageService.GetMessages(userID, conversationID, page, limit)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, pagination)
	mockConversationRepo.AssertExpectations(t)
}
