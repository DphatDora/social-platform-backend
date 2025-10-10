package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"time"

	"gorm.io/gorm"
)

type CommunityRepositoryImpl struct {
	db *gorm.DB
}

func NewCommunityRepository(db *gorm.DB) repository.CommunityRepository {
	return &CommunityRepositoryImpl{db: db}
}

func (r *CommunityRepositoryImpl) CreateCommunity(community *model.Community) error {
	return r.db.Create(community).Error
}

func (r *CommunityRepositoryImpl) GetCommunityByID(id uint64) (*model.Community, error) {
	var community model.Community
	err := r.db.Where("id = ?", id).First(&community).Error
	if err != nil {
		return nil, err
	}
	return &community, nil
}

func (r *CommunityRepositoryImpl) UpdateCommunity(id uint64, updateCommunity *request.UpdateCommunityRequest) error {
	updates := make(map[string]interface{})
	if updateCommunity.Name != nil {
		updates["name"] = *updateCommunity.Name
	}
	if updateCommunity.ShortDescription != nil {
		updates["short_description"] = *updateCommunity.ShortDescription
	}
	if updateCommunity.Description != nil {
		updates["description"] = *updateCommunity.Description
	}
	if updateCommunity.CoverImage != nil {
		updates["cover_image"] = *updateCommunity.CoverImage
	}
	if updateCommunity.IsPrivate != nil {
		updates["is_private"] = *updateCommunity.IsPrivate
	}
	return r.db.Model(&model.Community{}).Where("id = ?", id).Updates(updates).Error
}

func (r *CommunityRepositoryImpl) DeleteCommunity(id uint64) error {
	return r.db.Model(&model.Community{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
