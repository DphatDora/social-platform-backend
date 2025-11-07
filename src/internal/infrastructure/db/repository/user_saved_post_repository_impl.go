package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/package/util"

	"gorm.io/gorm"
)

type UserSavedPostRepositoryImpl struct {
	db *gorm.DB
}

func NewUserSavedPostRepository(db *gorm.DB) repository.UserSavedPostRepository {
	return &UserSavedPostRepositoryImpl{db: db}
}

func (r *UserSavedPostRepositoryImpl) GetUserSavedPosts(userID uint64, searchTitle string, isFollowed *bool, page, limit int) ([]*model.UserSavedPost, int64, error) {
	var savedPosts []*model.UserSavedPost
	var total int64

	countQuery := r.db.Model(&model.UserSavedPost{}).Where("user_id = ?", userID)
	query := r.db.Where("user_id = ?", userID)

	if searchTitle != "" {
		patterns := util.BuildSearchPattern(searchTitle)
		for _, p := range patterns {
			countQuery = countQuery.Where("unaccent(lower(post_title)) LIKE unaccent(lower(?))", p)
			query = query.Where("unaccent(lower(post_title)) LIKE unaccent(lower(?))", p)
		}
	}

	if isFollowed != nil {
		countQuery = countQuery.Where("is_followed = ?", *isFollowed)
		query = query.Where("is_followed = ?", *isFollowed)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC")

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Preload("Community").Find(&savedPosts).Error; err != nil {
		return nil, 0, err
	}

	return savedPosts, total, nil
}

func (r *UserSavedPostRepositoryImpl) CreateUserSavedPost(userID uint64, savedPost *request.UserSavedPostRequest) error {
	var post model.Post
	if err := r.db.Select("id, title, created_at, author_id, community_id").Where("id = ?", savedPost.PostID).First(&post).Error; err != nil {
		return err
	}

	var author model.User
	if err := r.db.Select("id, username, avatar").Where("id = ?", post.AuthorID).First(&author).Error; err != nil {
		return err
	}

	userSavedPost := &model.UserSavedPost{
		UserID:        userID,
		PostID:        savedPost.PostID,
		PostTitle:     post.Title,
		PostCreatedAt: post.CreatedAt,
		AuthorID:      author.ID,
		AuthorName:    author.Username,
		AuthorAvatar:  author.Avatar,
		CommunityID:   post.CommunityID,
		IsFollowed:    savedPost.IsFollowed,
	}

	return r.db.Create(userSavedPost).Error
}

func (r *UserSavedPostRepositoryImpl) UpdateFollowedStatus(userID, postID uint64, isFollowed bool) error {
	return r.db.Model(&model.UserSavedPost{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Update("is_followed", isFollowed).Error
}

func (r *UserSavedPostRepositoryImpl) DeleteUserSavedPost(userID, postID uint64) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.UserSavedPost{}).Error
}

func (r *UserSavedPostRepositoryImpl) CheckUserSavedPostExists(userID, postID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.UserSavedPost{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}

func (r *UserSavedPostRepositoryImpl) GetFollowersByPostID(postID uint64) ([]uint64, error) {
	var userIDs []uint64
	err := r.db.Model(&model.UserSavedPost{}).
		Where("post_id = ? AND is_followed = ?", postID, true).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
