package request

import "github.com/lib/pq"

type CreateCommunityRequest struct {
	Name             string         `json:"name" binding:"required"`
	ShortDescription string         `json:"shortDescription" binding:"required"`
	Description      *string        `json:"description,omitempty"`
	Topic            pq.StringArray `json:"topic,omitempty"`
	CommunityAvatar  *string        `json:"communityAvatar,omitempty"`
	CoverImage       *string        `json:"coverImage,omitempty"`
	IsPrivate        bool           `json:"isPrivate"`
}

type UpdateCommunityRequest struct {
	Name             *string         `json:"name,omitempty"`
	ShortDescription *string         `json:"shortDescription,omitempty"`
	Description      *string         `json:"description,omitempty"`
	Topic            *pq.StringArray `json:"topic,omitempty"`
	CommunityAvatar  *string         `json:"communityAvatar,omitempty"`
	CoverImage       *string         `json:"coverImage,omitempty"`
	IsPrivate        *bool           `json:"isPrivate,omitempty"`
}

// verify community name is unique request
type VerifyCommunityNameRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}

type UpdateRequiresApprovalRequest struct {
	RequiresApproval bool `json:"requiresApproval"`
}
