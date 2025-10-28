package repository

import "social-platform-backend/internal/domain/model"

type PostVoteRepository interface {
	UpsertPostVote(postVote *model.PostVote) error
	GetPostVote(userID, postID uint64) (*model.PostVote, error)
	DeletePostVote(userID, postID uint64) error
}
