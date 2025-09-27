package model

import (
	"encoding/json"
	"time"
)

type BotTask struct {
	ID         uint64           `gorm:"column:id;primaryKey"`
	Action     string           `gorm:"column:action"`
	Payload    *json.RawMessage `gorm:"column:payload"`
	CreatedAt  time.Time        `gorm:"column:created_at;autoCreateTime"`
	ExecutedAt *time.Time       `gorm:"column:executed_at"`
}

func (BotTask) TableName() string {
	return "bot_tasks"
}
