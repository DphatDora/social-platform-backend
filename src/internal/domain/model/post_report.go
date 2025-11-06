package model

import (
	"time"

	"github.com/lib/pq"
)

type PostReport struct {
	ID         uint64         `gorm:"column:id;primaryKey"`
	PostID     uint64         `gorm:"column:post_id"`
	ReporterID uint64         `gorm:"column:reporter_id"`
	Reasons    pq.StringArray `gorm:"column:reasons;type:text[]"`
	Note       *string        `gorm:"column:note"`
	CreatedAt  time.Time      `gorm:"column:created_at"`

	// relations
	Post     *Post `gorm:"foreignKey:PostID"`
	Reporter *User `gorm:"foreignKey:ReporterID"`
}

func (PostReport) TableName() string {
	return "post_reports"
}
