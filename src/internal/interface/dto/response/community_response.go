package response

import (
	"social-platform-backend/internal/domain/model"
	"time"

	"github.com/lib/pq"
)

type CommunityDetailResponse struct {
	ID               uint64         `json:"id"`
	Name             string         `json:"name"`
	ShortDescription string         `json:"shortDescription"`
	Description      *string        `json:"description,omitempty"`
	Topic            pq.StringArray `json:"topic,omitempty"`
	CommunityAvatar  *string        `json:"communityAvatar,omitempty"`
	CoverImage       *string        `json:"coverImage,omitempty"`
	IsPrivate        bool           `json:"isPrivate"`
	CreatedAt        time.Time      `json:"createdAt"`
	TotalMembers     int64          `json:"totalMembers"`

	// List of moderators
	Moderators []ModeratorResponse `json:"moderators,omitempty"`
}

func NewCommunityDetailResponse(community *model.Community) *CommunityDetailResponse {
	return &CommunityDetailResponse{
		ID:               community.ID,
		Name:             community.Name,
		ShortDescription: community.ShortDescription,
		Description:      community.Description,
		Topic:            community.Topic,
		CommunityAvatar:  community.CommunityAvatar,
		CoverImage:       community.CoverImage,
		IsPrivate:        community.IsPrivate,
		CreatedAt:        community.CreatedAt,
	}
}

type CommunityListResponse struct {
	ID               uint64         `json:"id"`
	Name             string         `json:"name"`
	ShortDescription string         `json:"shortDescription"`
	Topic            pq.StringArray `json:"topic,omitempty"`
	CommunityAvatar  *string        `json:"communityAvatar,omitempty"`
	IsPrivate        bool           `json:"isPrivate"`
	TotalMembers     int64          `json:"totalMembers"`
	IsFollow         *bool          `json:"isFollow,omitempty"`
}

func NewCommunityListResponse(community *model.Community) *CommunityListResponse {
	return &CommunityListResponse{
		ID:               community.ID,
		Name:             community.Name,
		ShortDescription: community.ShortDescription,
		Topic:            community.Topic,
		CommunityAvatar:  community.CommunityAvatar,
		IsPrivate:        community.IsPrivate,
	}
}

type MemberListResponse struct {
	UserID       uint64    `json:"userId"`
	Username     string    `json:"username"`
	Avatar       *string   `json:"avatar,omitempty"`
	Karma        uint64    `json:"karma"`
	Role         string    `json:"role"`
	SubscribedAt time.Time `json:"subscribedAt"`
}

func NewMemberListResponse(user *model.User, subscribedAt time.Time, role string) *MemberListResponse {
	return &MemberListResponse{
		UserID:       user.ID,
		Username:     user.Username,
		Avatar:       user.Avatar,
		Karma:        user.Karma,
		Role:         role,
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
