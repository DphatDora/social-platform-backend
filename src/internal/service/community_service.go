package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/template/payload"
	"strings"
	"time"
)

type CommunityService struct {
	communityRepo          repository.CommunityRepository
	subscriptionRepo       repository.SubscriptionRepository
	communityModeratorRepo repository.CommunityModeratorRepository
	postRepo               repository.PostRepository
	postReportRepo         repository.PostReportRepository
	commentRepo            repository.CommentRepository
	commentReportRepo      repository.CommentReportRepository
	topicRepo              repository.TopicRepository
	userRestrictionRepo    repository.UserRestrictionRepository
	notificationService    *NotificationService
	botTaskService         *BotTaskService
}

func NewCommunityService(
	communityRepo repository.CommunityRepository,
	subscriptionRepo repository.SubscriptionRepository,
	communityModeratorRepo repository.CommunityModeratorRepository,
	postRepo repository.PostRepository,
	postReportRepo repository.PostReportRepository,
	commentRepo repository.CommentRepository,
	commentReportRepo repository.CommentReportRepository,
	topicRepo repository.TopicRepository,
	userRestrictionRepo repository.UserRestrictionRepository,
	notificationService *NotificationService,
	botTaskService *BotTaskService,
) *CommunityService {
	return &CommunityService{
		communityRepo:          communityRepo,
		subscriptionRepo:       subscriptionRepo,
		communityModeratorRepo: communityModeratorRepo,
		postRepo:               postRepo,
		postReportRepo:         postReportRepo,
		commentRepo:            commentRepo,
		commentReportRepo:      commentReportRepo,
		topicRepo:              topicRepo,
		userRestrictionRepo:    userRestrictionRepo,
		notificationService:    notificationService,
		botTaskService:         botTaskService,
	}
}

func (s *CommunityService) CreateCommunity(userID uint64, req *request.CreateCommunityRequest) error {
	community := &model.Community{
		Name:             req.Name,
		ShortDescription: req.ShortDescription,
		Description:      req.Description,
		Topic:            req.Topic,
		CommunityAvatar:  req.CommunityAvatar,
		CoverImage:       req.CoverImage,
		IsPrivate:        req.IsPrivate,
		CreatedBy:        userID,
		CreatedAt:        time.Now(),
	}

	if err := s.communityRepo.CreateCommunity(community); err != nil {
		log.Printf("[Err] Error creating community in CommunityService.CreateCommunity: %v", err)
		return fmt.Errorf("failed to create community")
	}

	// Add creator as super admin
	moderator := &model.CommunityModerator{
		CommunityID: community.ID,
		UserID:      userID,
		Role:        constant.ROLE_SUPER_ADMIN,
		JoinedAt:    time.Now(),
	}

	if err := s.communityModeratorRepo.CreateModerator(moderator); err != nil {
		log.Printf("[Err] Error creating moderator in CommunityService.CreateCommunity: %v", err)
		return fmt.Errorf("failed to create moderator")
	}

	return nil
}

func (s *CommunityService) GetCommunityByID(id uint64, userID *uint64) (*response.CommunityDetailResponse, error) {
	community, memberCount, err := s.communityRepo.GetCommunityByIDWithUserSubscription(id, userID)
	if err != nil {
		log.Printf("[Err] Error getting community by ID in CommunityService.GetCommunityByID: %v", err)
		return nil, fmt.Errorf("community not found")
	}

	// Get community moderators
	moderators, err := s.communityModeratorRepo.GetCommunityModerators(id)
	if err != nil {
		log.Printf("[Err] Error getting community moderators in CommunityService.GetCommunityByID: %v", err)
		moderators = []*model.CommunityModerator{}
	}

	moderatorResponses := make([]response.ModeratorResponse, len(moderators))
	for i, mod := range moderators {
		moderatorResponses[i] = *response.NewModeratorResponse(mod.User, mod.Role)
	}

	// Get posts count in last week
	postsLastWeek, err := s.postRepo.GetPostsLastWeekCount(id)
	if err != nil {
		log.Printf("[Err] Error getting posts last week count in CommunityService.GetCommunityByID: %v", err)
		postsLastWeek = 0
	}

	communityResponse := response.NewCommunityDetailResponse(community)
	communityResponse.TotalMembers = memberCount
	communityResponse.PostsLastWeek = postsLastWeek
	communityResponse.Moderators = moderatorResponses

	if community.IsSubscribed != nil {
		communityResponse.IsFollow = community.IsSubscribed
	}

	if community.IsRequestJoin != nil {
		communityResponse.IsRequestJoin = community.IsRequestJoin
	}

	return communityResponse, nil
}

