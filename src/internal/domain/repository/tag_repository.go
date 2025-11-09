package repository

import "social-platform-backend/internal/domain/model"

type TagRepository interface {
	GetAllTags(search *string) ([]*model.Tag, error)
}
