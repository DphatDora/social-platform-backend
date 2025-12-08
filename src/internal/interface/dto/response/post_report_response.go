package response

import "social-platform-backend/internal/domain/model"

type ReporterInfo struct {
	ID       uint64   `json:"id"`
	Username string   `json:"username"`
	Avatar   *string  `json:"avatar,omitempty"`
	Reasons  []string `json:"reasons"`
	Note     *string  `json:"note,omitempty"`
}

type PostReportResponse struct {
	ID           uint64         `json:"id"`
	PostID       uint64         `json:"postId"`
	PostTitle    string         `json:"postTitle"`
	Author       AuthorInfo     `json:"author"`
	Reporters    []ReporterInfo `json:"reporters"`
	TotalReports int            `json:"totalReports"`
}

func NewPostReportResponse(postID uint64, postTitle string, author *model.User) *PostReportResponse {
	return &PostReportResponse{
		PostID:    postID,
		PostTitle: postTitle,
		Author: AuthorInfo{
			ID:        author.ID,
			Username:  author.Username,
			Avatar:    author.Avatar,
			Karma:     author.Karma,
			Bio:       author.Bio,
			CreatedAt: author.CreatedAt,
		},
		Reporters:    []ReporterInfo{},
		TotalReports: 0,
	}
}

func NewReporterInfo(reporter *model.User, reasons []string, note *string) ReporterInfo {
	return ReporterInfo{
		ID:       reporter.ID,
		Username: reporter.Username,
		Avatar:   reporter.Avatar,
		Reasons:  reasons,
		Note:     note,
	}
}
