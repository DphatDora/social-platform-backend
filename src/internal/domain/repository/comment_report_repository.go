package repository

import "social-platform-backend/internal/domain/model"

type CommentReportRepository interface {
	CreateCommentReport(report *model.CommentReport) error
	GetCommentReportsByCommunityID(communityID uint64, page, limit int) ([]*model.CommentReport, int64, error)
	DeleteCommentReport(id uint64) error
	DeleteCommentReportsByCommentID(commentID uint64) error
	IsUserReportedComment(userID, commentID uint64) (bool, error)
}
