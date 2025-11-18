package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type notificationSettingRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationSettingRepository(db *gorm.DB) repository.NotificationSettingRepository {
	return &notificationSettingRepositoryImpl{db: db}
}

func (r *notificationSettingRepositoryImpl) CreateNotificationSetting(setting *model.NotificationSetting) error {
	return r.db.Create(setting).Error
}

func (r *notificationSettingRepositoryImpl) CreateNotificationSettings(settings []*model.NotificationSetting) error {
	if len(settings) == 0 {
		return nil
	}
	return r.db.Create(settings).Error
}

func (r *notificationSettingRepositoryImpl) GetUserNotificationSetting(userID uint64, action string) (*model.NotificationSetting, error) {
	var setting model.NotificationSetting
	err := r.db.Where("user_id = ? AND action = ?", userID, action).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *notificationSettingRepositoryImpl) GetUserNotificationSettings(userID uint64) ([]*model.NotificationSetting, error) {
	var settings []*model.NotificationSetting
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&settings).Error
	return settings, err
}

func (r *notificationSettingRepositoryImpl) UpdateNotificationSetting(setting *model.NotificationSetting) error {
	return r.db.Save(setting).Error
}

func (r *notificationSettingRepositoryImpl) UpsertNotificationSetting(setting *model.NotificationSetting) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "action"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_push", "is_send_mail", "updated_at"}),
	}).Create(setting).Error
}
