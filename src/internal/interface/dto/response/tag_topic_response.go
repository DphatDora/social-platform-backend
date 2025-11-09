package response

type TagResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type TopicResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
