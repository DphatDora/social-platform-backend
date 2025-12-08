package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type CommentReportRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentReportRepository(db *gorm.DB) repository.CommentReportRepository {
	return &CommentReportRepositoryImpl{db: db}
}

func (r *CommentReportRepositoryImpl) CreateCommentReport(report *model.CommentReport) error {
	return r.db.Create(report).Error
}

func (r *CommentReportRepositoryImpl) GetCommentReportsByCommunityID(communityID uint64, page, limit int) ([]*model.CommentReport, int64, error) {
	var reports []*model.CommentReport
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Table("comment_reports").
		Joins("JOIN comments ON comment_reports.comment_id = comments.id").
		Joins("JOIN posts ON comments.post_id = posts.id").
		Where("posts.community_id = ?", communityID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Table("comment_reports").
		Select("comment_reports.*").
		Joins("JOIN comments ON comment_reports.comment_id = comments.id").
		Joins("JOIN posts ON comments.post_id = posts.id").
		Where("posts.community_id = ?", communityID).
		Order("comment_reports.created_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("Comment").
		Preload("Comment.Author").
		Preload("Reporter").
		Find(&reports).Error

	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *CommentReportRepositoryImpl) DeleteCommentReport(id uint64) error {
	var report model.CommentReport
	if err := r.db.Where("id = ?", id).First(&report).Error; err != nil {
		return err
	}

	return r.db.Where("comment_id = ?", report.CommentID).Delete(&model.CommentReport{}).Error
}

func (r *CommentReportRepositoryImpl) IsUserReportedComment(userID, commentID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.CommentReport{}).
		Where("reporter_id = ? AND comment_id = ?", userID, commentID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
