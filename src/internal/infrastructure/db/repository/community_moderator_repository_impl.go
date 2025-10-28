package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type CommunityModeratorRepositoryImpl struct {
	db *gorm.DB
}

func NewCommunityModeratorRepository(db *gorm.DB) repository.CommunityModeratorRepository {
	return &CommunityModeratorRepositoryImpl{db: db}
}

func (r *CommunityModeratorRepositoryImpl) CreateModerator(moderator *model.CommunityModerator) error {
	return r.db.Create(moderator).Error
}

func (r *CommunityModeratorRepositoryImpl) DeleteModerator(communityID, userID uint64) error {
	return r.db.Where("community_id = ? AND user_id = ?", communityID, userID).
		Delete(&model.CommunityModerator{}).Error
}

func (r *CommunityModeratorRepositoryImpl) GetModeratorRole(communityID, userID uint64) (string, error) {
	var role string
	err := r.db.Model(&model.CommunityModerator{}).
		Select("role").
		Where("community_id = ? AND user_id = ?", communityID, userID).
		Scan(&role).Error

	if err != nil {
		return "", err
	}

	return role, nil
}

func (r *CommunityModeratorRepositoryImpl) GetModeratorCommunitiesByUserID(userID uint64) ([]*model.CommunityModerator, error) {
	var moderators []*model.CommunityModerator
	err := r.db.Where("user_id = ?", userID).
		Find(&moderators).Error
	if err != nil {
		return nil, err
	}
	return moderators, nil
}

func (r *CommunityModeratorRepositoryImpl) GetCommunityModerators(communityID uint64) ([]*model.CommunityModerator, error) {
	var moderators []*model.CommunityModerator
	err := r.db.Preload("User").
		Where("community_id = ?", communityID).
		Find(&moderators).Error
	if err != nil {
		return nil, err
	}
	return moderators, nil
}
