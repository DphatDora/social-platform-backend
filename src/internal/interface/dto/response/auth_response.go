package response

type LoginResponse struct {
	Username    string `json:"username"`
	Avatar      string `json:"avatar,omitempty"`
	AccessToken string `json:"access_token"`

	// List of communities where the user is a moderator
	ModeratedCommunities []CommunityModerator `json:"moderatedCommunities,omitempty"`
}

type CommunityModerator struct {
	CommunityID uint64 `json:"communityId"`
	Role        string `json:"role"`
}
