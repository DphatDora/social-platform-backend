package request

import (
	"encoding/json"

	"github.com/lib/pq"
)

type CreatePostRequest struct {
	CommunityID uint64           `json:"communityId" binding:"required"`
	Title       string           `json:"title" binding:"required"`
	Type        string           `json:"type" binding:"required,oneof=text link media poll"`
	Content     string           `json:"content" binding:"required"`
	URL         *string          `json:"url,omitempty"`
	MediaURLs   *pq.StringArray  `json:"mediaUrls,omitempty"`
	PollData    *json.RawMessage `json:"pollData,omitempty"`
	Tags        *pq.StringArray  `json:"tags,omitempty"`
}

type UpdatePostTextRequest struct {
	Title   *string         `json:"title,omitempty"`
	Content *string         `json:"content,omitempty"`
	Tags    *pq.StringArray `json:"tags,omitempty"`
}

type UpdatePostLinkRequest struct {
	Title   *string         `json:"title,omitempty"`
	Content *string         `json:"content,omitempty"`
	URL     *string         `json:"url,omitempty"`
	Tags    *pq.StringArray `json:"tags,omitempty"`
}

type UpdatePostMediaRequest struct {
	Title     *string         `json:"title,omitempty"`
	Content   *string         `json:"content,omitempty"`
	MediaURLs *pq.StringArray `json:"mediaUrls,omitempty"`
	Tags      *pq.StringArray `json:"tags,omitempty"`
}

type UpdatePostPollRequest struct {
	Title    *string          `json:"title,omitempty"`
	Content  *string          `json:"content,omitempty"`
	PollData *json.RawMessage `json:"pollData,omitempty"`
	Tags     *pq.StringArray  `json:"tags,omitempty"`
}

type VotePostRequest struct {
	Vote bool `json:"vote"`
}

type VotePollRequest struct {
	OptionID int `json:"optionId" binding:"required"`
}

type UpdatePostStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending approved rejected"`
}

type ReportPostRequest struct {
	Reasons []string `json:"reasons" binding:"required,min=1"`
	Note    *string  `json:"note,omitempty"`
}
