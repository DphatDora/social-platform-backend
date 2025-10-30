package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type CommentVoteRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentVoteRepository(db *gorm.DB) repository.CommentVoteRepository {
	return &CommentVoteRepositoryImpl{db: db}
}

func (r *CommentVoteRepositoryImpl) UpsertCommentVote(commentVote *model.CommentVote) error {
	// Check if vote already exists
	var existingVote model.CommentVote
	err := r.db.Where("user_id = ? AND comment_id = ?", commentVote.UserID, commentVote.CommentID).First(&existingVote).Error

	if err == gorm.ErrRecordNotFound {
		// Create new vote
		commentVote.VotedAt = time.Now()
		return r.db.Create(commentVote).Error
	}

	if err != nil {
		return err
	}

	// Update existing vote
	return r.db.Model(&model.CommentVote{}).
		Where("user_id = ? AND comment_id = ?", commentVote.UserID, commentVote.CommentID).
		Updates(map[string]interface{}{
			"vote":     commentVote.Vote,
			"voted_at": time.Now(),
		}).Error
}

func (r *CommentVoteRepositoryImpl) GetCommentVote(userID, commentID uint64) (*model.CommentVote, error) {
	var commentVote model.CommentVote
	err := r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&commentVote).Error
	if err != nil {
		return nil, err
	}
	return &commentVote, nil
}

func (r *CommentVoteRepositoryImpl) DeleteCommentVote(userID, commentID uint64) error {
	return r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&model.CommentVote{}).Error
}
