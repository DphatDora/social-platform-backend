package model

import "time"

type UserBadge struct {
	UserID    uint64    `gorm:"column:user_id;primaryKey"`
	BadgeID   uint64    `gorm:"column:badge_id;primaryKey"`
	AwardedAt time.Time `gorm:"column:awarded_at"`

	// relation
	User  *User  `gorm:"foreignKey:UserID;references:ID"`
	Badge *Badge `gorm:"foreignKey:BadgeID;references:ID"`
}

func (UserBadge) TableName() string {
	return "user_badges"
}
