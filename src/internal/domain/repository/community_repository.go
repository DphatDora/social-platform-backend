package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"
)

type CommunityRepository interface {
	CreateCommunity(community *model.Community) error
	GetCommunityByID(id uint64) (*model.Community, error)
	UpdateCommunity(id uint64, updateCommunity *request.UpdateCommunityRequest) error
	DeleteCommunity(id uint64) error
}
