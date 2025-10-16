package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"
)

type CommunityRepository interface {
	CreateCommunity(community *model.Community) error
	GetCommunityByID(id uint64) (*model.Community, error)
	GetCommunityWithMemberCount(id uint64) (*model.Community, int64, error)
	UpdateCommunity(id uint64, updateCommunity *request.UpdateCommunityRequest) error
	DeleteCommunity(id uint64) error
	GetCommunities(page, limit int) ([]*model.Community, int64, error)
	SearchCommunitiesByName(name string, page, limit int) ([]*model.Community, int64, error)
	FilterCommunities(sortBy string, isPrivate *bool, page, limit int) ([]*model.Community, int64, error)
	GetCommunityRole(userID, communityID uint64) (string, error)
}