func (s *CommunityService) UpdateCommunity(userID, id uint64, req *request.UpdateCommunityRequest) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(id)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.UpdateCommunity: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (must be SUPER_ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(id, userID)
	if err != nil || role != constant.ROLE_SUPER_ADMIN {
		log.Printf("[Err] User does not have permission in CommunityService.UpdateCommunity: userID=%d, communityID=%d", userID, id)
		return fmt.Errorf("permission denied")
	}

	if err := s.communityRepo.UpdateCommunity(id, req); err != nil {
		log.Printf("[Err] Error updating community in CommunityService.UpdateCommunity: %v", err)
		return fmt.Errorf("failed to update community")
	}

	return nil
}

func (s *CommunityService) DeleteCommunity(userID, id uint64) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(id)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.DeleteCommunity: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (must be SUPER_ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(id, userID)
	if err != nil || role != constant.ROLE_SUPER_ADMIN {
		log.Printf("[Err] User does not have permission in CommunityService.DeleteCommunity: userID=%d, communityID=%d", userID, id)
		return fmt.Errorf("permission denied")
	}

	if err := s.communityRepo.DeleteCommunity(id); err != nil {
		log.Printf("[Err] Error deleting community in CommunityService.DeleteCommunity: %v", err)
		return fmt.Errorf("failed to delete community")
	}

	return nil
}

func (s *CommunityService) JoinCommunity(userID, communityID uint64) error {
	// Check if community exists
	community, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.JoinCommunity: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if already subscribed
	isSubscribed, err := s.subscriptionRepo.IsUserSubscribed(userID, communityID)
	if err != nil {
		log.Printf("[Err] Error checking subscription in CommunityService.JoinCommunity: %v", err)
		return fmt.Errorf("failed to check subscription")
	}

	if isSubscribed {
		return fmt.Errorf("already subscribed to this community")
	}

	subscriptionStatus := constant.SUBSCRIPTION_STATUS_APPROVED
	if community.RequiresMemberApproval {
		subscriptionStatus = constant.SUBSCRIPTION_STATUS_PENDING
	}

	// Create subscription
	subscription := &model.Subscription{
		UserID:       userID,
		CommunityID:  communityID,
		SubscribedAt: time.Now(),
		Status:       subscriptionStatus,
	}

	if err := s.subscriptionRepo.CreateSubscription(subscription); err != nil {
		log.Printf("[Err] Error creating subscription in CommunityService.JoinCommunity: %v", err)
		return fmt.Errorf("failed to join community")
	}

	// Create bot task for interest score
	if subscriptionStatus == constant.SUBSCRIPTION_STATUS_APPROVED {
		go func(userID, communityID uint64) {
			if err := s.botTaskService.CreateInterestScoreTask(userID, communityID, constant.INTEREST_ACTION_JOIN_COMMUNITY, nil); err != nil {
				log.Printf("[Err] Error creating interest score task in goroutine (JoinCommunity): %v", err)
			}
		}(userID, communityID)
	}

	return nil
}

func (s *CommunityService) UnjoinCommunity(userID, communityID uint64) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.UnjoinCommunity: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user is subscribed
	isSubscribed, err := s.subscriptionRepo.IsUserSubscribed(userID, communityID)
	if err != nil {
		log.Printf("[Err] Error checking subscription in CommunityService.UnjoinCommunity: %v", err)
		return fmt.Errorf("failed to check subscription")
	}

	if !isSubscribed {
		return fmt.Errorf("not subscribed to this community")
	}

	// Check if user is a moderator
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err == nil && role != "" {
		if err := s.communityModeratorRepo.DeleteModerator(communityID, userID); err != nil {
			log.Printf("[Err] Error deleting moderator in CommunityService.UnjoinCommunity: %v", err)
			return fmt.Errorf("failed to remove moderator role")
		}
	}

	// Delete subscription
	if err := s.subscriptionRepo.DeleteSubscription(userID, communityID); err != nil {
		log.Printf("[Err] Error deleting subscription in CommunityService.UnjoinCommunity: %v", err)
		return fmt.Errorf("failed to leave community")
	}

	// Create bot task for interest score
	go func(userID, communityID uint64) {
		if err := s.botTaskService.CreateInterestScoreTask(userID, communityID, constant.INTEREST_ACTION_LEAVE_COMMUNITY, nil); err != nil {
			log.Printf("[Err] Error creating interest score task in goroutine (UnjoinCommunity): %v", err)
		}
	}(userID, communityID)

	return nil
}

func (s *CommunityService) GetCommunities(page, limit int, userID *uint64) ([]*response.CommunityListResponse, *response.Pagination, error) {
	communities, total, err := s.communityRepo.GetCommunities(page, limit, userID)
	if err != nil {
		log.Printf("[Err] Error getting communities in CommunityService.GetCommunities: %v", err)
		return nil, nil, fmt.Errorf("failed to get communities")
	}

	communityResponses := make([]*response.CommunityListResponse, len(communities))
	for i, community := range communities {
		resp := response.NewCommunityListResponse(community)
		resp.TotalMembers = community.MemberCount

		if community.IsSubscribed != nil {
			resp.IsFollow = community.IsSubscribed
		}

		communityResponses[i] = resp
	}

	// Set pagination
	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/communities?page=%d&limit=%d", page+1, limit)
	}

	return communityResponses, pagination, nil
}

