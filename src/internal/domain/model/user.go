package model

import "time"

type User struct {
	ID                uint64     `gorm:"column:id;primaryKey"`
	Username          string     `gorm:"column:username"`
	Email             string     `gorm:"column:email"`
	Password          string     `gorm:"column:password"`
	Karma             uint64     `gorm:"column:karma"`
	Bio               *string    `gorm:"column:bio"`
	Avatar            *string    `gorm:"column:avatar"`
	IsActive          bool       `gorm:"column:is_active"`
	Role              string     `gorm:"column:role"`
	PasswordChangedAt *time.Time `gorm:"column:password_changed_at"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         *time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
