package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/util"
	"strings"
)

type UserService struct {
	userRepo               repository.UserRepository
	communityModeratorRepo repository.CommunityModeratorRepository
	userSavedPostRepo      repository.UserSavedPostRepository
}

func NewUserService(userRepo repository.UserRepository, communityModeratorRepo repository.CommunityModeratorRepository, userSavedPostRepo repository.UserSavedPostRepository) *UserService {
	return &UserService{
		userRepo:               userRepo,
		communityModeratorRepo: communityModeratorRepo,
		userSavedPostRepo:      userSavedPostRepo,
	}
}

func (s *UserService) GetUserProfile(userID uint64) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("[Err] Error getting user by ID in UserService.GetUserProfile: %v", err)
		return nil, fmt.Errorf("user not found")
	}

	// Get user achievement
	achievement := response.UserAchievement{
		Karma:         0,
		Badge:         "",
		TotalPosts:    0,
		TotalComments: 0,
	}

	// Get latest badge
	userBadge, err := s.userRepo.GetLatestUserBadge(userID)
	if err == nil && userBadge != nil {
		achievement.Karma = userBadge.Karma
		if userBadge.Badge != nil {
			achievement.Badge = userBadge.Badge.Name
		}
	}

	// Get total posts
	postCount, err := s.userRepo.GetUserPostCount(userID)
	if err == nil {
		achievement.TotalPosts = postCount
	}

	// Get total comments
	commentCount, err := s.userRepo.GetUserCommentCount(userID)
	if err == nil {
		achievement.TotalComments = commentCount
	}

	userProfile := response.NewUserProfileResponse(user, achievement)
	return userProfile, nil
}

func (s *UserService) UpdateUserProfile(userID uint64, updateReq *request.UpdateUserProfileRequest) error {
	err := s.userRepo.UpdateUserProfile(userID, updateReq)
	if err != nil {
		log.Printf("[Err] Error updating user profile in UserService.UpdateUserProfile: %v", err)
		return fmt.Errorf("failed to update user profile")
	}
	return nil
}

func (s *UserService) ChangePassword(userID uint64, changePasswordReq *request.ChangePasswordRequest) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("[Err] Error getting user by ID in UserService.ChangePassword: %v", err)
		return fmt.Errorf("user not found")
	}

	if err := util.ComparePassword(user.Password, changePasswordReq.OldPassword); err != nil {
		log.Printf("[Err] Old password is incorrect in UserService.ChangePassword for user ID: %d", userID)
		return fmt.Errorf("old password is incorrect")
	}

	hashedPassword, err := util.HashPassword(changePasswordReq.NewPassword)
	if err != nil {
		log.Printf("[Err] Error hashing password in UserService.ChangePassword: %v", err)
		return fmt.Errorf("failed to hash password")
	}

	if err := s.userRepo.UpdatePasswordAndSetChangedAt(userID, hashedPassword); err != nil {
		log.Printf("[Err] Error updating password in UserService.ChangePassword: %v", err)
		return fmt.Errorf("failed to update password")
	}

	return nil
}

func (s *UserService) GetUserConfig(userID uint64) (*response.UserConfigResponse, error) {
	// Get user information
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("[Err] Error getting user by ID in UserService.GetUserConfig: %v", err)
		return nil, fmt.Errorf("user not found")
	}

	// Get moderated communities
	moderators, err := s.communityModeratorRepo.GetModeratorCommunitiesByUserID(userID)
	if err != nil {
		log.Printf("[Err] Error getting moderated communities in UserService.GetUserConfig: %v", err)
		// Don't fail if we can't get moderated communities
		moderators = []*model.CommunityModerator{}
	}

	moderatedCommunities := make([]response.CommunityModerator, len(moderators))
	for i, mod := range moderators {
		moderatedCommunities[i] = response.CommunityModerator{
			CommunityID: mod.CommunityID,
			Role:        mod.Role,
		}
	}

	userConfig := &response.UserConfigResponse{
		Username:             user.Username,
		Avatar:               user.Avatar,
		ModeratedCommunities: moderatedCommunities,
	}

	return userConfig, nil
}

func (s *UserService) GetUserBadgeHistory(userID uint64) ([]*response.UserBadgeResponse, error) {
	userBadges, err := s.userRepo.GetUserBadgeHistory(userID)
	if err != nil {
		log.Printf("[Err] Error getting user badge history in UserService.GetUserBadgeHistory: %v", err)
		return nil, fmt.Errorf("failed to get user badge history")
	}

	badgeResponses := make([]*response.UserBadgeResponse, len(userBadges))
	for i, userBadge := range userBadges {
		badgeName := ""
		if userBadge.Badge != nil {
			badgeName = userBadge.Badge.Name
		}
		badgeResponses[i] = response.NewUserBadgeResponse(badgeName, userBadge.Badge.IconURL, userBadge.MonthYear, userBadge.Karma)
	}

	return badgeResponses, nil
}

func (s *UserService) GetUserSavedPosts(userID uint64, searchTitle string, isFollowed *bool, page, limit int) ([]*response.SavedPostResponse, *response.Pagination, error) {
	savedPosts, total, err := s.userSavedPostRepo.GetUserSavedPosts(userID, searchTitle, isFollowed, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting user saved posts in UserService.GetUserSavedPosts: %v", err)
		return nil, nil, fmt.Errorf("failed to get saved posts")
	}

	savedPostResponses := make([]*response.SavedPostResponse, len(savedPosts))
	for i, savedPost := range savedPosts {
		savedPostResponses[i] = response.NewSavedPostResponse(savedPost)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		nextURL := fmt.Sprintf("/api/v1/users/saved-posts?page=%d&limit=%d", page+1, limit)
		if searchTitle != "" {
			nextURL += fmt.Sprintf("&search=%s", searchTitle)
		}
		if isFollowed != nil {
			nextURL += fmt.Sprintf("&isFollowed=%t", *isFollowed)
		}
		pagination.NextURL = nextURL
	}

	return savedPostResponses, pagination, nil
}

func (s *UserService) CreateUserSavedPost(userID uint64, savedPostReq *request.UserSavedPostRequest) error {
	if err := s.userSavedPostRepo.CreateUserSavedPost(userID, savedPostReq); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return fmt.Errorf("post already saved")
		}

		log.Printf("[Err] Error creating user saved post in UserService.CreateUserSavedPost: %v", err)
		return fmt.Errorf("failed to save post")
	}
	return nil
}

func (s *UserService) UpdateUserSavedPostFollowStatus(userID, postID uint64, updateReq *request.UpdateUserSavedPostRequest) error {
	if err := s.userSavedPostRepo.UpdateFollowedStatus(userID, postID, updateReq.IsFollowed); err != nil {
		log.Printf("[Err] Error updating user saved post follow status in UserService.UpdateUserSavedPostFollowStatus: %v", err)
		return fmt.Errorf("failed to update follow status")
	}
	return nil
}

func (s *UserService) DeleteUserSavedPost(userID, postID uint64) error {
	if err := s.userSavedPostRepo.DeleteUserSavedPost(userID, postID); err != nil {
		log.Printf("[Err] Error deleting user saved post in UserService.DeleteUserSavedPost: %v", err)
		return fmt.Errorf("failed to delete saved post")
	}
	return nil
}
