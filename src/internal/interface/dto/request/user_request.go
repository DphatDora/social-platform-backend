package request

type UpdateUserProfileRequest struct {
	Username *string `json:"username,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}
