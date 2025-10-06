package response

import (
	"social-platform-backend/internal/domain/model"
	"time"
)

type UserProfileResponse struct {
	ID        uint64     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Karma     uint64     `json:"karma"`
	Bio       *string    `json:"bio,omitempty"`
	Avatar    *string    `json:"avatar,omitempty"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func NewUserProfileResponse(user *model.User) *UserProfileResponse {
	return &UserProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Karma:     user.Karma,
		Bio:       user.Bio,
		Avatar:    user.Avatar,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
