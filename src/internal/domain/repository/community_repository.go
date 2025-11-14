package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"
)

type CommunityRepository interface {
	CreateCommunity(community *model.Community) error
	GetCommunityByID(id uint64) (*model.Community, error)
	GetCommunityWithMemberCount(id uint64) (*model.Community, int64, error)
	GetCommunityByIDWithUserSubscription(communityID uint64, userID *uint64) (*model.Community, int64, error)
	UpdateCommunity(id uint64, updateCommunity *request.UpdateCommunityRequest) error
	DeleteCommunity(id uint64) error
	GetCommunities(page, limit int, userID *uint64) ([]*model.Community, int64, error)
	SearchCommunitiesByName(name string, page, limit int, userID *uint64) ([]*model.Community, int64, error)
	FilterCommunities(sortBy string, isPrivate *bool, topics []string, page, limit int, userID *uint64) ([]*model.Community, int64, error)
	GetCommunitiesByCreatorID(creatorID uint64) ([]*model.Community, error)
	IsCommunityNameExists(name string) (bool, error)
	UpdateRequiresPostApproval(id uint64, requiresPostApproval bool) error
	UpdateRequiresMemberApproval(id uint64, requiresMemberApproval bool) error
}
