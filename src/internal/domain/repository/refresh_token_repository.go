package repository

import "social-platform-backend/internal/domain/model"

type RefreshTokenRepository interface {
	CreateRefreshToken(token *model.RefreshToken) error
	GetRefreshTokenByToken(token string) (*model.RefreshToken, error)
	DeleteRefreshToken(id uint64) error
	DeleteRefreshTokenByToken(token string) error
	DeleteAllRefreshTokensByUserID(userID uint64) error
	DeleteExpiredTokens() error
}
