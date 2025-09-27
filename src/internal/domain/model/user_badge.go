package model

import "time"

type UserBadge struct {
	UserID    uint64    `gorm:"column:user_id"`
	User      *User     `gorm:"foreignKey:UserID;references:ID"`
	BadgeID   uint64    `gorm:"column:badge_id"`
	AwardedAt time.Time `gorm:"column:awarded_at;autoCreateTime"`

	// relation
	Badge *Badge `gorm:"foreignKey:BadgeID;references:ID"`
}

func (UserBadge) TableName() string {
	return "user_badges"
}
