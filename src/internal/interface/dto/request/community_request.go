package request

type CreateCommunityRequest struct {
	Name             string  `json:"name" binding:"required"`
	ShortDescription string  `json:"shortDescription" binding:"required"`
	Description      *string `json:"description,omitempty"`
	CoverImage       *string `json:"coverImage,omitempty"`
	IsPrivate        bool    `json:"isPrivate"`
}

type UpdateCommunityRequest struct {
	Name             *string `json:"name,omitempty"`
	ShortDescription *string `json:"shortDescription,omitempty"`
	Description      *string `json:"description,omitempty"`
	CoverImage       *string `json:"coverImage,omitempty"`
	IsPrivate        *bool   `json:"isPrivate,omitempty"`
}

// verify community name is unique request
type VerifyCommunityNameRequest struct {
	Name string `json:"name" binding:"required"`
}
