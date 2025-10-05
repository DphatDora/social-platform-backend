package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type PasswordResetRepositoryImpl struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) repository.PasswordResetRepository {
	return &PasswordResetRepositoryImpl{db: db}
}

func (r *PasswordResetRepositoryImpl) CreatePasswordReset(reset *model.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *PasswordResetRepositoryImpl) GetPasswordResetByToken(token string) (*model.PasswordReset, error) {
	var reset model.PasswordReset
	err := r.db.Where("token = ?", token).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *PasswordResetRepositoryImpl) DeletePasswordReset(id uint64) error {
	return r.db.Delete(&model.PasswordReset{}, id).Error
}

func (r *PasswordResetRepositoryImpl) DeletePasswordResetByUserID(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.PasswordReset{}).Error
}
