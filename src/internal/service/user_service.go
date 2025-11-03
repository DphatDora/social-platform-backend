package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/util"
)

type UserService struct {
	userRepo               repository.UserRepository
	communityModeratorRepo repository.CommunityModeratorRepository
}

func NewUserService(userRepo repository.UserRepository, communityModeratorRepo repository.CommunityModeratorRepository) *UserService {
	return &UserService{
		userRepo:               userRepo,
		communityModeratorRepo: communityModeratorRepo,
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
