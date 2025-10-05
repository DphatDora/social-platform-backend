package model

import "time"

type PasswordReset struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	UserID    uint64    `gorm:"column:user_id"`
	Token     string    `gorm:"column:token"`
	ExpiredAt time.Time `gorm:"column:expired_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}
