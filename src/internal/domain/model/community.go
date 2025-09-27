package model

import (
	"time"

	"gorm.io/gorm"
)

type Community struct {
	ID          uint64         `gorm:"column:id;primaryKey"`
	Name        string         `gorm:"column:name"`
	Description string         `gorm:"column:description"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	CreatedBy   uint64         `gorm:"column:created_by"`
	IsPrivate   bool           `gorm:"column:is_private"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// relation
	Creator *User `gorm:"foreignKey:CreatedBy;references:ID"`
}

func (Community) TableName() string {
	return "communities"
}
