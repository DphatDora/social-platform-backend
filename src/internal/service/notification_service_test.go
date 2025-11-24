package service

import (
	"encoding/json"
	"errors"
	"testing"

	"social-platform-backend/internal/domain/model"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/template/payload"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificationService_CreateNotification_Success(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)
	mockNotificationSettingRepo := new(MockNotificationSettingRepository)
	mockUserRepo := new(MockUserRepository)
	sseService := NewSSEService()

	notificationService := NewNotificationService(
		mockNotificationRepo,
		mockNotificationSettingRepo,
		nil,
		mockUserRepo,
		sseService,
		nil,
	)

	userID := uint64(123)
	action := constant.NOTIFICATION_ACTION_GET_POST_VOTE
	notifPayload := payload.PostVoteNotificationPayload{
		PostID:   456,
		UserName: "testuser",
		VoteType: true,
	}

	user := &model.User{
		ID:       userID,
		Email:    "test@example.com",
		Username: "testuser",
	}

	setting := &model.NotificationSetting{
		UserID:     userID,
		Action:     action,
		IsPush:     true,
		IsSendMail: false,
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockNotificationSettingRepo.On("GetUserNotificationSetting", userID, action).Return(setting, nil)
	mockNotificationRepo.On("CreateNotification", mock.AnythingOfType("*model.Notification")).Return(nil)
	// Background goroutine call
	mockNotificationRepo.On("GetUnreadCount", userID).Return(int64(1), nil).Maybe()

	err := notificationService.CreateNotification(userID, action, notifPayload)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockNotificationSettingRepo.AssertExpectations(t)
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_CreateNotification_UserNotFound(t *testing.T) {
	mockNotificationSettingRepo := new(MockNotificationSettingRepository)
	mockUserRepo := new(MockUserRepository)
	sseService := NewSSEService()

	notificationService := NewNotificationService(
		nil,
		mockNotificationSettingRepo,
		nil,
		mockUserRepo,
		sseService,
		nil,
	)

	userID := uint64(999)
	action := constant.NOTIFICATION_ACTION_GET_POST_VOTE
	notifPayload := payload.PostVoteNotificationPayload{}

	setting := &model.NotificationSetting{
		UserID:     userID,
		Action:     action,
		IsPush:     true,
		IsSendMail: false,
	}

	mockNotificationSettingRepo.On("GetUserNotificationSetting", userID, action).Return(setting, nil)
	mockUserRepo.On("GetUserByID", userID).Return(nil, errors.New("not found"))

	err := notificationService.CreateNotification(userID, action, notifPayload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user")
	mockUserRepo.AssertExpectations(t)
	mockNotificationSettingRepo.AssertExpectations(t)
}

func TestNotificationService_CreateNotification_DefaultSettings(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)
	mockNotificationSettingRepo := new(MockNotificationSettingRepository)
	mockUserRepo := new(MockUserRepository)
	sseService := NewSSEService()

	notificationService := NewNotificationService(
		mockNotificationRepo,
		mockNotificationSettingRepo,
		nil,
		mockUserRepo,
		sseService,
		nil,
	)

	userID := uint64(123)
	action := constant.NOTIFICATION_ACTION_GET_POST_NEW_COMMENT
	notifPayload := payload.PostCommentNotificationPayload{
		PostID:   456,
		UserName: "commenter",
	}

	user := &model.User{
		ID:       userID,
		Email:    "test@example.com",
		Username: "testuser",
	}

	// Return error to simulate no setting found
	mockNotificationSettingRepo.On("GetUserNotificationSetting", userID, action).Return(nil, errors.New("not found"))
	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockNotificationRepo.On("CreateNotification", mock.AnythingOfType("*model.Notification")).Return(nil)
	// Background goroutine call
	mockNotificationRepo.On("GetUnreadCount", userID).Return(int64(1), nil).Maybe()

	err := notificationService.CreateNotification(userID, action, notifPayload)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockNotificationSettingRepo.AssertExpectations(t)
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_GetUserNotifications_Success(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)

	notificationService := NewNotificationService(
		mockNotificationRepo,
		nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	page := 1
	limit := 10

	rawPayload := json.RawMessage(`{"postId":456,"userName":"testuser"}`)
	notifications := []*model.Notification{
		{
			ID:      1,
			UserID:  userID,
			Body:    "User testuser upvoted your post",
			Action:  constant.NOTIFICATION_ACTION_GET_POST_VOTE,
			Payload: &rawPayload,
			IsRead:  false,
		},
		{
			ID:      2,
			UserID:  userID,
			Body:    "User testuser commented on your post",
			Action:  constant.NOTIFICATION_ACTION_GET_POST_NEW_COMMENT,
			Payload: &rawPayload,
			IsRead:  true,
		},
	}

	mockNotificationRepo.On("GetUserNotifications", userID, limit, 0).Return(notifications, int64(2), nil)

	result, pagination, err := notificationService.GetUserNotifications(userID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, pagination)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), pagination.Total)
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_GetUserNotifications_EmptyResult(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)

	notificationService := NewNotificationService(
		mockNotificationRepo,
		nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	page := 1
	limit := 10

	mockNotificationRepo.On("GetUserNotifications", userID, limit, 0).Return([]*model.Notification{}, int64(0), nil)

	result, pagination, err := notificationService.GetUserNotifications(userID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, pagination)
	assert.Len(t, result, 0)
	assert.Equal(t, int64(0), pagination.Total)
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAsRead_Success(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)

	notificationService := NewNotificationService(
		mockNotificationRepo,
		nil, nil, nil, nil, nil,
	)

	notificationID := uint64(456)
	userID := uint64(123)

	rawPayload := json.RawMessage(`{"postId":789}`)
	notification := &model.Notification{
		ID:      notificationID,
		UserID:  userID,
		Body:    "Test notification",
		Action:  constant.NOTIFICATION_ACTION_GET_POST_VOTE,
		Payload: &rawPayload,
		IsRead:  false,
	}

	mockNotificationRepo.On("GetNotificationByID", notificationID).Return(notification, nil)
	mockNotificationRepo.On("MarkAsRead", notificationID).Return(nil)

	err := notificationService.MarkAsRead(userID, notificationID)

	assert.NoError(t, err)
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAsRead_NotificationNotFound(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)

	notificationService := NewNotificationService(
		mockNotificationRepo,
		nil, nil, nil, nil, nil,
	)

	notificationID := uint64(999)
	userID := uint64(123)

	mockNotificationRepo.On("GetNotificationByID", notificationID).Return(nil, errors.New("not found"))

	err := notificationService.MarkAsRead(userID, notificationID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "notification not found")
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAsRead_Unauthorized(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)

	notificationService := NewNotificationService(
		mockNotificationRepo,
		nil, nil, nil, nil, nil,
	)

	notificationID := uint64(456)
	userID := uint64(123)

	rawPayload := json.RawMessage(`{"postId":789}`)
	notification := &model.Notification{
		ID:      notificationID,
		UserID:  999, // Different user
		Body:    "Test notification",
		Action:  constant.NOTIFICATION_ACTION_GET_POST_VOTE,
		Payload: &rawPayload,
		IsRead:  false,
	}

	mockNotificationRepo.On("GetNotificationByID", notificationID).Return(notification, nil)

	err := notificationService.MarkAsRead(userID, notificationID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAllAsRead_Success(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)

	notificationService := NewNotificationService(
		mockNotificationRepo,
		nil, nil, nil, nil, nil,
	)

	userID := uint64(123)

	mockNotificationRepo.On("MarkAllAsRead", userID).Return(nil)

	err := notificationService.MarkAllAsRead(userID)

	assert.NoError(t, err)
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_DeleteNotification_Success(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)

	notificationService := NewNotificationService(
		mockNotificationRepo,
		nil, nil, nil, nil, nil,
	)

	notificationID := uint64(456)
	userID := uint64(123)

	rawPayload := json.RawMessage(`{"postId":789}`)
	notification := &model.Notification{
		ID:      notificationID,
		UserID:  userID,
		Body:    "Test notification",
		Action:  constant.NOTIFICATION_ACTION_GET_POST_VOTE,
		Payload: &rawPayload,
		IsRead:  false,
	}

	mockNotificationRepo.On("GetNotificationByID", notificationID).Return(notification, nil)
	mockNotificationRepo.On("DeleteNotification", notificationID).Return(nil)

	err := notificationService.DeleteNotification(userID, notificationID)

	assert.NoError(t, err)
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_DeleteNotification_Unauthorized(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)

	notificationService := NewNotificationService(
		mockNotificationRepo,
		nil, nil, nil, nil, nil,
	)

	notificationID := uint64(456)
	userID := uint64(123)

	rawPayload := json.RawMessage(`{"postId":789}`)
	notification := &model.Notification{
		ID:      notificationID,
		UserID:  999, // Different user
		Body:    "Test notification",
		Action:  constant.NOTIFICATION_ACTION_GET_POST_VOTE,
		Payload: &rawPayload,
		IsRead:  false,
	}

	mockNotificationRepo.On("GetNotificationByID", notificationID).Return(notification, nil)

	err := notificationService.DeleteNotification(userID, notificationID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_GetUnreadCount_Success(t *testing.T) {
	mockNotificationRepo := new(MockNotificationRepository)

	notificationService := NewNotificationService(
		mockNotificationRepo,
		nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	expectedCount := int64(5)

	mockNotificationRepo.On("GetUnreadCount", userID).Return(expectedCount, nil)

	count, err := notificationService.GetUnreadCount(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockNotificationRepo.AssertExpectations(t)
}
