package model

import "time"

type Notification struct {
	ID          uint64    `gorm:"column:id;primaryKey"`
	UserID      uint64    `gorm:"column:user_id"`
	Type        string    `gorm:"column:type"`
	ReferenceID *uint64   `gorm:"column:reference_id"`
	IsRead      bool      `gorm:"column:is_read"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`

	// relations
	User *User `gorm:"foreignKey:UserID;references:ID"`
}

func (Notification) TableName() string {
	return "notifications"
}
