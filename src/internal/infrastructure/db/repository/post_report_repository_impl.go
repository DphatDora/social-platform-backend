package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type PostReportRepositoryImpl struct {
	db *gorm.DB
}

func NewPostReportRepository(db *gorm.DB) repository.PostReportRepository {
	return &PostReportRepositoryImpl{db: db}
}

func (r *PostReportRepositoryImpl) CreatePostReport(report *model.PostReport) error {
	return r.db.Create(report).Error
}

func (r *PostReportRepositoryImpl) GetPostReportsByCommunityID(communityID uint64, page, limit int) ([]*model.PostReport, int64, error) {
	var reports []*model.PostReport
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Table("post_reports").
		Joins("JOIN posts ON post_reports.post_id = posts.id").
		Where("posts.community_id = ?", communityID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get reports with post and reporter info, ordered by newest first
	err := r.db.Table("post_reports").
		Select("post_reports.*").
		Joins("JOIN posts ON post_reports.post_id = posts.id").
		Where("posts.community_id = ?", communityID).
		Order("post_reports.created_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("Post").
		Preload("Post.Author").
		Preload("Reporter").
		Find(&reports).Error

	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *PostReportRepositoryImpl) DeletePostReport(id uint64) error {
	return r.db.Where("id = ?", id).Delete(&model.PostReport{}).Error
}

func (r *PostReportRepositoryImpl) IsUserReportedPost(userID, postID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.PostReport{}).
		Where("reporter_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
