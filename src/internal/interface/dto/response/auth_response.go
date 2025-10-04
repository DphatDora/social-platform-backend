package response

type LoginResponse struct {
	Username    string `json:"username"`
	Avatar      string `json:"avatar,omitempty"`
	AccessToken string `json:"access_token"`
}
