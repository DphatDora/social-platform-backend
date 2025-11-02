package response

import (
	"encoding/json"
	"social-platform-backend/internal/domain/model"
	"time"
)

type NotificationResponse struct {
	ID        uint64           `json:"id"`
	Body      string           `json:"body"`
	Action    string           `json:"action"`
	Payload   *json.RawMessage `json:"payload"`
	IsRead    bool             `json:"isRead"`
	CreatedAt time.Time        `json:"createdAt"`
}

func NewNotificationResponse(notification *model.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID:        notification.ID,
		Body:      notification.Body,
		Action:    notification.Action,
		Payload:   notification.Payload,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt,
	}
}

type NewNotificationEvent struct {
	Notification NotificationResponse `json:"notification"`
	UnreadCount  int64                `json:"unreadCount"`
}
