package model

import (
	"time"

	"gorm.io/gorm"
)

type MetaData struct {
	ID        uint64   `json:"id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Content   string   `json:"content"`
	MediaURLs []string `json:"mediaUrls"`
}

type Message struct {
	ID             uint64         `gorm:"column:id;primaryKey"`
	ConversationID uint64         `gorm:"column:conversation_id"`
	SenderID       uint64         `gorm:"column:sender_id"`
	Content        string         `gorm:"column:content"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	IsRead         bool           `gorm:"column:is_read"`
	ReadAt         *time.Time     `gorm:"column:read_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
	MetaData       *MetaData      `gorm:"column:meta_data;type:jsonb;serializer:json"`
	// relation
	Sender       *User               `gorm:"foreignKey:SenderID;references:ID"`
	Conversation *Conversation       `gorm:"foreignKey:ConversationID;references:ID"`
	Attachments  []MessageAttachment `gorm:"foreignKey:MessageID;references:ID"`
}

func (Message) TableName() string {
	return "messages"
}
