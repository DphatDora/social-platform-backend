package service

import (
	"errors"
	"testing"

	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentService_CreateComment_Success(t *testing.T) {
	mockCommentRepo := new(MockCommentRepository)
	mockPostRepo := new(MockPostRepository)

	commentService := NewCommentService(
		mockCommentRepo,
		mockPostRepo,
		nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	req := &request.CreateCommentRequest{
		PostID:  456,
		Content: "Test comment",
	}

	post := &model.Post{
		ID:       456,
		AuthorID: 789,
	}

	mockPostRepo.On("GetPostByID", req.PostID).Return(post, nil)
	mockCommentRepo.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(nil)

	err := commentService.CreateComment(userID, req)

	assert.NoError(t, err)
	mockPostRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_CreateComment_PostNotFound(t *testing.T) {
	mockPostRepo := new(MockPostRepository)

	commentService := NewCommentService(
		nil,
		mockPostRepo,
		nil, nil, nil, nil, nil, nil,
	)

	req := &request.CreateCommentRequest{
		PostID:  999,
		Content: "Test comment",
	}

	mockPostRepo.On("GetPostByID", req.PostID).Return(nil, errors.New("not found"))

	err := commentService.CreateComment(123, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "post not found")
	mockPostRepo.AssertExpectations(t)
}

func TestCommentService_CreateComment_WithParent(t *testing.T) {
	mockCommentRepo := new(MockCommentRepository)
	mockPostRepo := new(MockPostRepository)

	commentService := NewCommentService(
		mockCommentRepo,
		mockPostRepo,
		nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	parentCommentID := uint64(111)
	req := &request.CreateCommentRequest{
		PostID:          456,
		ParentCommentID: &parentCommentID,
		Content:         "Reply comment",
	}

	post := &model.Post{
		ID:       456,
		AuthorID: 789,
	}

	parentComment := &model.Comment{
		ID:       parentCommentID,
		PostID:   456,
		AuthorID: 999,
	}

	mockPostRepo.On("GetPostByID", req.PostID).Return(post, nil)
	mockCommentRepo.On("GetCommentByID", parentCommentID).Return(parentComment, nil)
	mockCommentRepo.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(nil)

	err := commentService.CreateComment(userID, req)

	assert.NoError(t, err)
	mockPostRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_CreateComment_ParentNotFound(t *testing.T) {
	mockCommentRepo := new(MockCommentRepository)
	mockPostRepo := new(MockPostRepository)

	commentService := NewCommentService(
		mockCommentRepo,
		mockPostRepo,
		nil, nil, nil, nil, nil, nil,
	)

	parentCommentID := uint64(999)
	req := &request.CreateCommentRequest{
		PostID:          456,
		ParentCommentID: &parentCommentID,
		Content:         "Reply comment",
	}

	post := &model.Post{ID: 456}

	mockPostRepo.On("GetPostByID", req.PostID).Return(post, nil)
	mockCommentRepo.On("GetCommentByID", parentCommentID).Return(nil, errors.New("not found"))

	err := commentService.CreateComment(123, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parent comment not found")
	mockPostRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_CreateComment_ParentDifferentPost(t *testing.T) {
	mockCommentRepo := new(MockCommentRepository)
	mockPostRepo := new(MockPostRepository)

	commentService := NewCommentService(
		mockCommentRepo,
		mockPostRepo,
		nil, nil, nil, nil, nil, nil,
	)

	parentCommentID := uint64(111)
	req := &request.CreateCommentRequest{
		PostID:          456,
		ParentCommentID: &parentCommentID,
		Content:         "Reply comment",
	}

	post := &model.Post{ID: 456}
	parentComment := &model.Comment{
		ID:     parentCommentID,
		PostID: 789,
	}

	mockPostRepo.On("GetPostByID", req.PostID).Return(post, nil)
	mockCommentRepo.On("GetCommentByID", parentCommentID).Return(parentComment, nil)

	err := commentService.CreateComment(123, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parent comment does not belong to this post")
	mockPostRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_UpdateComment_Success(t *testing.T) {
	mockCommentRepo := new(MockCommentRepository)

	commentService := NewCommentService(
		mockCommentRepo,
		nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	commentID := uint64(456)
	req := &request.UpdateCommentRequest{
		Content: "Updated content",
	}

	comment := &model.Comment{
		ID:       commentID,
		AuthorID: userID,
	}

	mockCommentRepo.On("GetCommentByID", commentID).Return(comment, nil)
	mockCommentRepo.On("UpdateComment", commentID, req.Content, req.MediaURL).Return(nil)

	err := commentService.UpdateComment(userID, commentID, req)

	assert.NoError(t, err)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_UpdateComment_NotAuthor(t *testing.T) {
	mockCommentRepo := new(MockCommentRepository)

	commentService := NewCommentService(
		mockCommentRepo,
		nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	commentID := uint64(456)
	req := &request.UpdateCommentRequest{
		Content: "Updated content",
	}

	comment := &model.Comment{
		ID:       commentID,
		AuthorID: 999,
	}

	mockCommentRepo.On("GetCommentByID", commentID).Return(comment, nil)

	err := commentService.UpdateComment(userID, commentID, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_DeleteComment_Success(t *testing.T) {
	mockCommentRepo := new(MockCommentRepository)

	commentService := NewCommentService(
		mockCommentRepo,
		nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	commentID := uint64(456)

	comment := &model.Comment{
		ID:              commentID,
		AuthorID:        userID,
		ParentCommentID: nil,
	}

	mockCommentRepo.On("GetCommentByID", commentID).Return(comment, nil)
	mockCommentRepo.On("DeleteComment", commentID, comment.ParentCommentID).Return(nil)

	err := commentService.DeleteComment(userID, commentID)

	assert.NoError(t, err)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_DeleteComment_NotFound(t *testing.T) {
	mockCommentRepo := new(MockCommentRepository)

	commentService := NewCommentService(
		mockCommentRepo,
		nil, nil, nil, nil, nil, nil, nil,
	)

	commentID := uint64(999)
	mockCommentRepo.On("GetCommentByID", commentID).Return(nil, errors.New("not found"))

	err := commentService.DeleteComment(123, commentID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "comment not found")
	mockCommentRepo.AssertExpectations(t)
}
