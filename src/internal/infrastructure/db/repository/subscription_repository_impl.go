package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type SubscriptionRepositoryImpl struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) repository.SubscriptionRepository {
	return &SubscriptionRepositoryImpl{db: db}
}

func (r *SubscriptionRepositoryImpl) CreateSubscription(subscription *model.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *SubscriptionRepositoryImpl) DeleteSubscription(userID, communityID uint64) error {
	return r.db.Where("user_id = ? AND community_id = ?", userID, communityID).
		Delete(&model.Subscription{}).Error
}

func (r *SubscriptionRepositoryImpl) IsUserSubscribed(userID, communityID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Subscription{}).
		Where("user_id = ? AND community_id = ?", userID, communityID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
