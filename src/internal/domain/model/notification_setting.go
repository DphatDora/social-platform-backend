package model

import "time"

type NotificationSetting struct {
	ID         uint64    `gorm:"column:id;primaryKey"`
	UserID     uint64    `gorm:"column:user_id"`
	Action     string    `gorm:"column:action"`
	IsPush     bool      `gorm:"column:is_push"`
	IsSendMail bool      `gorm:"column:is_send_mail"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// relations
	User *User `gorm:"foreignKey:UserID;references:ID"`
}

func (NotificationSetting) TableName() string {
	return "notification_settings"
}
