package model

import "time"

type UserVerification struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	UserID    uint64    `gorm:"column:user_id;not null"`
	Token     string    `gorm:"column:token;not null;unique"`
	ExpiredAt time.Time `gorm:"column:expired_at;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (UserVerification) TableName() string {
	return "user_verifications"
}
