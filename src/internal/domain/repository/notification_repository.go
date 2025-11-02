package repository

import "social-platform-backend/internal/domain/model"

type NotificationRepository interface {
	CreateNotification(notification *model.Notification) error
	GetNotificationByID(id uint64) (*model.Notification, error)
	GetUserNotifications(userID uint64, limit, offset int) ([]*model.Notification, int64, error)
	MarkAsRead(id uint64) error
	MarkAllAsRead(userID uint64) error
	DeleteNotification(id uint64) error
	GetUnreadCount(userID uint64) (int64, error)
}
