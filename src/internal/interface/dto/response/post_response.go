package response

import (
	"encoding/json"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/package/template/payload"
	"time"

	"github.com/lib/pq"
)

type CommunityInfo struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar,omitempty"`
}

type AuthorInfo struct {
	ID       uint64  `json:"id"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar,omitempty"`
}

type PollOptionResponse struct {
	ID     int      `json:"id"`
	Text   string   `json:"text"`
	Votes  int      `json:"votes"`
	Voters []uint64 `json:"voters"`
}

type PollDataResponse struct {
	Question       string               `json:"question"`
	Options        []PollOptionResponse `json:"options"`
	MultipleChoice bool                 `json:"multipleChoice"`
	ExpiresAt      *time.Time           `json:"expiresAt,omitempty"`
	TotalVotes     int                  `json:"totalVotes"`
}

func convertPollDataToResponse(pollDataRaw *json.RawMessage) *PollDataResponse {
	if pollDataRaw == nil {
		return nil
	}

	var pollData payload.PollData
	if err := json.Unmarshal(*pollDataRaw, &pollData); err != nil {
		return nil
	}

	options := make([]PollOptionResponse, len(pollData.Options))
	for i, opt := range pollData.Options {
		options[i] = PollOptionResponse{
			ID:     opt.ID,
			Text:   opt.Text,
			Votes:  opt.Votes,
			Voters: opt.Voters,
		}
	}

	return &PollDataResponse{
		Question:       pollData.Question,
		Options:        options,
		MultipleChoice: pollData.MultipleChoice,
		ExpiresAt:      pollData.ExpiresAt,
		TotalVotes:     pollData.TotalVotes,
	}
}

type PostListResponse struct {
	ID          uint64            `json:"id"`
	CommunityID uint64            `json:"communityId"`
	Community   *CommunityInfo    `json:"community,omitempty"`
	AuthorID    uint64            `json:"authorId"`
	Author      *AuthorInfo       `json:"author,omitempty"`
	Title       string            `json:"title"`
	Type        string            `json:"type"`
	Content     string            `json:"content"`
	URL         *string           `json:"url,omitempty"`
	MediaURLs   *pq.StringArray   `json:"mediaUrls,omitempty"`
	PollData    *PollDataResponse `json:"pollData,omitempty"`
	Tags        *pq.StringArray   `json:"tags,omitempty"`
	Vote        int64             `json:"vote"`
	IsVoted     *bool             `json:"isVoted,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt,omitempty"`
}

func NewPostListResponse(post *model.Post) *PostListResponse {
	response := &PostListResponse{
		ID:          post.ID,
		CommunityID: post.CommunityID,
		AuthorID:    post.AuthorID,
		Title:       post.Title,
		Type:        post.Type,
		Content:     post.Content,
		URL:         post.URL,
		MediaURLs:   post.MediaURLs,
		PollData:    convertPollDataToResponse(post.PollData),
		Tags:        post.Tags,
		Vote:        post.Vote,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}

	if post.UserVote != nil {
		isVoted := *post.UserVote == 1
		response.IsVoted = &isVoted
	}

	if post.Community != nil {
		response.Community = &CommunityInfo{
			ID:     post.Community.ID,
			Name:   post.Community.Name,
			Avatar: post.Community.CommunityAvatar,
		}
	}
	if post.Author != nil {
		response.Author = &AuthorInfo{
			ID:       post.Author.ID,
			Username: post.Author.Username,
			Avatar:   post.Author.Avatar,
		}
	}

	return response
}

