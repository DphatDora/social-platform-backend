package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type TagRepositoryImpl struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &TagRepositoryImpl{db: db}
}

func (r *TagRepositoryImpl) GetAllTags(search *string) ([]*model.Tag, error) {
	var tags []*model.Tag
	query := r.db.Model(&model.Tag{})

	if search != nil && *search != "" {
		query = query.Where("name ILIKE ?", "%"+*search+"%")
	}

	err := query.Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}
