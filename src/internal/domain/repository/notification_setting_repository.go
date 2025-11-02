package repository

import "social-platform-backend/internal/domain/model"

type NotificationSettingRepository interface {
	CreateNotificationSetting(setting *model.NotificationSetting) error
	GetUserNotificationSetting(userID uint64, action string) (*model.NotificationSetting, error)
	GetUserNotificationSettings(userID uint64) ([]*model.NotificationSetting, error)
	UpdateNotificationSetting(setting *model.NotificationSetting) error
	UpsertNotificationSetting(setting *model.NotificationSetting) error
}
