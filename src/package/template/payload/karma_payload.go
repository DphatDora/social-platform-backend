package payload

import "time"

type UpdateUserKarmaPayload struct {
	UserId    uint64    `json:"user_id"`
	TargetId  *uint64   `json:"target_id,omitempty"`
	Action    string    `json:"action"`
	UpdatedAt time.Time `json:"updated_at"`
}
