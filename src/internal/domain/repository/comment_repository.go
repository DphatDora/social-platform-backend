package repository

import "social-platform-backend/internal/domain/model"

type CommentRepository interface {
	CreateComment(comment *model.Comment) error
	GetCommentByID(id uint64) (*model.Comment, error)
	GetCommentsByPostID(postID uint64, limit, offset int) ([]*model.Comment, int64, error)
	UpdateComment(id uint64, content string, mediaURL *string) error
	DeleteComment(commentID uint64, parentCommentID *uint64) error
	GetRepliesByParentID(parentID uint64) ([]*model.Comment, error)
	GetCommentsByUserID(userID uint64, sortBy string, page, limit int) ([]*model.Comment, int64, error)
}
