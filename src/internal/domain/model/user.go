package model

import "time"

type User struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	Username    string     `gorm:"column:username"`
	Email       string     `gorm:"column:email"`
	DateOfBirth *time.Time `gorm:"column:date_of_birth"`
	Gender      *string    `gorm:"column:gender"`
	Phone       *string    `gorm:"column:phone"`
	Address     *string    `gorm:"column:address"`
	Bio         *string    `gorm:"column:bio"`
	Avatar      *string    `gorm:"column:avatar"`
	CoverImage  *string    `gorm:"column:cover_image"`
	Karma       uint64     `gorm:"column:karma"`

	GoogleID     *string `gorm:"column:google_id;unique"`
	AuthProvider string  `gorm:"column:auth_provider;default:'email'"`
	IsActive     bool    `gorm:"column:is_active"`

	Password          *string    `gorm:"column:password"`
	PasswordChangedAt *time.Time `gorm:"column:password_changed_at"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         *time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
