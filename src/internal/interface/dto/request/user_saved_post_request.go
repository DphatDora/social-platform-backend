package request

type UserSavedPostRequest struct {
	PostID     uint64 `json:"postId" binding:"required"`
	IsFollowed bool   `json:"isFollowed"`
}

type UpdateUserSavedPostRequest struct {
	IsFollowed bool `json:"isFollowed"`
}
