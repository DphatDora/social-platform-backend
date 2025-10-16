package payload

import "time"

type PollOption struct {
	ID     int      `json:"id"`
	Text   string   `json:"text"`
	Votes  int      `json:"votes"`
	Voters []uint64 `json:"voters"`
}

type PollData struct {
	Question       string       `json:"question"`
	Options        []PollOption `json:"options"`
	MultipleChoice bool         `json:"multiple_choice"`
	ExpiresAt      *time.Time   `json:"expires_at,omitempty"`
	TotalVotes     int          `json:"total_votes"`
}
