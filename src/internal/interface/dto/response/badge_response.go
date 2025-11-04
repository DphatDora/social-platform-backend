package response

type UserBadgeResponse struct {
	BadgeName string `json:"badgeName"`
	IconURL   string `json:"iconUrl"`
	Karma     uint64 `json:"karma"`
	MonthYear string `json:"monthYear"`
}

func NewUserBadgeResponse(badgeName, iconURL, monthYear string, karma uint64) *UserBadgeResponse {
	return &UserBadgeResponse{
		BadgeName: badgeName,
		IconURL:   iconURL,
		Karma:     karma,
		MonthYear: monthYear,
	}
}
