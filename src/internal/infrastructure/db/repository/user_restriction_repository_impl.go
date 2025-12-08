package repository

import (
	"social-platform-backend/internal/domain/model"
	"time"

	"gorm.io/gorm"
)

type UserRestrictionRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRestrictionRepository(db *gorm.DB) *UserRestrictionRepositoryImpl {
	return &UserRestrictionRepositoryImpl{db: db}
}

func (r *UserRestrictionRepositoryImpl) CreateRestriction(restriction *model.UserRestriction) error {
	return r.db.Create(restriction).Error
}

func (r *UserRestrictionRepositoryImpl) GetActiveRestrictionByUserAndCommunity(userID, communityID uint64) (*model.UserRestriction, error) {
	var restriction model.UserRestriction
	now := time.Now()

	err := r.db.Where("user_id = ? AND community_id = ? AND (expires_at IS NULL OR expires_at > ?)",
		userID, communityID, now).
		Order("created_at DESC").
		First(&restriction).Error

	if err != nil {
		return nil, err
	}

	return &restriction, nil
}

func (r *UserRestrictionRepositoryImpl) GetUserRestrictionHistory(userID uint64, page, limit int) ([]*model.UserRestriction, int64, error) {
	var restrictions []*model.UserRestriction
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&model.UserRestriction{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&restrictions).Error

	if err != nil {
		return nil, 0, err
	}

	return restrictions, total, nil
}

func (r *UserRestrictionRepositoryImpl) DeleteRestriction(id uint64) error {
	return r.db.Delete(&model.UserRestriction{}, id).Error
}
