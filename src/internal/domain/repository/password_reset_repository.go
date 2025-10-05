package repository

import "social-platform-backend/internal/domain/model"

type PasswordResetRepository interface {
	CreatePasswordReset(reset *model.PasswordReset) error
	GetPasswordResetByToken(token string) (*model.PasswordReset, error)
	DeletePasswordReset(id uint64) error
	DeletePasswordResetByUserID(userID uint64) error
}
