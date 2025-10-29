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
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
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

// user config, return after login
type UserConfigResponse struct {
	Username string  `json:"username"`
	Avatar   *string `json:"avatar,omitempty"`
	// List of communities where the user is a moderator
	ModeratedCommunities []CommunityModerator `json:"moderatedCommunities,omitempty"`

	// more config fields can be added later
}

type CommunityModerator struct {
	CommunityID uint64 `json:"communityId"`
	Role        string `json:"role"`
}
