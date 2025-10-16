package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/util"
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

func (r *CommunityRepositoryImpl) GetCommunityWithMemberCount(id uint64) (*model.Community, int64, error) {
	var community model.Community
	err := r.db.Where("id = ?", id).First(&community).Error
	if err != nil {
		return nil, 0, err
	}

	var memberCount int64
	r.db.Model(&model.Subscription{}).
		Where("community_id = ?", id).
		Count(&memberCount)

	return &community, memberCount, nil
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

func (r *CommunityRepositoryImpl) GetCommunities(page, limit int) ([]*model.Community, int64, error) {
	var communities []*model.Community
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&model.Community{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Select("communities.*, COUNT(subscriptions.user_id) as member_count").
		Joins("LEFT JOIN subscriptions ON subscriptions.community_id = communities.id").
		Group("communities.id").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&communities).Error

	return communities, total, err
}

func (r *CommunityRepositoryImpl) SearchCommunitiesByName(name string, page, limit int) ([]*model.Community, int64, error) {
	var communities []*model.Community
	var total int64

	offset := (page - 1) * limit

	patterns := util.BuildSearchPattern(name)
	query := r.db.Model(&model.Community{})

	for _, p := range patterns {
		query = query.Where("unaccent(lower(name)) LIKE unaccent(lower(?))", p)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Reset query for actual data fetch with JOIN
	query = r.db.Model(&model.Community{})
	for _, p := range patterns {
		query = query.Where("unaccent(lower(name)) LIKE unaccent(lower(?))", p)
	}

	err := query.Select("communities.*, COUNT(subscriptions.user_id) as member_count").
		Joins("LEFT JOIN subscriptions ON subscriptions.community_id = communities.id").
		Group("communities.id").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&communities).Error
	if err != nil {
		return nil, 0, err
	}

	return communities, total, nil
}

func (r *CommunityRepositoryImpl) FilterCommunities(sortBy string, isPrivate *bool, page, limit int) ([]*model.Community, int64, error) {
	var communities []*model.Community
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&model.Community{}).
		Select("communities.*, COUNT(subscriptions.user_id) as member_count").
		Joins("LEFT JOIN subscriptions ON subscriptions.community_id = communities.id").
		Group("communities.id")

	if isPrivate != nil {
		query = query.Where("is_private = ?", *isPrivate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch sortBy {
	case constant.SORT_MEMBER_COUNT:
		query = query.Order("member_count DESC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Limit(limit).Offset(offset).Find(&communities).Error

	return communities, total, err
}
