package payload

import "time"

type UpdateInterestScorePayload struct {
	UserID      uint64    `json:"user_id"`
	CommunityID uint64    `json:"community_id"`
	Action      string    `json:"action"` // upvote, downvote, follow, join
	PostID      *uint64   `json:"post_id,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}