func (s *CommunityService) SearchCommunitiesByName(name string, sortBy string, page, limit int, userID *uint64) ([]*response.CommunityListResponse, *response.Pagination, error) {
	// Validate sortBy
	if sortBy != constant.SORT_NEWEST && sortBy != constant.SORT_MEMBER_COUNT {
		sortBy = constant.SORT_NEWEST
	}

	communities, total, err := s.communityRepo.SearchCommunitiesByName(name, sortBy, page, limit, userID)
	if err != nil {
		log.Printf("[Err] Error searching communities in CommunityService.SearchCommunitiesByName: %v", err)
		return nil, nil, fmt.Errorf("failed to search communities")
	}

	communityResponses := make([]*response.CommunityListResponse, len(communities))
	for i, community := range communities {
		resp := response.NewCommunityListResponse(community)
		resp.TotalMembers = community.MemberCount

		if community.IsSubscribed != nil {
			resp.IsFollow = community.IsSubscribed
		}

		communityResponses[i] = resp
	}

	// Get pagination
	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/communities/search?name=%s&sortBy=%s&page=%d&limit=%d", name, sortBy, page+1, limit)
	}

	return communityResponses, pagination, nil
}

func (s *CommunityService) FilterCommunities(sortBy string, isPrivate *bool, topics []string, page, limit int, userID *uint64) ([]*response.CommunityListResponse, *response.Pagination, error) {
	// Validate sortBy
	if sortBy != constant.SORT_NEWEST && sortBy != constant.SORT_MEMBER_COUNT {
		sortBy = constant.SORT_NEWEST
	}

	communities, total, err := s.communityRepo.FilterCommunities(sortBy, isPrivate, topics, page, limit, userID)
	if err != nil {
		log.Printf("[Err] Error filtering communities in CommunityService.FilterCommunities: %v", err)
		return nil, nil, fmt.Errorf("failed to filter communities")
	}

	communityResponses := make([]*response.CommunityListResponse, len(communities))
	for i, community := range communities {
		resp := response.NewCommunityListResponse(community)
		resp.TotalMembers = community.MemberCount

		if community.IsSubscribed != nil {
			resp.IsFollow = community.IsSubscribed
		}

		communityResponses[i] = resp
	}

	// Get pagination
	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		nextURL := fmt.Sprintf("/api/v1/communities/filter?sortBy=%s&page=%d&limit=%d", sortBy, page+1, limit)
		if isPrivate != nil {
			nextURL += fmt.Sprintf("&isPrivate=%t", *isPrivate)
		}
		if len(topics) > 0 {
			nextURL += fmt.Sprintf("&topics=%s", strings.Join(topics, ","))
		}
		pagination.NextURL = nextURL
	}

	return communityResponses, pagination, nil
}

func (s *CommunityService) GetCommunityMembers(userID, communityID uint64, sortBy, searchName, status string, page, limit int) ([]*response.MemberListResponse, *response.Pagination, error) {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.GetCommunityMembers: %v", err)
		return nil, nil, fmt.Errorf("community not found")
	}

	// Check if user has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.GetCommunityMembers: userID=%d, communityID=%d", userID, communityID)
		return nil, nil, fmt.Errorf("permission denied")
	}

	// Validate sortBy
	if sortBy != constant.SORT_NEWEST && sortBy != constant.SORT_OLDEST && sortBy != constant.SORT_KARMA {
		sortBy = constant.SORT_NEWEST
	}

	// Default status is 'approved' if not specified
	if status == "" {
		status = constant.SUBSCRIPTION_STATUS_APPROVED
	}

	subscriptions, total, err := s.subscriptionRepo.GetCommunityMembers(communityID, sortBy, searchName, status, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting community members in CommunityService.GetCommunityMembers: %v", err)
		return nil, nil, fmt.Errorf("failed to get community members")
	}

	memberResponses := make([]*response.MemberListResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		// Default role is "user"
		role := constant.ROLE_USER
		if subscription.ModeratorRole != nil && *subscription.ModeratorRole != "" {
			role = *subscription.ModeratorRole
		}

		memberResponses[i] = response.NewMemberListResponse(subscription.User, subscription.SubscribedAt, role, subscription.Status, subscription.IsBannedBefore)
	}

	// Set pagination
	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		nextURL := fmt.Sprintf("/api/v1/communities/%d/members?sortBy=%s&status=%s&page=%d&limit=%d", communityID, sortBy, status, page+1, limit)
		if searchName != "" {
			nextURL += fmt.Sprintf("&search=%s", searchName)
		}
		pagination.NextURL = nextURL
	}

	return memberResponses, pagination, nil
}