type PostDetailResponse struct {
	ID        uint64            `json:"id"`
	Community *CommunityInfo    `json:"community,omitempty"`
	Author    *AuthorInfo       `json:"author,omitempty"`
	Title     string            `json:"title"`
	Type      string            `json:"type"`
	Content   string            `json:"content"`
	URL       *string           `json:"url,omitempty"`
	MediaURLs *pq.StringArray   `json:"mediaUrls,omitempty"`
	PollData  *PollDataResponse `json:"pollData,omitempty"`
	Tags      *pq.StringArray   `json:"tags,omitempty"`
	Vote      int64             `json:"vote"`
	IsVoted   *bool             `json:"isVoted,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt *time.Time        `json:"updatedAt,omitempty"`
}

func NewPostDetailResponse(post *model.Post) *PostDetailResponse {
	response := &PostDetailResponse{
		ID:        post.ID,
		Title:     post.Title,
		Type:      post.Type,
		Content:   post.Content,
		URL:       post.URL,
		MediaURLs: post.MediaURLs,
		PollData:  convertPollDataToResponse(post.PollData),
		Tags:      post.Tags,
		Vote:      post.Vote,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	if post.UserVote != nil {
		isVoted := *post.UserVote == 1
		response.IsVoted = &isVoted
	}

	if post.Community != nil {
		response.Community = &CommunityInfo{
			ID:   post.Community.ID,
			Name: post.Community.Name,
		}
	}
	if post.Author != nil {
		response.Author = &AuthorInfo{
			ID:       post.Author.ID,
			Username: post.Author.Username,
			Avatar:   post.Author.Avatar,
		}
	}
	return response
}

type CommunityPostListResponse struct {
	ID          uint64            `json:"id"`
	CommunityID uint64            `json:"communityId"`
	AuthorID    uint64            `json:"authorId"`
	Author      *AuthorInfo       `json:"author,omitempty"`
	Title       string            `json:"title"`
	Type        string            `json:"type"`
	Content     string            `json:"content"`
	URL         *string           `json:"url,omitempty"`
	MediaURLs   *pq.StringArray   `json:"mediaUrls,omitempty"`
	PollData    *PollDataResponse `json:"pollData,omitempty"`
	Tags        *pq.StringArray   `json:"tags,omitempty"`
	Status      string            `json:"status"`
	Vote        int64             `json:"vote"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt,omitempty"`
}

func NewCommunityPostListResponse(post *model.Post) *CommunityPostListResponse {
	response := &CommunityPostListResponse{
		ID:          post.ID,
		CommunityID: post.CommunityID,
		AuthorID:    post.AuthorID,
		Title:       post.Title,
		Type:        post.Type,
		Content:     post.Content,
		URL:         post.URL,
		MediaURLs:   post.MediaURLs,
		PollData:    convertPollDataToResponse(post.PollData),
		Tags:        post.Tags,
		Status:      post.Status,
		Vote:        post.Vote,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
	if post.Author != nil {
		response.Author = &AuthorInfo{
			ID:       post.Author.ID,
			Username: post.Author.Username,
			Avatar:   post.Author.Avatar,
		}
	}
	return response
}

type SavedPostResponse struct {
	PostID     uint64         `json:"postId"`
	Title      string         `json:"title"`
	Community  *CommunityInfo `json:"community,omitempty"`
	Author     AuthorInfo     `json:"author"`
	IsFollowed bool           `json:"isFollowed"`
	CreatedAt  time.Time      `json:"createdAt"`
}

func NewSavedPostResponse(savedPost *model.UserSavedPost) *SavedPostResponse {
	response := &SavedPostResponse{
		PostID: savedPost.PostID,
		Title:  savedPost.PostTitle,
		Author: AuthorInfo{
			ID:       savedPost.AuthorID,
			Username: savedPost.AuthorName,
			Avatar:   savedPost.AuthorAvatar,
		},
		IsFollowed: savedPost.IsFollowed,
		CreatedAt:  savedPost.PostCreatedAt,
	}

	if savedPost.Community != nil {
		response.Community = &CommunityInfo{
			ID:     savedPost.Community.ID,
			Name:   savedPost.Community.Name,
			Avatar: savedPost.Community.CommunityAvatar,
		}
	}

	return response
}
