package response

type ReporterInfo struct {
	ID       uint64   `json:"id"`
	Username string   `json:"username"`
	Avatar   *string  `json:"avatar,omitempty"`
	Reasons  []string `json:"reasons"`
	Note     *string  `json:"note,omitempty"`
}

type PostReportResponse struct {
	PostID       uint64         `json:"postId"`
	PostTitle    string         `json:"postTitle"`
	Author       AuthorInfo     `json:"author"`
	Reporters    []ReporterInfo `json:"reporters"`
	TotalReports int            `json:"totalReports"`
}
