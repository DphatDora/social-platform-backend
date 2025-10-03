package model

import "time"

type CommentVote struct {
	UserID    uint64    `gorm:"column:user_id;primaryKey"`
	CommentID uint64    `gorm:"column:comment_id;primaryKey"`
	Vote      bool      `gorm:"column:vote"`
	VotedAt   time.Time `gorm:"column:voted_at"`

	// relation
	User    *User    `gorm:"foreignKey:UserID;references:ID"`
	Comment *Comment `gorm:"foreignKey:CommentID;references:ID"`
}

func (CommentVote) TableName() string {
	return "comment_votes"
}
