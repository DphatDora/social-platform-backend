package response

type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Total   int64  `json:"total"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	NextURL string `json:"nextUrl,omitempty"`
}
