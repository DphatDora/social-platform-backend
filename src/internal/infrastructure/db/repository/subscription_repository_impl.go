package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/package/util"

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

func (r *SubscriptionRepositoryImpl) GetCommunityMembers(communityID uint64, sortBy, searchName string, page, limit int) ([]*model.Subscription, int64, error) {
	var subscriptions []*model.Subscription
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&model.Subscription{}).
		Preload("User").
		Where("community_id = ?", communityID)

	// Apply search
	if searchName != "" {
		patterns := util.BuildSearchPattern(searchName)
		query = query.Joins("JOIN users ON users.id = subscriptions.user_id")
		for _, p := range patterns {
			query = query.Where("unaccent(lower(users.username)) LIKE unaccent(lower(?))", p)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Reset query for actual data fetch with sorting
	query = r.db.Model(&model.Subscription{}).
		Preload("User").
		Where("community_id = ?", communityID)

	// Re-apply search filter for data fetch
	if searchName != "" {
		patterns := util.BuildSearchPattern(searchName)
		query = query.Joins("JOIN users ON users.id = subscriptions.user_id")
		for _, p := range patterns {
			query = query.Where("unaccent(lower(users.username)) LIKE unaccent(lower(?))", p)
		}
	}

	// Apply sorting
	switch sortBy {
	case "oldest":
		query = query.Order("subscriptions.subscribed_at ASC")
	case "karma":
		if searchName == "" {
			query = query.Joins("JOIN users ON users.id = subscriptions.user_id")
		}
		query = query.Order("users.karma DESC")
	default:
		query = query.Order("subscriptions.subscribed_at DESC")
	}

	err := query.Limit(limit).Offset(offset).Find(&subscriptions).Error
	if err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}
