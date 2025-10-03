package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type UserVerificationRepositoryImpl struct {
	db *gorm.DB
}

func NewUserVerificationRepository(db *gorm.DB) repository.UserVerificationRepository {
	return &UserVerificationRepositoryImpl{db: db}
}

func (r *UserVerificationRepositoryImpl) CreateVerification(verification *model.UserVerification) error {
	return r.db.Create(verification).Error
}

func (r *UserVerificationRepositoryImpl) GetVerificationByToken(token string) (*model.UserVerification, error) {
	var verification model.UserVerification
	err := r.db.Where("token = ?", token).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *UserVerificationRepositoryImpl) DeleteVerification(id uint64) error {
	return r.db.Delete(&model.UserVerification{}, id).Error
}
