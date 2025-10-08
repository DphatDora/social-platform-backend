package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"time"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) IsEmailExisted(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) GetUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) ActivateUser(id uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("is_active", true).Error
}

func (r *UserRepositoryImpl) UpdatePasswordAndSetChangedAt(id uint64, hashedPassword string) error {
	now := time.Now()
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password":            hashedPassword,
		"password_changed_at": now,
	}).Error
}

func (r *UserRepositoryImpl) UpdateUserProfile(id uint64, updateUser *request.UpdateUserProfileRequest) error {
	updates := make(map[string]interface{})
	if updateUser.Username != nil {
		updates["username"] = *updateUser.Username
	}
	if updateUser.Bio != nil {
		updates["bio"] = *updateUser.Bio
	}
	if updateUser.Avatar != nil {
		updates["avatar"] = *updateUser.Avatar
	}
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}
