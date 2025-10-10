package response

import "social-platform-backend/internal/domain/model"

type CommunityDetailResponse struct {
	ID               uint64  `json:"id"`
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	Description      *string `json:"description,omitempty"`
	CoverImage       *string `json:"coverImage,omitempty"`
	IsPrivate        bool    `json:"isPrivate"`
}

func NewCommunityDetailResponse(community *model.Community) *CommunityDetailResponse {
	return &CommunityDetailResponse{
		ID:               community.ID,
		Name:             community.Name,
		ShortDescription: community.ShortDescription,
		Description:      community.Description,
		CoverImage:       community.CoverImage,
		IsPrivate:        community.IsPrivate,
	}
}

type CommunityListResponse struct {
	ID               uint64 `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"shortDescription"`
}

func NewCommunityListResponse(community *model.Community) *CommunityListResponse {
	return &CommunityListResponse{
		ID:               community.ID,
		Name:             community.Name,
		ShortDescription: community.ShortDescription,
	}
}
