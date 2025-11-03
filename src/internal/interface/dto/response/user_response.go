package response

import (
	"social-platform-backend/internal/domain/model"
	"time"
)

type UserProfileResponse struct {
	ID              uint64          `json:"id"`
	Username        string          `json:"username"`
	Email           string          `json:"email"`
	Bio             *string         `json:"bio,omitempty"`
	Avatar          *string         `json:"avatar,omitempty"`
	UserAchievement UserAchievement `json:"achievement"`
	CreatedAt       time.Time       `json:"createdAt"`
}

func NewUserProfileResponse(user *model.User, achievement UserAchievement) *UserProfileResponse {
	return &UserProfileResponse{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		Bio:             user.Bio,
		Avatar:          user.Avatar,
		UserAchievement: achievement,
		CreatedAt:       user.CreatedAt,
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

type UserAchievement struct {
	Karma         uint64 `json:"karma"`
	Badge         string `json:"badge"`
	TotalPosts    uint64 `json:"totalPosts"`
	TotalComments uint64 `json:"totalComments"`
}
