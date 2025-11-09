package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type TopicRepositoryImpl struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) repository.TopicRepository {
	return &TopicRepositoryImpl{db: db}
}

func (r *TopicRepositoryImpl) GetAllTopics(search *string) ([]*model.Topic, error) {
	var topics []*model.Topic
	query := r.db.Model(&model.Topic{})

	if search != nil && *search != "" {
		query = query.Where("name ILIKE ?", "%"+*search+"%")
	}

	err := query.Order("name ASC").Find(&topics).Error
	if err != nil {
		return nil, err
	}

	return topics, nil
}
