package model

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	ID          uint64           `gorm:"column:id;primaryKey"`
	CommunityID uint64           `gorm:"column:community_id"`
	AuthorID    uint64           `gorm:"column:author_id"`
	Title       string           `gorm:"column:title"`
	Type        string           `gorm:"column:type"`
	Content     string           `gorm:"column:content"`
	URL         *string          `gorm:"column:url"`
	MediaURLs   *pq.StringArray  `gorm:"column:media_urls;type:text[]"`
	PollData    *json.RawMessage `gorm:"column:poll_data"`
	Tags        *pq.StringArray  `gorm:"column:tags;type:text[]"`
	Status      string           `gorm:"column:status;default:'pending'"`
	CreatedAt   time.Time        `gorm:"column:created_at"`
	UpdatedAt   *time.Time       `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt   `gorm:"column:deleted_at"`

	// Total vote
	Vote int64 `gorm:"column:vote;<-:false"`

	// User's vote status (1=upvote, 0=downvote, NULL=no vote)
	UserVote *int `gorm:"column:user_vote;<-:false"`

	// relation
	Community *Community `gorm:"foreignKey:CommunityID;references:ID"`
	Author    *User      `gorm:"foreignKey:AuthorID;references:ID"`
	Comments  []*Comment `gorm:"foreignKey:PostID"`
}

func (Post) TableName() string {
	return "posts"
}
