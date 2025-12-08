package model

import (
	"time"

	"github.com/lib/pq"
)

type CommentReport struct {
	ID         uint64         `gorm:"column:id;primaryKey"`
	CommentID  uint64         `gorm:"column:comment_id"`
	ReporterID uint64         `gorm:"column:reporter_id"`
	Reasons    pq.StringArray `gorm:"column:reasons;type:text[]"`
	Note       *string        `gorm:"column:note"`
	CreatedAt  time.Time      `gorm:"column:created_at"`

	// relations
	Comment  *Comment `gorm:"foreignKey:CommentID"`
	Reporter *User    `gorm:"foreignKey:ReporterID"`
}

func (CommentReport) TableName() string {
	return "comment_reports"
}
