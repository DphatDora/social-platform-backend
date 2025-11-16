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

func (r *UserRepositoryImpl) GetUserByGoogleID(googleID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
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
	if updateUser.CoverImage != nil {
		updates["cover_image"] = *updateUser.CoverImage
	}
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepositoryImpl) UpdateAuthProvider(userID uint64, provider string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("auth_provider", provider).Error
}

func (r *UserRepositoryImpl) LinkGoogleAccount(userID uint64, googleID string, provider string) error {
	updates := map[string]interface{}{
		"google_id":     googleID,
		"auth_provider": provider,
	}
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *UserRepositoryImpl) GetLatestUserBadge(userID uint64) (*model.UserBadge, error) {
	var userBadge model.UserBadge
	err := r.db.Where("user_id = ?", userID).
		Order("awarded_at DESC").
		Preload("Badge").
		First(&userBadge).Error
	if err != nil {
		return nil, err
	}
	return &userBadge, nil
}

func (r *UserRepositoryImpl) GetUserPostCount(userID uint64) (uint64, error) {
	var count int64
	err := r.db.Model(&model.Post{}).
		Where("author_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return uint64(count), err
}

func (r *UserRepositoryImpl) GetUserCommentCount(userID uint64) (uint64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).
		Where("author_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return uint64(count), err
}

func (r *UserRepositoryImpl) GetUserBadgeHistory(userID uint64) ([]*model.UserBadge, error) {
	var userBadges []*model.UserBadge
	err := r.db.Where("user_id = ?", userID).
		Order("awarded_at DESC").
		Preload("Badge").
		Find(&userBadges).
		Limit(5).Error
	if err != nil {
		return nil, err
	}
	return userBadges, nil
}

func (r *UserRepositoryImpl) SearchUsers(searchTerm string, page, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// Build query with username search
	query := r.db.Model(&model.User{}).Where("is_active = ?", true)
	countQuery := r.db.Model(&model.User{}).Where("is_active = ?", true)

	if searchTerm != "" {
		searchPattern := "%" + searchTerm + "%"
		query = query.Where("LOWER(username) LIKE LOWER(?)", searchPattern)
		countQuery = countQuery.Where("LOWER(username) LIKE LOWER(?)", searchPattern)
	}

	// Count total
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get users with pagination
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
