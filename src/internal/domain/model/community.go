package model

import (
	"time"

	"gorm.io/gorm"
)

type Community struct {
	ID               uint64         `gorm:"column:id;primaryKey"`
	Name             string         `gorm:"column:name"`
	ShortDescription string         `gorm:"column:short_description"`
	Description      *string        `gorm:"column:description"`
	CoverImage       *string        `gorm:"column:cover_image"`
	CreatedAt        time.Time      `gorm:"column:created_at"`
	CreatedBy        uint64         `gorm:"column:created_by"`
	IsPrivate        bool           `gorm:"column:is_private"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at"`

	// computed column
	MemberCount int64 `gorm:"column:member_count;<-:false"`

	// relation
	Creator *User `gorm:"foreignKey:CreatedBy;references:ID"`
}

func (Community) TableName() string {
	return "communities"
}
