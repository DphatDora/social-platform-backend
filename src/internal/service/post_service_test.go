package service

import (
	"errors"
	"testing"

	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/package/constant"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostService_CreatePost_TextSuccess(t *testing.T) {
	mockPostRepo := new(MockPostRepository)
	mockCommunityRepo := new(MockCommunityRepository)

	postService := NewPostService(
		mockPostRepo,
		mockCommunityRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	req := &request.CreatePostRequest{
		CommunityID: 1,
		Title:       "Test Post",
		Type:        constant.PostTypeText,
		Content:     *stringPtr("Test content"),
	}

	community := &model.Community{
		ID:                   1,
		Name:                 "test-community",
		RequiresPostApproval: false,
	}

	mockCommunityRepo.On("GetCommunityByID", req.CommunityID).Return(community, nil)
	mockPostRepo.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(nil)

	err := postService.CreatePost(userID, req)

	assert.NoError(t, err)
	mockCommunityRepo.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
}

func TestPostService_CreatePost_CommunityNotFound(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)

	postService := NewPostService(
		nil,
		mockCommunityRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	req := &request.CreatePostRequest{
		CommunityID: 999,
		Title:       "Test Post",
		Type:        constant.PostTypeText,
	}

	mockCommunityRepo.On("GetCommunityByID", req.CommunityID).Return(nil, errors.New("not found"))

	err := postService.CreatePost(123, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "community not found")
	mockCommunityRepo.AssertExpectations(t)
}

func TestPostService_CreatePost_LinkWithoutURL(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)

	postService := NewPostService(
		nil,
		mockCommunityRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	req := &request.CreatePostRequest{
		CommunityID: 1,
		Title:       "Test Link Post",
		Type:        constant.PostTypeLink,
		URL:         nil,
	}

	community := &model.Community{ID: 1}
	mockCommunityRepo.On("GetCommunityByID", req.CommunityID).Return(community, nil)

	err := postService.CreatePost(123, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "url is required")
	mockCommunityRepo.AssertExpectations(t)
}

func TestPostService_CreatePost_MediaWithoutURLs(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)

	postService := NewPostService(
		nil,
		mockCommunityRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	req := &request.CreatePostRequest{
		CommunityID: 1,
		Title:       "Test Media Post",
		Type:        constant.PostTypeMedia,
		MediaURLs:   nil,
	}

	community := &model.Community{ID: 1}
	mockCommunityRepo.On("GetCommunityByID", req.CommunityID).Return(community, nil)

	err := postService.CreatePost(123, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "media_urls are required")
	mockCommunityRepo.AssertExpectations(t)
}

func TestPostService_CreatePost_InvalidType(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)

	postService := NewPostService(
		nil,
		mockCommunityRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	req := &request.CreatePostRequest{
		CommunityID: 1,
		Title:       "Test Post",
		Type:        "invalid_type",
	}

	community := &model.Community{ID: 1}
	mockCommunityRepo.On("GetCommunityByID", req.CommunityID).Return(community, nil)

	err := postService.CreatePost(123, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid post type")
	mockCommunityRepo.AssertExpectations(t)
}

func TestPostService_UpdatePost_Success(t *testing.T) {
	mockPostRepo := new(MockPostRepository)

	postService := NewPostService(
		mockPostRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	postID := uint64(456)
	updateReq := &request.UpdatePostTextRequest{
		Title:   stringPtr("Updated Title"),
		Content: stringPtr("Updated Content"),
	}

	post := &model.Post{
		ID:       postID,
		AuthorID: userID,
		Type:     constant.PostTypeText,
	}

	mockPostRepo.On("GetPostByID", postID).Return(post, nil)
	mockPostRepo.On("UpdatePostText", postID, updateReq).Return(nil)

	err := postService.UpdatePost(userID, postID, constant.PostTypeText, updateReq)

	assert.NoError(t, err)
	mockPostRepo.AssertExpectations(t)
}

func TestPostService_UpdatePost_NotAuthor(t *testing.T) {
	mockPostRepo := new(MockPostRepository)

	postService := NewPostService(
		mockPostRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	postID := uint64(456)
	updateReq := &request.UpdatePostTextRequest{
		Title: stringPtr("Updated Title"),
	}

	post := &model.Post{
		ID:       postID,
		AuthorID: 999,
		Type:     constant.PostTypeText,
	}

	mockPostRepo.On("GetPostByID", postID).Return(post, nil)

	err := postService.UpdatePost(userID, postID, constant.PostTypeText, updateReq)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
	mockPostRepo.AssertExpectations(t)
}

func TestPostService_UpdatePost_TypeMismatch(t *testing.T) {
	mockPostRepo := new(MockPostRepository)

	postService := NewPostService(
		mockPostRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	postID := uint64(456)
	updateReq := &request.UpdatePostTextRequest{
		Title: stringPtr("Updated Title"),
	}

	post := &model.Post{
		ID:       postID,
		AuthorID: userID,
		Type:     constant.PostTypeLink,
	}

	mockPostRepo.On("GetPostByID", postID).Return(post, nil)

	err := postService.UpdatePost(userID, postID, constant.PostTypeText, updateReq)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "post type mismatch")
	mockPostRepo.AssertExpectations(t)
}

func TestPostService_DeletePost_Success(t *testing.T) {
	mockPostRepo := new(MockPostRepository)

	postService := NewPostService(
		mockPostRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	postID := uint64(456)

	post := &model.Post{
		ID:       postID,
		AuthorID: userID,
	}

	mockPostRepo.On("GetPostByID", postID).Return(post, nil)
	mockPostRepo.On("DeletePost", postID).Return(nil)

	err := postService.DeletePost(userID, postID)

	assert.NoError(t, err)
	mockPostRepo.AssertExpectations(t)
}

func TestPostService_DeletePost_NotFound(t *testing.T) {
	mockPostRepo := new(MockPostRepository)

	postService := NewPostService(
		mockPostRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	postID := uint64(999)
	mockPostRepo.On("GetPostByID", postID).Return(nil, errors.New("not found"))

	err := postService.DeletePost(123, postID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "post not found")
	mockPostRepo.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
