package model

import "time"

type PostVote struct {
	UserID  uint64    `gorm:"column:user_id;primaryKey"`
	PostID  uint64    `gorm:"column:post_id;primaryKey"`
	Vote    bool      `gorm:"column:vote"`
	VotedAt time.Time `gorm:"column:voted_at"`

	// relation
	User *User `gorm:"foreignKey:UserID;references:ID"`
	Post *Post `gorm:"foreignKey:PostID;references:ID"`
}

func (PostVote) TableName() string {
	return "post_votes"
}
