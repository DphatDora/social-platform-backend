package response

import (
	"social-platform-backend/internal/domain/model"
	"time"
)

type CommunityDetailResponse struct {
	ID               uint64    `json:"id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"shortDescription"`
	Description      *string   `json:"description,omitempty"`
	CoverImage       *string   `json:"coverImage,omitempty"`
	IsPrivate        bool      `json:"isPrivate"`
	CreatedAt        time.Time `json:"createdAt"`
	TotalMembers     int64     `json:"totalMembers"`

	// List of moderators
	Moderators []ModeratorResponse `json:"moderators,omitempty"`
}

func NewCommunityDetailResponse(community *model.Community) *CommunityDetailResponse {
	return &CommunityDetailResponse{
		ID:               community.ID,
		Name:             community.Name,
		ShortDescription: community.ShortDescription,
		Description:      community.Description,
		CoverImage:       community.CoverImage,
		IsPrivate:        community.IsPrivate,
		CreatedAt:        community.CreatedAt,
	}
}

type CommunityListResponse struct {
	ID               uint64 `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"shortDescription"`
	IsPrivate        bool   `json:"isPrivate"`
	TotalMembers     int64  `json:"totalMembers"`
}

func NewCommunityListResponse(community *model.Community) *CommunityListResponse {
	return &CommunityListResponse{
		ID:               community.ID,
		Name:             community.Name,
		ShortDescription: community.ShortDescription,
		IsPrivate:        community.IsPrivate,
	}
}

type MemberListResponse struct {
	UserID       uint64    `json:"userId"`
	Username     string    `json:"username"`
	Avatar       *string   `json:"avatar,omitempty"`
	Karma        uint64    `json:"karma"`
	SubscribedAt time.Time `json:"subscribedAt"`
}

func NewMemberListResponse(user *model.User, subscribedAt time.Time) *MemberListResponse {
	return &MemberListResponse{
		UserID:       user.ID,
		Username:     user.Username,
		Avatar:       user.Avatar,
		Karma:        user.Karma,
		SubscribedAt: subscribedAt,
	}
}

// Moderator of community response
type ModeratorResponse struct {
	UserID   uint64  `json:"userId"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar,omitempty"`
	Role     string  `json:"role"`
}

func NewModeratorResponse(user *model.User, role string) *ModeratorResponse {
	return &ModeratorResponse{
		UserID:   user.ID,
		Username: user.Username,
		Avatar:   user.Avatar,
		Role:     role,
	}
}

// verify community name is unique response
type VerifyCommunityNameResponse struct {
	IsUnique bool `json:"isUnique"`
}
