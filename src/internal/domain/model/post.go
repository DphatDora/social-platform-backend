package model

import "time"

type Post struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	CommunityID uint64     `gorm:"column:community_id"`
	AuthorID    uint64     `gorm:"column:author_id"`
	Title       string     `gorm:"column:title"`
	Content     string     `gorm:"column:content"`
	URL         *string    `gorm:"column:url"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`

	// relation
	Community *Community `gorm:"foreignKey:CommunityID;references:ID"`
	Author    *User      `gorm:"foreignKey:AuthorID;references:ID"`
	Comments  []*Comment `gorm:"foreignKey:PostID"`
}

func (Post) TableName() string {
	return "posts"
}
