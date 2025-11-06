package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"
)

type UserSavedPostRepository interface {
	GetUserSavedPosts(userID uint64, searchTitle string, isFollowed *bool, page, limit int) ([]*model.UserSavedPost, int64, error)
	CreateUserSavedPost(userID uint64, savedPost *request.UserSavedPostRequest) error
	UpdateFollowedStatus(userID, postID uint64, isFollowed bool) error
	DeleteUserSavedPost(userID, postID uint64) error
	CheckUserSavedPostExists(userID, postID uint64) (bool, error)
}
