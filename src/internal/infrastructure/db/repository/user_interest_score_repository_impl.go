package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type UserInterestScoreRepositoryImpl struct {
	db *gorm.DB
}

func NewUserInterestScoreRepository(db *gorm.DB) repository.UserInterestScoreRepository {
	return &UserInterestScoreRepositoryImpl{db: db}
}

func (r *UserInterestScoreRepositoryImpl) UpsertInterestScore(score *model.UserInterestScore) error {
	return r.db.Save(score).Error
}

func (r *UserInterestScoreRepositoryImpl) GetUserInterestScores(userID uint64, limit int) ([]*model.UserInterestScore, error) {
	var scores []*model.UserInterestScore
	err := r.db.Where("user_id = ?", userID).
		Order("score DESC").
		Limit(limit).
		Preload("Community").
		Find(&scores).Error
	return scores, err
}

func (r *UserInterestScoreRepositoryImpl) GetTopCommunitiesByScore(userID uint64, limit int) ([]uint64, error) {
	var communityIDs []uint64
	err := r.db.Model(&model.UserInterestScore{}).
		Where("user_id = ? AND score > 0", userID).
		Order("score DESC").
		Limit(limit).
		Pluck("community_id", &communityIDs).Error
	return communityIDs, err
}

func (r *UserInterestScoreRepositoryImpl) GetUserInterestScoresWithScores(userID uint64, limit int) (map[uint64]float64, error) {
	var scores []*model.UserInterestScore
	err := r.db.Where("user_id = ? AND score > 0", userID).
		Order("score DESC").
		Limit(limit).
		Find(&scores).Error

	if err != nil {
		return nil, err
	}

	scoresMap := make(map[uint64]float64)
	for _, score := range scores {
		scoresMap[score.CommunityID] = score.Score
	}

	return scoresMap, nil
}

func (r *UserInterestScoreRepositoryImpl) UpdateScoreByAction(userID, communityID uint64, scoreDelta float64) error {
	now := time.Now()

	// Try to find existing record
	var existingScore model.UserInterestScore
	err := r.db.Where("user_id = ? AND community_id = ?", userID, communityID).
		First(&existingScore).Error

	if err == gorm.ErrRecordNotFound {
		// Create new record
		newScore := &model.UserInterestScore{
			UserID:      userID,
			CommunityID: communityID,
			Score:       scoreDelta,
			UpdatedAt:   now,
		}
		return r.db.Create(newScore).Error
	} else if err != nil {
		return err
	}

	// Update existing record
	return r.db.Model(&model.UserInterestScore{}).
		Where("user_id = ? AND community_id = ?", userID, communityID).
		Updates(map[string]interface{}{
			"score":      gorm.Expr("score + ?", scoreDelta),
			"updated_at": now,
		}).Error
}
