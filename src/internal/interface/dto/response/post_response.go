package response

import (
	"encoding/json"
	"social-platform-backend/internal/domain/model"
	"time"

	"github.com/lib/pq"
)

type CommunityInfo struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type AuthorInfo struct {
	ID       uint64  `json:"id"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar,omitempty"`
}

type PostListResponse struct {
	ID          uint64           `json:"id"`
	CommunityID uint64           `json:"communityId"`
	Community   *CommunityInfo   `json:"community,omitempty"`
	AuthorID    uint64           `json:"authorId"`
	Author      *AuthorInfo      `json:"author,omitempty"`
	Title       string           `json:"title"`
	Type        string           `json:"type"`
	Content     string           `json:"content"`
	URL         *string          `json:"url,omitempty"`
	MediaURLs   *pq.StringArray  `json:"mediaUrls,omitempty"`
	PollData    *json.RawMessage `json:"pollData,omitempty"`
	Tags        *pq.StringArray  `json:"tags,omitempty"`
	Vote        int64            `json:"vote"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   *time.Time       `json:"updatedAt,omitempty"`
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
		PollData:    post.PollData,
		Tags:        post.Tags,
		Vote:        post.Vote,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
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

type PostDetailResponse struct {
	ID        uint64           `json:"id"`
	Community *CommunityInfo   `json:"community,omitempty"`
	Author    *AuthorInfo      `json:"author,omitempty"`
	Title     string           `json:"title"`
	Type      string           `json:"type"`
	Content   string           `json:"content"`
	URL       *string          `json:"url,omitempty"`
	MediaURLs *pq.StringArray  `json:"mediaUrls,omitempty"`
	PollData  *json.RawMessage `json:"pollData,omitempty"`
	Tags      *pq.StringArray  `json:"tags,omitempty"`
	Vote      int64            `json:"vote"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt *time.Time       `json:"updatedAt,omitempty"`
}

func NewPostDetailResponse(post *model.Post) *PostDetailResponse {
	response := &PostDetailResponse{
		ID:        post.ID,
		Title:     post.Title,
		Type:      post.Type,
		Content:   post.Content,
		URL:       post.URL,
		MediaURLs: post.MediaURLs,
		PollData:  post.PollData,
		Tags:      post.Tags,
		Vote:      post.Vote,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
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
	ID          uint64           `json:"id"`
	CommunityID uint64           `json:"communityId"`
	AuthorID    uint64           `json:"authorId"`
	Author      *AuthorInfo      `json:"author,omitempty"`
	Title       string           `json:"title"`
	Type        string           `json:"type"`
	Content     string           `json:"content"`
	URL         *string          `json:"url,omitempty"`
	MediaURLs   *pq.StringArray  `json:"mediaUrls,omitempty"`
	PollData    *json.RawMessage `json:"pollData,omitempty"`
	Tags        *pq.StringArray  `json:"tags,omitempty"`
	Status      string           `json:"status"`
	Vote        int64            `json:"vote"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   *time.Time       `json:"updatedAt,omitempty"`
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
		PollData:    post.PollData,
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
	PostID     uint64     `json:"postId"`
	Title      string     `json:"title"`
	Author     AuthorInfo `json:"author"`
	IsFollowed bool       `json:"isFollowed"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func NewSavedPostResponse(savedPost *model.UserSavedPost) *SavedPostResponse {
	return &SavedPostResponse{
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
}
