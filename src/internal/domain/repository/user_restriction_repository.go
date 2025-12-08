package repository

import "social-platform-backend/internal/domain/model"

type UserRestrictionRepository interface {
	CreateRestriction(restriction *model.UserRestriction) error
	GetActiveRestrictionByUserAndCommunity(userID, communityID uint64) (*model.UserRestriction, error)
	GetUserRestrictionHistory(userID uint64, page, limit int) ([]*model.UserRestriction, int64, error)
	DeleteRestriction(id uint64) error
}
