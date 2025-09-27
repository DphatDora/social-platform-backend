package model

import "time"

type Conversation struct {
	ID      uint64 `gorm:"column:id;primaryKey"`
	User1ID uint64 `gorm:"column:user1_id"`
	User2ID uint64 `gorm:"column:user2_id"`

	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// relations
	User1 *User `gorm:"foreignKey:User1ID;references:ID"`
	User2 *User `gorm:"foreignKey:User2ID;references:ID"`
}

func (Conversation) TableName() string {
	return "conversations"
}