func (s *CommunityService) RemoveMember(userID, communityID, memberID uint64) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.RemoveMember: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (must be SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.RemoveMember: userID=%d, communityID=%d", userID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Check if member is subscribed
	isSubscribed, err := s.subscriptionRepo.IsUserSubscribed(memberID, communityID)
	if err != nil {
		log.Printf("[Err] Error checking subscription in CommunityService.RemoveMember: %v", err)
		return fmt.Errorf("failed to check subscription")
	}
	if !isSubscribed {
		return fmt.Errorf("member not found in this community")
	}

	// Remove member
	if err := s.subscriptionRepo.DeleteSubscription(memberID, communityID); err != nil {
		log.Printf("[Err] Error removing member in CommunityService.RemoveMember: %v", err)
		return fmt.Errorf("failed to remove member")
	}

	return nil
}

func (s *CommunityService) UpdateMemberRole(adminUserID, communityID, targetUserID uint64, role string) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.UpdateMemberRole: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if admin has permission (must be SUPER_ADMIN)
	adminRole, err := s.communityModeratorRepo.GetModeratorRole(communityID, adminUserID)
	if err != nil || adminRole != constant.ROLE_SUPER_ADMIN {
		log.Printf("[Err] User does not have permission in CommunityService.UpdateMemberRole: userID=%d, communityID=%d", adminUserID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Check if target user is a member
	isSubscribed, err := s.subscriptionRepo.IsUserSubscribed(targetUserID, communityID)
	if err != nil {
		log.Printf("[Err] Error checking subscription in CommunityService.UpdateMemberRole: %v", err)
		return fmt.Errorf("failed to check subscription")
	}
	if !isSubscribed {
		return fmt.Errorf("user is not a member of this community")
	}

	// If role is "user", remove from moderators
	if role == constant.ROLE_USER {
		if err := s.communityModeratorRepo.DeleteModerator(communityID, targetUserID); err != nil {
			log.Printf("[Err] Error removing moderator in CommunityService.UpdateMemberRole: %v", err)
			return fmt.Errorf("failed to remove moderator role")
		}
		return nil
	}

	// If role is "admin", upsert moderator
	if role == constant.ROLE_ADMIN {
		moderator := &model.CommunityModerator{
			CommunityID: communityID,
			UserID:      targetUserID,
			Role:        constant.ROLE_ADMIN,
			JoinedAt:    time.Now(),
		}
		if err := s.communityModeratorRepo.UpsertModerator(moderator); err != nil {
			log.Printf("[Err] Error upserting moderator in CommunityService.UpdateMemberRole: %v", err)
			return fmt.Errorf("failed to update moderator role")
		}
		return nil
	}

	return fmt.Errorf("invalid role")
}

func (s *CommunityService) GetUserRoleInCommunity(userID, communityID uint64) (string, error) {
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil {
		log.Printf("[Err] Error getting user role in CommunityService.GetUserRoleInCommunity: %v", err)
		return "", fmt.Errorf("failed to get user role")
	}
	return role, nil
}

func (s *CommunityService) VerifyCommunityName(name string) (bool, error) {
	exists, err := s.communityRepo.IsCommunityNameExists(name)
	if err != nil {
		log.Printf("[Err] Error checking community name in CommunityService.VerifyCommunityName: %v", err)
		return false, fmt.Errorf("failed to check community name")
	}

	return !exists, nil
}

func (s *CommunityService) GetCommunityPostsForModerator(userID, communityID uint64, status, searchTitle string, page, limit int) ([]*response.CommunityPostListResponse, *response.Pagination, error) {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.GetCommunityPostsForModerator: %v", err)
		return nil, nil, fmt.Errorf("community not found")
	}

	// Check if user has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.GetCommunityPostsForModerator: userID=%d, communityID=%d", userID, communityID)
		return nil, nil, fmt.Errorf("permission denied")
	}

	posts, total, err := s.postRepo.GetCommunityPostsForModerator(communityID, status, searchTitle, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting community posts for moderator in CommunityService.GetCommunityPostsForModerator: %v", err)
		return nil, nil, fmt.Errorf("failed to get posts")
	}

	postResponses := make([]*response.CommunityPostListResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = response.NewCommunityPostListResponse(post)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		nextURL := fmt.Sprintf("/api/v1/communities/%d/manage/posts?page=%d&limit=%d", communityID, page+1, limit)
		if status != "" {
			nextURL += fmt.Sprintf("&status=%s", status)
		}
		if searchTitle != "" {
			nextURL += fmt.Sprintf("&search=%s", searchTitle)
		}
		pagination.NextURL = nextURL
	}

	return postResponses, pagination, nil
}

