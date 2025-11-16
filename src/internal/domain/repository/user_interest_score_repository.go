package repository

import "social-platform-backend/internal/domain/model"

type UserInterestScoreRepository interface {
	UpsertInterestScore(score *model.UserInterestScore) error
	GetUserInterestScores(userID uint64, limit int) ([]*model.UserInterestScore, error)
	GetTopCommunitiesByScore(userID uint64, limit int) ([]uint64, error)
	GetUserInterestScoresWithScores(userID uint64, limit int) (map[uint64]float64, error)
	UpdateScoreByAction(userID, communityID uint64, scoreDelta float64) error
}
