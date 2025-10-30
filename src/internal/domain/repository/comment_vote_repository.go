package repository

import "social-platform-backend/internal/domain/model"

type CommentVoteRepository interface {
	UpsertCommentVote(commentVote *model.CommentVote) error
	GetCommentVote(userID, commentID uint64) (*model.CommentVote, error)
	DeleteCommentVote(userID, commentID uint64) error
}
