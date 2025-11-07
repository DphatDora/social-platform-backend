package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"
)

type UserRepository interface {
	IsEmailExisted(email string) (bool, error)
	CreateUser(user *model.User) error
	GetUserByID(id uint64) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByGoogleID(googleID string) (*model.User, error)
	ActivateUser(id uint64) error
	UpdatePasswordAndSetChangedAt(id uint64, hashedPassword string) error
	UpdateUserProfile(id uint64, updateUser *request.UpdateUserProfileRequest) error
	UpdateAuthProvider(userID uint64, provider string) error
	LinkGoogleAccount(userID uint64, googleID string, provider string) error
	GetLatestUserBadge(userID uint64) (*model.UserBadge, error)
	GetUserPostCount(userID uint64) (uint64, error)
	GetUserCommentCount(userID uint64) (uint64, error)
	GetUserBadgeHistory(userID uint64) ([]*model.UserBadge, error)
	SearchUsers(searchTerm string, page, limit int) ([]*model.User, int64, error)
}
