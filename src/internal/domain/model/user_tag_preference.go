package model

import (
	"time"

	"github.com/lib/pq"
)

type UserTagPreference struct {
	UserID        uint64         `gorm:"column:user_id;primaryKey"`
	PreferredTags pq.StringArray `gorm:"column:preferred_tags;type:text[]"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`

	// Relation
	User *User `gorm:"foreignKey:UserID;references:ID"`
}

func (UserTagPreference) TableName() string {
	return "user_tag_preferences"
}
