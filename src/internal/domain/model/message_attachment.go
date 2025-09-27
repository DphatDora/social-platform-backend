package model

import "time"

type MessageAttachment struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	MessageID uint64    `gorm:"column:message_id"`
	FileURL   string    `gorm:"column:file_url"`
	FileType  string    `gorm:"column:file_type"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	// relation
	Message *Message `gorm:"foreignKey:MessageID;references:ID"`
}

func (MessageAttachment) TableName() string {
	return "message_attachments"
}
