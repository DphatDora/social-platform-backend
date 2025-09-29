package model

import "time"

type Badge struct {
	ID          uint64    `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	IconURL     string    `gorm:"column:icon_url"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (Badge) TableName() string {
	return "badges"
}
