package service

import (
	"errors"
	"testing"

	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommunityService_CreateCommunity_Success(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)
	mockCommunityModeratorRepo := new(MockCommunityModeratorRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil,
		mockCommunityModeratorRepo,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	desc := "Test Description"
	req := &request.CreateCommunityRequest{
		Name:        "Test Community",
		Description: &desc,
	}

	mockCommunityRepo.On("CreateCommunity", mock.AnythingOfType("*model.Community")).Return(nil)
	mockCommunityModeratorRepo.On("CreateModerator", mock.AnythingOfType("*model.CommunityModerator")).Return(nil)

	err := communityService.CreateCommunity(userID, req)

	assert.NoError(t, err)
	mockCommunityRepo.AssertExpectations(t)
	mockCommunityModeratorRepo.AssertExpectations(t)
}

func TestCommunityService_CreateCommunity_RepositoryError(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	desc := "Test"
	req := &request.CreateCommunityRequest{
		Name:        "Test Community",
		Description: &desc,
	}

	mockCommunityRepo.On("CreateCommunity", mock.AnythingOfType("*model.Community")).Return(errors.New("db error"))

	err := communityService.CreateCommunity(123, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create community")
	mockCommunityRepo.AssertExpectations(t)
}

func TestCommunityService_GetCommunityByID_Success(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)
	mockCommunityModeratorRepo := new(MockCommunityModeratorRepository)
	mockPostRepo := new(MockPostRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil, // subscriptionRepo
		mockCommunityModeratorRepo,
		mockPostRepo,
		nil, nil, nil, nil, nil, nil, nil,
	)

	communityID := uint64(456)
	userID := uint64(123)

	community := &model.Community{
		ID:   communityID,
		Name: "Test Community",
	}

	mockCommunityRepo.On("GetCommunityByIDWithUserSubscription", communityID, &userID).Return(community, uint64(100), nil)
	mockCommunityModeratorRepo.On("GetCommunityModerators", communityID).Return([]*model.CommunityModerator{}, nil)
	mockPostRepo.On("GetPostsLastWeekCount", communityID).Return(int64(5), nil)

	result, err := communityService.GetCommunityByID(communityID, &userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(5), result.PostsLastWeek)
	mockCommunityRepo.AssertExpectations(t)
	mockCommunityModeratorRepo.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
}

func TestCommunityService_GetCommunityByID_NotFound(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	communityID := uint64(999)
	userID := uint64(123)
	mockCommunityRepo.On("GetCommunityByIDWithUserSubscription", communityID, &userID).Return(nil, uint64(0), errors.New("not found"))

	result, err := communityService.GetCommunityByID(communityID, &userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "community not found")
	mockCommunityRepo.AssertExpectations(t)
}

func TestCommunityService_GetCommunityByID_WithModerators(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)
	mockCommunityModeratorRepo := new(MockCommunityModeratorRepository)
	mockPostRepo := new(MockPostRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil, // subscriptionRepo
		mockCommunityModeratorRepo,
		mockPostRepo,
		nil, nil, nil, nil, nil, nil, nil,
	)

	communityID := uint64(456)
	userID := uint64(123)

	community := &model.Community{
		ID:   communityID,
		Name: "Test Community",
	}

	moderators := []*model.CommunityModerator{
		{
			CommunityID: communityID,
			UserID:      1,
			Role:        "super_admin",
			User: &model.User{
				ID:       1,
				Username: "admin",
			},
		},
	}

	mockCommunityRepo.On("GetCommunityByIDWithUserSubscription", communityID, &userID).Return(community, uint64(100), nil)
	mockCommunityModeratorRepo.On("GetCommunityModerators", communityID).Return(moderators, nil)
	mockPostRepo.On("GetPostsLastWeekCount", communityID).Return(int64(10), nil)

	result, err := communityService.GetCommunityByID(communityID, &userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Moderators, 1)
	assert.Equal(t, int64(10), result.PostsLastWeek)
	mockCommunityRepo.AssertExpectations(t)
	mockCommunityModeratorRepo.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
}

