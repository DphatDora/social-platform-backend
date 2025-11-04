package request

type UpdateNotificationSettingRequest struct {
	Action     string `json:"action" binding:"required"`
	IsPush     *bool  `json:"isPush"`
	IsSendMail *bool  `json:"isSendMail"`
}
