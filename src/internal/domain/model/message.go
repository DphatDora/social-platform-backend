package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID             uint64         `gorm:"column:id;primaryKey"`
	ConversationID uint64         `gorm:"column:conversation_id"`
	SenderID       uint64         `gorm:"column:sender_id"`
	Sender         *User          `gorm:"foreignKey:SenderID;references:ID"`
	Type           string         `gorm:"column:type"`
	Content        string         `gorm:"column:content"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	IsRead         bool           `gorm:"column:is_read"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// relation
	Conversation *Conversation `gorm:"foreignKey:ConversationID;references:ID"`
}

func (Message) TableName() string {
	return "messages"
}
