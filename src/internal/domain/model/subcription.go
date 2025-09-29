package model

import "time"

type Subscription struct {
	UserID       uint64    `gorm:"column:user_id;primaryKey"`
	CommunityID  uint64    `gorm:"column:community_id;primaryKey"`
	SubscribedAt time.Time `gorm:"column:subscribed_at"`

	// relation
	User      *User      `gorm:"foreignKey:UserID;references:ID"`
	Community *Community `gorm:"foreignKey:CommunityID;references:ID"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
