package repository

import "social-platform-backend/internal/domain/model"

type BotTaskRepository interface {
	CreateBotTask(task *model.BotTask) error
}
