package request

import "time"

type BanUserRequest struct {
	UserID          uint64     `json:"userId" binding:"required"`
	RestrictionType string     `json:"restrictionType" binding:"required,oneof=warning temporary_ban permanent_ban"`
	Reason          string     `json:"reason" binding:"required,min=1,max=500"`
	ExpiresAt       *time.Time `json:"expiresAt"`
}
