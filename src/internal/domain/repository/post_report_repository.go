package repository

import "social-platform-backend/internal/domain/model"

type PostReportRepository interface {
	CreatePostReport(report *model.PostReport) error
	GetPostReportsByCommunityID(communityID uint64, page, limit int) ([]*model.PostReport, int64, error)
	DeletePostReport(id uint64) error
	DeletePostReportsByPostID(postID uint64) error
	IsUserReportedPost(userID, postID uint64) (bool, error)
}
