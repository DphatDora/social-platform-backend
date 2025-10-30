package response

import (
	"social-platform-backend/internal/domain/model"
	"time"
)

type CommentResponse struct {
	ID              uint64             `json:"id"`
	PostID          uint64             `json:"postId"`
	Author          *AuthorInfo        `json:"author,omitempty"`
	ParentCommentID *uint64            `json:"parentCommentId,omitempty"`
	Content         string             `json:"content"`
	MediaURL        *string            `json:"mediaUrl,omitempty"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       *time.Time         `json:"updatedAt,omitempty"`
	Replies         []*CommentResponse `json:"replies,omitempty"`
}

func NewCommentResponse(comment *model.Comment) *CommentResponse {
	response := &CommentResponse{
		ID:              comment.ID,
		PostID:          comment.PostID,
		ParentCommentID: comment.ParentCommentID,
		Content:         comment.Content,
		MediaURL:        comment.MediaURL,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
	}

	if comment.Author != nil {
		response.Author = &AuthorInfo{
			ID:       comment.Author.ID,
			Username: comment.Author.Username,
			Avatar:   comment.Author.Avatar,
		}
	}

	return response
}
