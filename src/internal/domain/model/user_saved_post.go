package model

import "time"

type UserSavedPost struct {
	UserID        uint64    `gorm:"column:user_id;primaryKey"`
	PostID        uint64    `gorm:"column:post_id;primaryKey"`
	PostTitle     string    `gorm:"column:post_title"`
	PostCreatedAt time.Time `gorm:"column:post_created_at"`
	AuthorID      uint64    `gorm:"column:author_id"`
	AuthorName    string    `gorm:"column:author_name"`
	AuthorAvatar  *string   `gorm:"column:author_avatar"`
	IsFollowed    bool      `gorm:"column:is_followed"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (UserSavedPost) TableName() string {
	return "user_saved_posts"
}
