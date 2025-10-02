package response

type RegisterResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

type VerifyEmailResponse struct {
	Message string `json:"message"`
}
