package repository

import (
	"fmt"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type notificationRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &notificationRepositoryImpl{db: db}
}

func (r *notificationRepositoryImpl) CreateNotification(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepositoryImpl) GetNotificationByID(id uint64) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.Preload("User").Where("id = ?", id).First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepositoryImpl) GetUserNotifications(userID uint64, limit, offset int) ([]*model.Notification, int64, error) {
	var notifications []*model.Notification
	var total int64

	if err := r.db.Model(&model.Notification{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *notificationRepositoryImpl) MarkAsRead(id uint64) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *notificationRepositoryImpl) MarkAllAsRead(userID uint64) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r *notificationRepositoryImpl) DeleteNotification(id uint64) error {
	result := r.db.Delete(&model.Notification{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}
	return nil
}

func (r *notificationRepositoryImpl) GetUnreadCount(userID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}
