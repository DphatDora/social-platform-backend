package repository

import "social-platform-backend/internal/domain/model"

type TopicRepository interface {
	GetAllTopics(search *string) ([]*model.Topic, error)
}
