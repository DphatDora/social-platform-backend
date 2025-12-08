package response

import "social-platform-backend/internal/domain/model"

type CommentReportResponse struct {
	ID             uint64         `json:"id"`
	CommentID      uint64         `json:"commentId"`
	CommentContent string         `json:"commentContent"`
	CommentAuthor  AuthorInfo     `json:"commentAuthor"`
	PostID         uint64         `json:"postId"`
	PostTitle      string         `json:"postTitle"`
	Reporters      []ReporterInfo `json:"reporters"`
	TotalReports   int            `json:"totalReports"`
}

func NewCommentReportResponse(commentID uint64, commentContent string, commentAuthor *model.User, postID uint64, postTitle string) *CommentReportResponse {
	return &CommentReportResponse{
		CommentID:      commentID,
		CommentContent: commentContent,
		CommentAuthor: AuthorInfo{
			ID:        commentAuthor.ID,
			Username:  commentAuthor.Username,
			Avatar:    commentAuthor.Avatar,
			Karma:     commentAuthor.Karma,
			Bio:       commentAuthor.Bio,
			CreatedAt: commentAuthor.CreatedAt,
		},
		PostID:       postID,
		PostTitle:    postTitle,
		Reporters:    []ReporterInfo{},
		TotalReports: 0,
	}
}
