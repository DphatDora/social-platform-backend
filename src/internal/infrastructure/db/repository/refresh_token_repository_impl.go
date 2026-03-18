package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) repository.RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{db: db}
}

func (r *RefreshTokenRepositoryImpl) CreateRefreshToken(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenRepositoryImpl) GetRefreshTokenByToken(token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := r.db.Where("token = ? AND expired_at > ?", token, time.Now()).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepositoryImpl) DeleteRefreshToken(id uint64) error {
	return r.db.Delete(&model.RefreshToken{}, id).Error
}

func (r *RefreshTokenRepositoryImpl) DeleteRefreshTokenByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&model.RefreshToken{}).Error
}

func (r *RefreshTokenRepositoryImpl) DeleteAllRefreshTokensByUserID(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}

func (r *RefreshTokenRepositoryImpl) DeleteExpiredTokens() error {
	return r.db.Where("expired_at <= ?", time.Now()).Delete(&model.RefreshToken{}).Error
}
