package request

import "time"

type UpdateUserProfileRequest struct {
	Username    *string    `json:"username,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Address     *string    `json:"address,omitempty"`
	Avatar      *string    `json:"avatar,omitempty"`
	CoverImage  *string    `json:"coverImage,omitempty"`
}
