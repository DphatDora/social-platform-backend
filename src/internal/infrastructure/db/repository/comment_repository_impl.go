package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/package/constant"

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

func (r *CommentRepositoryImpl) GetCommentsByPostID(postID uint64, sortBy string, limit, offset int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	// Get total count of top-level comments (parent_comment_id IS NULL)
	if err := r.db.Model(&model.Comment{}).
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orderClause string
	switch sortBy {
	case constant.COMMENT_SORT_OLDEST:
		orderClause = "comments.created_at ASC"
	case constant.COMMENT_SORT_POPULAR:
		orderClause = "vote DESC, comments.created_at DESC"
	default: // newest
		orderClause = "comments.created_at DESC"
	}

	// Get top-level comments with author info and vote count
	err := r.db.Table("comments").
		Select(`comments.*,
			COALESCE(SUM(CASE WHEN comment_votes.vote = true THEN 1 WHEN comment_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`).
		Joins("LEFT JOIN comment_votes ON comments.id = comment_votes.comment_id").
		Where("comments.post_id = ? AND comments.parent_comment_id IS NULL", postID).
		Group("comments.id").
		Order(orderClause).
		Preload("Author").
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

	// Get replies with vote count
	err := r.db.Table("comments").
		Select(`comments.*,
			COALESCE(SUM(CASE WHEN comment_votes.vote = true THEN 1 WHEN comment_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`).
		Joins("LEFT JOIN comment_votes ON comments.id = comment_votes.comment_id").
		Where("comments.parent_comment_id = ?", parentID).
		Group("comments.id").
		Order("comments.created_at ASC").
		Preload("Author").
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

func (r *CommentRepositoryImpl) GetCommentsByUserID(userID uint64, sortBy string, page, limit int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	// Count total comments by user
	if err := r.db.Model(&model.Comment{}).
		Where("author_id = ? AND deleted_at IS NULL", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Table("comments").
		Select(`comments.*,
			COALESCE(SUM(CASE WHEN comment_votes.vote = true THEN 1 WHEN comment_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`).
		Joins("LEFT JOIN comment_votes ON comments.id = comment_votes.comment_id").
		Where("comments.author_id = ? AND comments.deleted_at IS NULL", userID).
		Group("comments.id").
		Preload("Author").
		Preload("Post")

	switch sortBy {
	case constant.SORT_TOP:
		query = query.Order("vote DESC, comments.created_at DESC")
	case constant.SORT_NEW:
		fallthrough
	default:
		query = query.Order("comments.created_at DESC")
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
