package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type BotTaskRepositoryImpl struct {
	db *gorm.DB
}

func NewBotTaskRepository(db *gorm.DB) repository.BotTaskRepository {
	return &BotTaskRepositoryImpl{db: db}
}

func (r *BotTaskRepositoryImpl) CreateBotTask(task *model.BotTask) error {
	return r.db.Create(task).Error
}
