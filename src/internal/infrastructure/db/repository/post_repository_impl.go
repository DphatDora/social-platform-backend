package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"time"

	"gorm.io/gorm"
)

type PostRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) CreatePost(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepositoryImpl) GetPostByID(id uint64) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepositoryImpl) UpdatePostText(id uint64, updatePost *request.UpdatePostTextRequest) error {
	updates := make(map[string]interface{})
	if updatePost.Title != nil {
		updates["title"] = *updatePost.Title
	}
	if updatePost.Content != nil {
		updates["content"] = *updatePost.Content
	}
	if updatePost.Tags != nil {
		updates["tags"] = *updatePost.Tags
	}
	if len(updates) > 0 {
		now := time.Now()
		updates["updated_at"] = now
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepositoryImpl) UpdatePostLink(id uint64, updatePost *request.UpdatePostLinkRequest) error {
	updates := make(map[string]interface{})
	if updatePost.Title != nil {
		updates["title"] = *updatePost.Title
	}
	if updatePost.Content != nil {
		updates["content"] = *updatePost.Content
	}
	if updatePost.URL != nil {
		updates["url"] = *updatePost.URL
	}
	if updatePost.Tags != nil {
		updates["tags"] = *updatePost.Tags
	}
	if len(updates) > 0 {
		now := time.Now()
		updates["updated_at"] = now
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepositoryImpl) UpdatePostMedia(id uint64, updatePost *request.UpdatePostMediaRequest) error {
	updates := make(map[string]interface{})
	if updatePost.Title != nil {
		updates["title"] = *updatePost.Title
	}
	if updatePost.Content != nil {
		updates["content"] = *updatePost.Content
	}
	if updatePost.MediaURLs != nil {
		updates["media_urls"] = *updatePost.MediaURLs
	}
	if updatePost.Tags != nil {
		updates["tags"] = *updatePost.Tags
	}
	if len(updates) > 0 {
		now := time.Now()
		updates["updated_at"] = now
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepositoryImpl) UpdatePostPoll(id uint64, updatePost *request.UpdatePostPollRequest) error {
	updates := make(map[string]interface{})
	if updatePost.Title != nil {
		updates["title"] = *updatePost.Title
	}
	if updatePost.Content != nil {
		updates["content"] = *updatePost.Content
	}
	if updatePost.PollData != nil {
		updates["poll_data"] = *updatePost.PollData
	}
	if updatePost.Tags != nil {
		updates["tags"] = *updatePost.Tags
	}
	if len(updates) > 0 {
		now := time.Now()
		updates["updated_at"] = now
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepositoryImpl) DeletePost(id uint64) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
