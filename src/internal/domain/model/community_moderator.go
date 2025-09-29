package model

import "time"

type CommunityModerator struct {
	CommunityID uint64    `gorm:"column:community_id;primaryKey"`
	UserID      uint64    `gorm:"column:user_id;primaryKey"`
	Role        string    `gorm:"column:role"`
	JoinedAt    time.Time `gorm:"column:joined_at"`

	// relation
	Community *Community `gorm:"foreignKey:CommunityID;references:ID"`
	User      *User      `gorm:"foreignKey:UserID;references:ID"`
}

func (CommunityModerator) TableName() string {
	return "community_moderators"
}
