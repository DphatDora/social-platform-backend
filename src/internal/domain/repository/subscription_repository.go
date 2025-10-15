package repository

import "social-platform-backend/internal/domain/model"

type SubscriptionRepository interface {
	CreateSubscription(subscription *model.Subscription) error
	DeleteSubscription(userID, communityID uint64) error
	IsUserSubscribed(userID, communityID uint64) (bool, error)
	GetCommunityMembers(communityID uint64, sortBy, searchName string, page, limit int) ([]*model.Subscription, int64, error)
}
