package repository

import "social-platform-backend/internal/domain/model"

type UserRepository interface {
	IsEmailExisted(email string) (bool, error)
	CreateUser(user *model.User) error
	GetUserByID(id uint64) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	ActivateUser(id uint64) error
	UpdatePasswordAndSetChangedAt(id uint64, hashedPassword string) error
}
