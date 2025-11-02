package model

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID        uint64           `gorm:"column:id;primaryKey"`
	UserID    uint64           `gorm:"column:user_id"`
	Body      string           `gorm:"column:body"`
	Action    string           `gorm:"column:action"`
	Payload   *json.RawMessage `gorm:"column:payload"`
	IsRead    bool             `gorm:"column:is_read"`
	CreatedAt time.Time        `gorm:"column:created_at"`

	// relations
	User *User `gorm:"foreignKey:UserID;references:ID"`
}

func (Notification) TableName() string {
	return "notifications"
}
