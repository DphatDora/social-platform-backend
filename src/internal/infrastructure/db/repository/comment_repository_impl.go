package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type CommentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repository.CommentRepository {
	return &CommentRepositoryImpl{db: db}
}

func (r *CommentRepositoryImpl) CreateComment(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepositoryImpl) GetCommentByID(id uint64) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepositoryImpl) GetCommentsByPostID(postID uint64, limit, offset int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	// Get total count of top-level comments (parent_comment_id IS NULL)
	if err := r.db.Model(&model.Comment{}).
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get top-level comments with author info
	err := r.db.Preload("Author").
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *CommentRepositoryImpl) GetRepliesByParentID(parentID uint64) ([]*model.Comment, error) {
	var replies []*model.Comment
	err := r.db.Preload("Author").
		Where("parent_comment_id = ?", parentID).
		Order("created_at ASC").
		Find(&replies).Error
	if err != nil {
		return nil, err
	}
	return replies, nil
}

func (r *CommentRepositoryImpl) UpdateComment(id uint64, content string, mediaURL *string) error {
	updates := map[string]interface{}{
		"content":   content,
		"media_url": mediaURL,
	}
	return r.db.Model(&model.Comment{}).Where("id = ?", id).Updates(updates).Error
}

func (r *CommentRepositoryImpl) DeleteComment(commentID uint64, parentCommentID *uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update parent_comment_id of all direct replies (promotes 1 level up)
		updates := map[string]interface{}{
			"parent_comment_id": parentCommentID,
		}
		if err := tx.Model(&model.Comment{}).
			Where("parent_comment_id = ?", commentID).
			Updates(updates).Error; err != nil {
			return err
		}

		// Delete the comment
		if err := tx.Delete(&model.Comment{}, commentID).Error; err != nil {
			return err
		}

		return nil
	})
}
