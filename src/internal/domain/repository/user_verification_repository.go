package repository

import "social-platform-backend/internal/domain/model"

type UserVerificationRepository interface {
	CreateVerification(verification *model.UserVerification) error
	GetVerificationByToken(token string) (*model.UserVerification, error)
	DeleteVerification(id uint64) error
}
