package model

import "time"

type UserInterestScore struct {
	UserID      uint64     `gorm:"column:user_id;primaryKey"`
	CommunityID uint64     `gorm:"column:community_id;primaryKey"`
	Score       float64    `gorm:"column:score;default:0;index"`
	LastVoteAt  *time.Time `gorm:"column:last_vote_at"`
	LastJoinAt  *time.Time `gorm:"column:last_join_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`

	// Relations
	User      *User      `gorm:"foreignKey:UserID;references:ID"`
	Community *Community `gorm:"foreignKey:CommunityID;references:ID"`
}

func (UserInterestScore) TableName() string {
	return "user_interest_scores"
}
