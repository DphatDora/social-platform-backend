package response

import (
	"social-platform-backend/internal/domain/model"
	"time"
)

type UserRestrictionResponse struct {
	ID              uint64     `json:"id"`
	RestrictionType string     `json:"restrictionType"`
	Reason          string     `json:"reason"`
	ExpiresAt       *time.Time `json:"expiresAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
}

func NewUserRestrictionResponse(restriction *model.UserRestriction) *UserRestrictionResponse {
	response := &UserRestrictionResponse{
		ID:              restriction.ID,
		RestrictionType: restriction.RestrictionType,
		Reason:          restriction.Reason,
		ExpiresAt:       restriction.ExpiresAt,
		CreatedAt:       restriction.CreatedAt,
	}

	return response
}
