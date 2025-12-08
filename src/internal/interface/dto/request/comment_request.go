package request

type CreateCommentRequest struct {
	PostID          uint64  `json:"postId" binding:"required"`
	Content         string  `json:"content" binding:"required"`
	ParentCommentID *uint64 `json:"parentCommentId,omitempty"`
	MediaURL        *string `json:"mediaUrl,omitempty"`
}

type UpdateCommentRequest struct {
	Content  string  `json:"content" binding:"required"`
	MediaURL *string `json:"mediaUrl,omitempty"`
}

type ReportCommentRequest struct {
	Reasons []string `json:"reasons" binding:"required,min=1"`
	Note    *string  `json:"note,omitempty"`
}