func (s *CommunityService) UpdatePostStatusByModerator(userID, communityID, postID uint64, status string) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.UpdatePostStatusByModerator: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.UpdatePostStatusByModerator: userID=%d, communityID=%d", userID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Get post
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("[Err] Post not found in CommunityService.UpdatePostStatusByModerator: %v", err)
		return fmt.Errorf("post not found")
	}

	// Verify post belongs to this community
	if post.CommunityID != communityID {
		log.Printf("[Err] Post does not belong to community in CommunityService.UpdatePostStatusByModerator: postID=%d, communityID=%d", postID, communityID)
		return fmt.Errorf("post not found in this community")
	}

	if err := s.postRepo.UpdatePostStatus(postID, status); err != nil {
		log.Printf("[Err] Error updating post status in CommunityService.UpdatePostStatusByModerator: %v", err)
		return fmt.Errorf("failed to update post status")
	}

	// Send notification to post author
	if status == constant.POST_STATUS_APPROVED || status == constant.POST_STATUS_REJECTED {
		go func(authorID uint64, postID uint64, postTitle string, status string) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Panic] Recovered in CommunityService.UpdatePostStatusByModerator notification: %v", r)
				}
			}()

			statusStr := ""
			if status == constant.POST_STATUS_APPROVED {
				statusStr = "approved"
			} else {
				statusStr = "rejected"
			}

			notifPayload := payload.PostStatusNotificationPayload{
				PostID: postID,
				Status: statusStr,
			}

			s.notificationService.CreateNotification(authorID, constant.NOTIFICATION_ACTION_POST_STATUS_UPDATED, notifPayload)
		}(post.AuthorID, postID, post.Title, status)
	}

	return nil
}

func (s *CommunityService) DeletePostByModerator(userID, communityID, postID uint64) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.DeletePostByModerator: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.DeletePostByModerator: userID=%d, communityID=%d", userID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Get post
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("[Err] Post not found in CommunityService.DeletePostByModerator: %v", err)
		return fmt.Errorf("post not found")
	}

	// Verify post belongs to this community
	if post.CommunityID != communityID {
		log.Printf("[Err] Post does not belong to community in CommunityService.DeletePostByModerator: postID=%d, communityID=%d", postID, communityID)
		return fmt.Errorf("post not found in this community")
	}

	// Delete post
	if err := s.postRepo.DeletePost(postID); err != nil {
		log.Printf("[Err] Error deleting post in CommunityService.DeletePostByModerator: %v", err)
		return fmt.Errorf("failed to delete post")
	}

	// Send notification to post author
	go func(authorID uint64, postID uint64) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Recovered in CommunityService.DeletePostByModerator notification: %v", r)
			}
		}()

		notifPayload := map[string]interface{}{
			"postId": postID,
		}

		s.notificationService.CreateNotification(authorID, constant.NOTIFICATION_ACTION_POST_DELETED, notifPayload)
	}(post.AuthorID, postID)

	return nil
}

func (s *CommunityService) DeleteCommentByModerator(userID, communityID, commentID uint64) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.DeleteCommentByModerator: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.DeleteCommentByModerator: userID=%d, communityID=%d", userID, communityID)
		return fmt.Errorf("permission denied")
	}

	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		log.Printf("[Err] Comment not found in CommunityService.DeleteCommentByModerator: %v", err)
		return fmt.Errorf("comment not found")
	}

	post, err := s.postRepo.GetPostByID(comment.PostID)
	if err != nil {
		log.Printf("[Err] Post not found in CommunityService.DeleteCommentByModerator: %v", err)
		return fmt.Errorf("post not found")
	}

	// Verify post belongs to this community
	if post.CommunityID != communityID {
		log.Printf("[Err] Comment does not belong to community in CommunityService.DeleteCommentByModerator: commentID=%d, communityID=%d", commentID, communityID)
		return fmt.Errorf("comment not found in this community")
	}

	if err := s.commentRepo.DeleteComment(commentID, comment.ParentCommentID); err != nil {
		log.Printf("[Err] Error deleting comment in CommunityService.DeleteCommentByModerator: %v", err)
		return fmt.Errorf("failed to delete comment")
	}

	// Send notification to comment author
	go func(authorID uint64, commentID uint64, postID uint64) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Recovered in CommunityService.DeleteCommentByModerator notification: %v", r)
			}
		}()

		notifPayload := payload.CommentDeletedNotificationPayload{
			CommentID: commentID,
			PostID:    postID,
		}

		s.notificationService.CreateNotification(authorID, constant.NOTIFICATION_ACTION_COMMENT_DELETED, notifPayload)
	}(comment.AuthorID, commentID, comment.PostID)

	return nil
}

