package model

import "time"

type Comment struct {
	ID              uint64     `gorm:"column:id;primaryKey"`
	PostID          uint64     `gorm:"column:post_id"`
	AuthorID        uint64     `gorm:"column:author_id"`
	ParentCommentID *uint64    `gorm:"column:parent_comment_id"`
	Content         string     `gorm:"column:content"`
	MediaURL        *string    `gorm:"column:media_url"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Total vote
	Vote int64 `gorm:"column:vote;<-:false"`
	// User's vote status (1=upvote, 0=downvote, NULL=no vote)
	UserVote *int `gorm:"column:user_vote;<-:false"`

	// relation
	Post          *Post      `gorm:"foreignKey:PostID;references:ID"`
	Author        *User      `gorm:"foreignKey:AuthorID;references:ID"`
	ParentComment *Comment   `gorm:"foreignKey:ParentCommentID;references:ID"`
	ChildComments []*Comment `gorm:"foreignKey:ParentCommentID"`
}

func (Comment) TableName() string {
	return "comments"
}
