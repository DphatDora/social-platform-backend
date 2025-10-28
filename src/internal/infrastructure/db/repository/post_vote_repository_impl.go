package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type PostVoteRepositoryImpl struct {
	db *gorm.DB
}

func NewPostVoteRepository(db *gorm.DB) repository.PostVoteRepository {
	return &PostVoteRepositoryImpl{db: db}
}

func (r *PostVoteRepositoryImpl) UpsertPostVote(postVote *model.PostVote) error {
	// Check if vote already exists
	var existingVote model.PostVote
	err := r.db.Where("user_id = ? AND post_id = ?", postVote.UserID, postVote.PostID).First(&existingVote).Error

	if err == gorm.ErrRecordNotFound {
		// Create new vote
		postVote.VotedAt = time.Now()
		return r.db.Create(postVote).Error
	}

	if err != nil {
		return err
	}

	// Update existing vote
	return r.db.Model(&model.PostVote{}).
		Where("user_id = ? AND post_id = ?", postVote.UserID, postVote.PostID).
		Updates(map[string]interface{}{
			"vote":     postVote.Vote,
			"voted_at": time.Now(),
		}).Error
}

func (r *PostVoteRepositoryImpl) GetPostVote(userID, postID uint64) (*model.PostVote, error) {
	var postVote model.PostVote
	err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&postVote).Error
	if err != nil {
		return nil, err
	}
	return &postVote, nil
}

func (r *PostVoteRepositoryImpl) DeletePostVote(userID, postID uint64) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.PostVote{}).Error
}