func (s *CommunityService) GetCommunityPostReports(userID, communityID uint64, page, limit int) ([]*response.PostReportResponse, *response.Pagination, error) {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.GetCommunityPostReports: %v", err)
		return nil, nil, fmt.Errorf("community not found")
	}

	// Check if user has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.GetCommunityPostReports: userID=%d, communityID=%d", userID, communityID)
		return nil, nil, fmt.Errorf("permission denied")
	}

	reports, total, err := s.postReportRepo.GetPostReportsByCommunityID(communityID, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting community post reports in CommunityService.GetCommunityPostReports: %v", err)
		return nil, nil, fmt.Errorf("failed to get post reports")
	}

	// Group reports by post
	reportMap := make(map[uint64]*response.PostReportResponse)
	for _, report := range reports {
		// Skip reports where post is nil (shouldn't happen but safety check)
		if report.Post == nil {
			log.Printf("[Warning] Post is nil for report ID %d, PostID %d", report.ID, report.PostID)
			continue
		}

		if reportMap[report.PostID] == nil {
			reportMap[report.PostID] = response.NewPostReportResponse(
				report.Post.ID,
				report.Post.Title,
				report.Post.Author,
			)
			// Set the first report ID as the response ID
			reportMap[report.PostID].ID = report.ID
		}

		reporterInfo := response.NewReporterInfo(report.Reporter, report.Reasons, report.Note)
		reportMap[report.PostID].Reporters = append(reportMap[report.PostID].Reporters, reporterInfo)
		reportMap[report.PostID].TotalReports++
	}

	// Convert map to slice
	reportResponses := make([]*response.PostReportResponse, 0, len(reportMap))
	for _, reportResp := range reportMap {
		reportResponses = append(reportResponses, reportResp)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/communities/%d/manage/reports?page=%d&limit=%d", communityID, page+1, limit)
	}

	return reportResponses, pagination, nil
}

func (s *CommunityService) DeletePostReport(userID, communityID, reportID uint64) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.DeletePostReport: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.DeletePostReport: userID=%d, communityID=%d", userID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Delete report
	if err := s.postReportRepo.DeletePostReport(reportID); err != nil {
		log.Printf("[Err] Error deleting post report in CommunityService.DeletePostReport: %v", err)
		return fmt.Errorf("failed to delete post report")
	}

	return nil
}

func (s *CommunityService) GetCommunityCommentReports(userID, communityID uint64, page, limit int) ([]*response.CommentReportResponse, *response.Pagination, error) {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.GetCommunityCommentReports: %v", err)
		return nil, nil, fmt.Errorf("community not found")
	}

	// Check if user has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.GetCommunityCommentReports: userID=%d, communityID=%d", userID, communityID)
		return nil, nil, fmt.Errorf("permission denied")
	}

	reports, total, err := s.commentReportRepo.GetCommentReportsByCommunityID(communityID, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting community comment reports in CommunityService.GetCommunityCommentReports: %v", err)
		return nil, nil, fmt.Errorf("failed to get comment reports")
	}

	// Group reports by comment
	reportMap := make(map[uint64]*response.CommentReportResponse)
	for _, report := range reports {
		if reportMap[report.CommentID] == nil {
			postTitle := ""
			postID := uint64(0)
			if report.Comment.Post != nil {
				postTitle = report.Comment.Post.Title
				postID = report.Comment.Post.ID
			}

			reportMap[report.CommentID] = response.NewCommentReportResponse(
				report.Comment.ID,
				report.Comment.Content,
				report.Comment.Author,
				postID,
				postTitle,
			)
			// Set the first report ID as the response ID
			reportMap[report.CommentID].ID = report.ID
		}

		reporterInfo := response.NewReporterInfo(report.Reporter, report.Reasons, report.Note)
		reportMap[report.CommentID].Reporters = append(reportMap[report.CommentID].Reporters, reporterInfo)
		reportMap[report.CommentID].TotalReports++
	}

	// Convert map to slice
	reportResponses := make([]*response.CommentReportResponse, 0, len(reportMap))
	for _, reportResp := range reportMap {
		reportResponses = append(reportResponses, reportResp)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/communities/%d/manage/comment-reports?page=%d&limit=%d", communityID, page+1, limit)
	}

	return reportResponses, pagination, nil
}

func (s *CommunityService) DeleteCommentReport(userID, communityID, reportID uint64) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.DeleteCommentReport: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.DeleteCommentReport: userID=%d, communityID=%d", userID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Delete report
	if err := s.commentReportRepo.DeleteCommentReport(reportID); err != nil {
		log.Printf("[Err] Error deleting comment report in CommunityService.DeleteCommentReport: %v", err)
		return fmt.Errorf("failed to delete comment report")
	}

	return nil
}

func (s *CommunityService) GetAllTopics(search *string) ([]*response.TopicResponse, error) {
	topics, err := s.topicRepo.GetAllTopics(search)
	if err != nil {
		log.Printf("[Err] Error getting topics in CommunityService.GetAllTopics: %v", err)
		return nil, err
	}

	topicResponses := make([]*response.TopicResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = &response.TopicResponse{
			ID:   topic.ID,
			Name: topic.Name,
		}
	}

	return topicResponses, nil
}

