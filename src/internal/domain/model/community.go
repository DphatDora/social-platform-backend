package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Community struct {
	ID                     uint64         `gorm:"column:id;primaryKey"`
	Name                   string         `gorm:"column:name"`
	ShortDescription       string         `gorm:"column:short_description"`
	Description            *string        `gorm:"column:description"`
	Topic                  pq.StringArray `gorm:"column:topic;type:text[]"`
	CommunityAvatar        *string        `gorm:"column:community_avatar"`
	CoverImage             *string        `gorm:"column:cover_image"`
	CreatedAt              time.Time      `gorm:"column:created_at"`
	CreatedBy              uint64         `gorm:"column:created_by"`
	IsPrivate              bool           `gorm:"column:is_private"`
	RequiresPostApproval   bool           `gorm:"column:requires_post_approval"`
	RequiresMemberApproval bool           `gorm:"column:requires_member_approval"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deleted_at"`

	// computed column
	MemberCount  int64 `gorm:"column:member_count;<-:false"`
	IsSubscribed *bool `gorm:"column:is_subscribed;<-:false"`

	// relation
	Creator *User `gorm:"foreignKey:CreatedBy;references:ID"`
}

func (Community) TableName() string {
	return "communities"
}
