package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/util"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepo               repository.UserRepository
	communityRepo          repository.CommunityRepository
	communityModeratorRepo repository.CommunityModeratorRepository
	userSavedPostRepo      repository.UserSavedPostRepository
	postRepo               repository.PostRepository
	botTaskService         *BotTaskService
	redisClient            *redis.Client
}

func NewUserService(
	userRepo repository.UserRepository,
	communityRepo repository.CommunityRepository,
	communityModeratorRepo repository.CommunityModeratorRepository,
	userSavedPostRepo repository.UserSavedPostRepository,
	postRepo repository.PostRepository,
	botTaskService *BotTaskService,
	redisClient *redis.Client,
) *UserService {
	return &UserService{
		userRepo:               userRepo,
		communityRepo:          communityRepo,
		communityModeratorRepo: communityModeratorRepo,
		userSavedPostRepo:      userSavedPostRepo,
		postRepo:               postRepo,
		botTaskService:         botTaskService,
		redisClient:            redisClient,
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

	if user.Password == nil {
		log.Printf("[Err] User ID %d registered with Google, cannot change password", userID)
		return fmt.Errorf("this account is registered with Google and does not have a password")
	}

	if err := util.ComparePassword(*user.Password, changePasswordReq.OldPassword); err != nil {
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

	// Invalidate password cache after successful password change
	if s.redisClient != nil {
		if err := util.InvalidatePasswordCache(s.redisClient, userID); err != nil {
			log.Printf("[Warn] Error invalidating password cache for user %d: %v", userID, err)
		}
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
	// Check if post is already saved
	exists, err := s.userSavedPostRepo.CheckUserSavedPostExists(userID, savedPostReq.PostID)
	if err != nil {
		log.Printf("[Err] Error checking if post is saved in UserService.CreateUserSavedPost: %v", err)
		return fmt.Errorf("failed to check saved post status")
	}

	// If already saved and trying to follow, just update the follow status
	if exists {
		if savedPostReq.IsFollowed {
			if err := s.userSavedPostRepo.UpdateFollowedStatus(userID, savedPostReq.PostID, savedPostReq.IsFollowed); err != nil {
				log.Printf("[Err] Error updating follow status in UserService.CreateUserSavedPost: %v", err)
				return fmt.Errorf("failed to update follow status")
			}
			return nil
		}
		return fmt.Errorf("post already saved")
	}

	// Create new saved post
	if err := s.userSavedPostRepo.CreateUserSavedPost(userID, savedPostReq); err != nil {
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

	// Create bot task for interest score if user is following the post
	if updateReq.IsFollowed {
		go func(userID, postID uint64) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[PanicRecovered] in goroutine UserService.UpdateUserSavedPostFollowStatus: %v", r)
				}
			}()

			// Get post to find community ID
			post, err := s.postRepo.GetPostByID(postID)
			if err != nil {
				log.Printf("[Warn] Error getting post for interest score in UserService.UpdateUserSavedPostFollowStatus: %v", err)
				return
			}

			postIDPtr := &postID
			if err := s.botTaskService.CreateInterestScoreTask(
				userID,
				post.CommunityID,
				"follow_post",
				postIDPtr,
			); err != nil {
				log.Printf("[Err] Error creating interest score task in goroutine (UserService.UpdateUserSavedPostFollowStatus): %v", err)
			}
		}(userID, postID)
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

func (s *UserService) SearchUsers(searchTerm string, page, limit int) ([]*response.UserSearchResponse, *response.Pagination, error) {
	users, total, err := s.userRepo.SearchUsers(searchTerm, page, limit)
	if err != nil {
		log.Printf("[Err] Error searching users in UserService.SearchUsers: %v", err)
		return nil, nil, fmt.Errorf("failed to search users")
	}

	userResponses := make([]*response.UserSearchResponse, len(users))
	for i, user := range users {
		// Get latest badge/karma for each user
		karma := uint64(0)
		userBadge, err := s.userRepo.GetLatestUserBadge(user.ID)
		if err == nil && userBadge != nil {
			karma = userBadge.Karma
		}

		userResponses[i] = response.NewUserSearchResponse(user, karma)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		nextURL := fmt.Sprintf("/api/v1/users/search?page=%d&limit=%d", page+1, limit)
		if searchTerm != "" {
			nextURL += fmt.Sprintf("&search=%s", searchTerm)
		}
		pagination.NextURL = nextURL
	}

	return userResponses, pagination, nil
}

func (s *UserService) GetUserSuperAdminCommunities(userID uint64) ([]*response.CommunityListResponse, error) {
	communities, err := s.communityRepo.GetCommunitiesByCreatorID(userID)
	if err != nil {
		log.Printf("[Err] Error getting communities by creator in UserService.GetUserSuperAdminCommunities: %v", err)
		return nil, fmt.Errorf("failed to get super admin communities")
	}

	communityResponses := make([]*response.CommunityListResponse, len(communities))
	for i, community := range communities {
		resp := response.NewCommunityListResponse(community)
		resp.TotalMembers = community.MemberCount
		communityResponses[i] = resp
	}

	return communityResponses, nil
}

func (s *UserService) GetUserAdminCommunities(userID uint64) ([]*response.CommunityListResponse, error) {
	communities, err := s.communityRepo.GetCommunitiesByModeratorID(userID, constant.ROLE_ADMIN)
	if err != nil {
		log.Printf("[Err] Error getting communities by moderator in UserService.GetUserAdminCommunities: %v", err)
		return nil, fmt.Errorf("failed to get admin communities")
	}

	communityResponses := make([]*response.CommunityListResponse, len(communities))
	for i, community := range communities {
		resp := response.NewCommunityListResponse(community)
		resp.TotalMembers = community.MemberCount
		communityResponses[i] = resp
	}

	return communityResponses, nil
}