func (s *CommunityService) UpdateRequiresPostApproval(userID, communityID uint64, req *request.UpdateRequiresPostApprovalRequest) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.UpdateRequiresPostApproval: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (must be SUPER_ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || role != constant.ROLE_SUPER_ADMIN {
		log.Printf("[Err] User does not have permission in CommunityService.UpdateRequiresPostApproval: userID=%d, communityID=%d", userID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Update requires post approval
	if err := s.communityRepo.UpdateRequiresPostApproval(communityID, req.RequiresPostApproval); err != nil {
		log.Printf("[Err] Error updating requires post approval in CommunityService.UpdateRequiresPostApproval: %v", err)
		return fmt.Errorf("failed to update requires post approval")
	}

	return nil
}

func (s *CommunityService) UpdateRequiresMemberApproval(userID, communityID uint64, req *request.UpdateRequiresMemberApprovalRequest) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.UpdateRequiresMemberApproval: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if user has permission (must be SUPER_ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, userID)
	if err != nil || role != constant.ROLE_SUPER_ADMIN {
		log.Printf("[Err] User does not have permission in CommunityService.UpdateRequiresMemberApproval: userID=%d, communityID=%d", userID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Update requires member approval
	if err := s.communityRepo.UpdateRequiresMemberApproval(communityID, req.RequiresMemberApproval); err != nil {
		log.Printf("[Err] Error updating requires member approval in CommunityService.UpdateRequiresMemberApproval: %v", err)
		return fmt.Errorf("failed to update requires member approval")
	}

	return nil
}

func (s *CommunityService) UpdateSubscriptionStatus(moderatorUserID, communityID, targetUserID uint64, status string) error {
	// Check if community exists
	community, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.UpdateSubscriptionStatus: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if moderator has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, moderatorUserID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.UpdateSubscriptionStatus: userID=%d, communityID=%d", moderatorUserID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Check if subscription exists
	isSubscribed, err := s.subscriptionRepo.IsUserSubscribed(targetUserID, communityID)
	if err != nil || !isSubscribed {
		log.Printf("[Err] Subscription not found in CommunityService.UpdateSubscriptionStatus: %v", err)
		return fmt.Errorf("subscription not found")
	}

	switch status {
	case constant.SUBSCRIPTION_STATUS_APPROVED:
		// Update status to approved
		if err := s.subscriptionRepo.UpdateSubscriptionStatus(targetUserID, communityID, status); err != nil {
			log.Printf("[Err] Error updating subscription status in CommunityService.UpdateSubscriptionStatus: %v", err)
			return fmt.Errorf("failed to approve subscription")
		}

		// Create bot task for interest score
		go func(userID, communityID uint64) {
			if err := s.botTaskService.CreateInterestScoreTask(userID, communityID, constant.INTEREST_ACTION_JOIN_COMMUNITY, nil); err != nil {
				log.Printf("[Err] Error creating interest score task in goroutine (UpdateSubscriptionStatus): %v", err)
			}
		}(targetUserID, communityID)

		// Send notification to user
		go func(targetUserID, communityID uint64, communityName string) {
			notificationPayload := payload.SubscriptionStatusNotificationPayload{
				CommunityID:   communityID,
				CommunityName: communityName,
				Status:        "approved",
			}

			if err := s.notificationService.CreateNotification(targetUserID, constant.NOTIFICATION_ACTION_SUBSCRIPTION_STATUS_UPDATED, notificationPayload); err != nil {
				log.Printf("[Err] Error sending notification in goroutine (UpdateSubscriptionStatus-Approved): %v", err)
			}
		}(targetUserID, communityID, community.Name)

	case constant.SUBSCRIPTION_STATUS_REJECTED:
		// Delete subscription
		if err := s.subscriptionRepo.DeleteSubscription(targetUserID, communityID); err != nil {
			log.Printf("[Err] Error deleting subscription in CommunityService.UpdateSubscriptionStatus: %v", err)
			return fmt.Errorf("failed to reject subscription")
		}

		// Send notification to user
		go func(targetUserID, communityID uint64, communityName string) {
			notificationPayload := payload.SubscriptionStatusNotificationPayload{
				CommunityID:   communityID,
				CommunityName: communityName,
				Status:        "rejected",
			}

			if err := s.notificationService.CreateNotification(targetUserID, constant.NOTIFICATION_ACTION_SUBSCRIPTION_STATUS_UPDATED, notificationPayload); err != nil {
				log.Printf("[Err] Error sending notification in goroutine (UpdateSubscriptionStatus-Rejected): %v", err)
			}
		}(targetUserID, communityID, community.Name)

	default:
		return fmt.Errorf("invalid status")
	}

	return nil

}

func (s *CommunityService) BanUser(moderatorID, communityID uint64, req *request.BanUserRequest) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.BanUser: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if moderator has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, moderatorID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.BanUser: moderatorID=%d, communityID=%d", moderatorID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Validate restriction type and expiry date
	if req.RestrictionType == constant.RESTRICTION_TEMPORARY_BAN {
		if req.ExpiresAt == nil {
			return fmt.Errorf("temporary ban requires expiry date")
		}
		if req.ExpiresAt.Before(time.Now()) {
			return fmt.Errorf("expiry date must be in the future")
		}
	} else if req.RestrictionType == constant.RESTRICTION_PERMANENT_BAN || req.RestrictionType == constant.RESTRICTION_WARNING {
		if req.ExpiresAt != nil {
			return fmt.Errorf("permanent ban and warning cannot have expiry date")
		}
	} else {
		return fmt.Errorf("invalid restriction type")
	}

	// Create user restriction
	restriction := &model.UserRestriction{
		UserID:          req.UserID,
		CommunityID:     communityID,
		RestrictionType: req.RestrictionType,
		Reason:          req.Reason,
		IssuedBy:        moderatorID,
		ExpiresAt:       req.ExpiresAt,
		CreatedAt:       time.Now(),
	}

	if err := s.userRestrictionRepo.CreateRestriction(restriction); err != nil {
		log.Printf("[Err] Error creating user restriction in CommunityService.BanUser: %v", err)
		return fmt.Errorf("failed to ban user")
	}

	// Get community name for notification
	community, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Failed to get community for notification in CommunityService.BanUser: %v", err)
	} else {
		// Send notification to banned user
		go func(targetUserID, communityID uint64, communityName, restrictionType, reason string, expiresAt *time.Time) {
			expiresAtStr := ""
			if expiresAt != nil {
				expiresAtStr = expiresAt.Format("2006-01-02 15:04:05")
			}

			notificationPayload := payload.UserBanNotificationPayload{
				CommunityID:     communityID,
				CommunityName:   communityName,
				RestrictionType: restrictionType,
				Reason:          reason,
				ExpiresAt:       expiresAtStr,
			}

			if err := s.notificationService.CreateNotification(targetUserID, constant.NOTIFICATION_ACTION_USER_BANNED, notificationPayload); err != nil {
				log.Printf("[Err] Error sending notification in goroutine (BanUser): %v", err)
			}
		}(req.UserID, communityID, community.Name, req.RestrictionType, req.Reason, req.ExpiresAt)
	}

	return nil
}

func (s *CommunityService) GetUserRestrictionHistory(moderatorID, communityID, targetUserID uint64, page, limit int) ([]*response.UserRestrictionResponse, *response.Pagination, error) {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.GetUserRestrictionHistory: %v", err)
		return nil, nil, fmt.Errorf("community not found")
	}

	// Check if moderator has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, moderatorID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.GetUserRestrictionHistory: moderatorID=%d, communityID=%d", moderatorID, communityID)
		return nil, nil, fmt.Errorf("permission denied")
	}

	// Get restriction history
	restrictions, total, err := s.userRestrictionRepo.GetUserRestrictionHistory(targetUserID, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting user restriction history in CommunityService.GetUserRestrictionHistory: %v", err)
		return nil, nil, fmt.Errorf("failed to get restriction history")
	}

	// Filter restrictions for this community only
	communityRestrictions := make([]*response.UserRestrictionResponse, 0)
	for _, restriction := range restrictions {
		if restriction.CommunityID == communityID {
			communityRestrictions = append(communityRestrictions, response.NewUserRestrictionResponse(restriction))
		}
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/communities/%d/manage/restrictions/user/%d?page=%d&limit=%d", communityID, targetUserID, page+1, limit)
	}

	return communityRestrictions, pagination, nil
}

func (s *CommunityService) RemoveUserRestriction(moderatorID, communityID, restrictionID uint64) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.RemoveUserRestriction: %v", err)
		return fmt.Errorf("community not found")
	}

	// Check if moderator has permission (SUPER_ADMIN or ADMIN)
	role, err := s.communityModeratorRepo.GetModeratorRole(communityID, moderatorID)
	if err != nil || (role != constant.ROLE_SUPER_ADMIN && role != constant.ROLE_ADMIN) {
		log.Printf("[Err] User does not have permission in CommunityService.RemoveUserRestriction: moderatorID=%d, communityID=%d", moderatorID, communityID)
		return fmt.Errorf("permission denied")
	}

	// Delete restriction
	if err := s.userRestrictionRepo.DeleteRestriction(restrictionID); err != nil {
		log.Printf("[Err] Error deleting user restriction in CommunityService.RemoveUserRestriction: %v", err)
		return fmt.Errorf("failed to remove restriction")
	}

	return nil
}
