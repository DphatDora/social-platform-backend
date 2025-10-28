package repository

import "social-platform-backend/internal/domain/model"

type CommunityModeratorRepository interface {
	CreateModerator(moderator *model.CommunityModerator) error
	DeleteModerator(communityID, userID uint64) error
	GetModeratorRole(communityID, userID uint64) (string, error)
	GetModeratorCommunitiesByUserID(userID uint64) ([]*model.CommunityModerator, error)
	GetCommunityModerators(communityID uint64) ([]*model.CommunityModerator, error)
}