func TestCommunityService_UpdateCommunity_Success(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)
	mockCommunityModeratorRepo := new(MockCommunityModeratorRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil,
		mockCommunityModeratorRepo,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	communityID := uint64(456)
	description := "Updated description"
	req := &request.UpdateCommunityRequest{
		Description: &description,
	}

	community := &model.Community{ID: communityID}

	mockCommunityRepo.On("GetCommunityByID", communityID).Return(community, nil)
	mockCommunityModeratorRepo.On("GetModeratorRole", communityID, userID).Return("super_admin", nil)
	mockCommunityRepo.On("UpdateCommunity", communityID, req).Return(nil)

	err := communityService.UpdateCommunity(userID, communityID, req)

	assert.NoError(t, err)
	mockCommunityRepo.AssertExpectations(t)
	mockCommunityModeratorRepo.AssertExpectations(t)
}

func TestCommunityService_UpdateCommunity_NotSuperAdmin(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)
	mockCommunityModeratorRepo := new(MockCommunityModeratorRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil,
		mockCommunityModeratorRepo,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	communityID := uint64(456)
	description := "Updated description"
	req := &request.UpdateCommunityRequest{
		Description: &description,
	}

	community := &model.Community{ID: communityID}

	mockCommunityRepo.On("GetCommunityByID", communityID).Return(community, nil)
	mockCommunityModeratorRepo.On("GetModeratorRole", communityID, userID).Return("moderator", nil)

	err := communityService.UpdateCommunity(userID, communityID, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
	mockCommunityRepo.AssertExpectations(t)
	mockCommunityModeratorRepo.AssertExpectations(t)
}

func TestCommunityService_UpdateCommunity_NotModerator(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)
	mockCommunityModeratorRepo := new(MockCommunityModeratorRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil,
		mockCommunityModeratorRepo,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	communityID := uint64(456)
	description := "Updated description"
	req := &request.UpdateCommunityRequest{
		Description: &description,
	}

	community := &model.Community{ID: communityID}

	mockCommunityRepo.On("GetCommunityByID", communityID).Return(community, nil)
	mockCommunityModeratorRepo.On("GetModeratorRole", communityID, userID).Return("", errors.New("not found"))

	err := communityService.UpdateCommunity(userID, communityID, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
	mockCommunityRepo.AssertExpectations(t)
	mockCommunityModeratorRepo.AssertExpectations(t)
}

func TestCommunityService_DeleteCommunity_Success(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)
	mockCommunityModeratorRepo := new(MockCommunityModeratorRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil,
		mockCommunityModeratorRepo,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	communityID := uint64(456)

	community := &model.Community{ID: communityID}

	mockCommunityRepo.On("GetCommunityByID", communityID).Return(community, nil)
	mockCommunityModeratorRepo.On("GetModeratorRole", communityID, userID).Return("super_admin", nil)
	mockCommunityRepo.On("DeleteCommunity", communityID).Return(nil)

	err := communityService.DeleteCommunity(userID, communityID)

	assert.NoError(t, err)
	mockCommunityRepo.AssertExpectations(t)
	mockCommunityModeratorRepo.AssertExpectations(t)
}

func TestCommunityService_DeleteCommunity_NotSuperAdmin(t *testing.T) {
	mockCommunityRepo := new(MockCommunityRepository)
	mockCommunityModeratorRepo := new(MockCommunityModeratorRepository)

	communityService := NewCommunityService(
		mockCommunityRepo,
		nil,
		mockCommunityModeratorRepo,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	userID := uint64(123)
	communityID := uint64(456)

	community := &model.Community{ID: communityID}

	mockCommunityRepo.On("GetCommunityByID", communityID).Return(community, nil)
	mockCommunityModeratorRepo.On("GetModeratorRole", communityID, userID).Return("moderator", nil)

	err := communityService.DeleteCommunity(userID, communityID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
	mockCommunityRepo.AssertExpectations(t)
	mockCommunityModeratorRepo.AssertExpectations(t)
}
