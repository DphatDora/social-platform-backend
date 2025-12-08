package model

import "time"

type UserRestriction struct {
	ID              uint64     `gorm:"column:id;primaryKey"`
	UserID          uint64     `gorm:"column:user_id"`
	CommunityID     uint64     `gorm:"column:community_id"`
	RestrictionType string     `gorm:"column:restriction_type"`
	Reason          string     `gorm:"column:reason"`
	IssuedBy        uint64     `gorm:"column:issued_by"`
	ExpiresAt       *time.Time `gorm:"column:expires_at"`
	CreatedAt       time.Time  `gorm:"column:created_at"`

	// relations
	User      *User      `gorm:"foreignKey:UserID"`
	Community *Community `gorm:"foreignKey:CommunityID"`
	Issuer    *User      `gorm:"foreignKey:IssuedBy"`
}

func (UserRestriction) TableName() string {
	return "user_restrictions"
}
