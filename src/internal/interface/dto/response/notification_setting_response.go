package response

import (
	"social-platform-backend/internal/domain/model"
)

type NotificationSettingResponse struct {
	ID         uint64 `json:"id"`
	Action     string `json:"action"`
	IsPush     bool   `json:"isPush"`
	IsSendMail bool   `json:"isSendMail"`
}

func NewNotificationSettingResponse(setting *model.NotificationSetting) *NotificationSettingResponse {
	return &NotificationSettingResponse{
		ID:         setting.ID,
		Action:     setting.Action,
		IsPush:     setting.IsPush,
		IsSendMail: setting.IsSendMail,
	}
}
