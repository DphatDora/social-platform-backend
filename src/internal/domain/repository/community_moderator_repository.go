package repository

import "social-platform-backend/internal/domain/model"

type CommunityModeratorRepository interface {
	CreateModerator(moderator *model.CommunityModerator) error
	DeleteModerator(communityID, userID uint64) error
	GetModeratorRole(communityID, userID uint64) (string, error)
}
