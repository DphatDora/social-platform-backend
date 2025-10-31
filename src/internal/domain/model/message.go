package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID             uint64         `gorm:"column:id;primaryKey"`
	ConversationID uint64         `gorm:"column:conversation_id"`
	SenderID       uint64         `gorm:"column:sender_id"`
	Type           string         `gorm:"column:type"`
	Content        string         `gorm:"column:content"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	IsRead         bool           `gorm:"column:is_read"`
	ReadAt         *time.Time     `gorm:"column:read_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`

	// relation
	Sender       *User               `gorm:"foreignKey:SenderID;references:ID"`
	Conversation *Conversation       `gorm:"foreignKey:ConversationID;references:ID"`
	Attachments  []MessageAttachment `gorm:"foreignKey:MessageID;references:ID"`
}

func (Message) TableName() string {
	return "messages"
}
