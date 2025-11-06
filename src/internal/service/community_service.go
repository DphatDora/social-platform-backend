package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
	"time"
)

type CommunityService struct {
	communityRepo          repository.CommunityRepository
	subscriptionRepo       repository.SubscriptionRepository
	communityModeratorRepo repository.CommunityModeratorRepository
	postRepo               repository.PostRepository
	postReportRepo         repository.PostReportRepository
	notificationService    *NotificationService
}

func NewCommunityService(
	communityRepo repository.CommunityRepository,
	subscriptionRepo repository.SubscriptionRepository,
	communityModeratorRepo repository.CommunityModeratorRepository,
	postRepo repository.PostRepository,
	postReportRepo repository.PostReportRepository,
	notificationService *NotificationService,
) *CommunityService {
	return &CommunityService{
		communityRepo:          communityRepo,
		subscriptionRepo:       subscriptionRepo,
		communityModeratorRepo: communityModeratorRepo,
		postRepo:               postRepo,
		postReportRepo:         postReportRepo,
		notificationService:    notificationService,
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

func (s *CommunityService) GetCommunityByID(id uint64) (*response.CommunityDetailResponse, error) {
	community, memberCount, err := s.communityRepo.GetCommunityWithMemberCount(id)
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

	communityResponse := response.NewCommunityDetailResponse(community)
	communityResponse.TotalMembers = memberCount
	communityResponse.Moderators = moderatorResponses

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
	_, err := s.communityRepo.GetCommunityByID(communityID)
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
		return fmt.Errorf("already joined this community")
	}

	// Create subscription
	subscription := &model.Subscription{
		UserID:       userID,
		CommunityID:  communityID,
		SubscribedAt: time.Now(),
	}

	if err := s.subscriptionRepo.CreateSubscription(subscription); err != nil {
		log.Printf("[Err] Error creating subscription in CommunityService.JoinCommunity: %v", err)
		return fmt.Errorf("failed to join community")
	}

	return nil
}

func (s *CommunityService) GetCommunities(page, limit int) ([]*response.CommunityListResponse, *response.Pagination, error) {
	communities, total, err := s.communityRepo.GetCommunities(page, limit)
	if err != nil {
		log.Printf("[Err] Error getting communities in CommunityService.GetCommunities: %v", err)
		return nil, nil, fmt.Errorf("failed to get communities")
	}

	communityResponses := make([]*response.CommunityListResponse, len(communities))
	for i, community := range communities {
		resp := response.NewCommunityListResponse(community)
		resp.TotalMembers = community.MemberCount
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

func (s *CommunityService) SearchCommunitiesByName(name string, page, limit int) ([]*response.CommunityListResponse, *response.Pagination, error) {
	communities, total, err := s.communityRepo.SearchCommunitiesByName(name, page, limit)
	if err != nil {
		log.Printf("[Err] Error searching communities in CommunityService.SearchCommunitiesByName: %v", err)
		return nil, nil, fmt.Errorf("failed to search communities")
	}

	communityResponses := make([]*response.CommunityListResponse, len(communities))
	for i, community := range communities {
		resp := response.NewCommunityListResponse(community)
		resp.TotalMembers = community.MemberCount
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
		pagination.NextURL = fmt.Sprintf("/api/v1/communities/search?name=%s&page=%d&limit=%d", name, page+1, limit)
	}

	return communityResponses, pagination, nil
}

func (s *CommunityService) FilterCommunities(sortBy string, isPrivate *bool, page, limit int) ([]*response.CommunityListResponse, *response.Pagination, error) {
	// Validate sortBy
	if sortBy != constant.SORT_NEWEST && sortBy != constant.SORT_MEMBER_COUNT {
		sortBy = constant.SORT_NEWEST
	}

	communities, total, err := s.communityRepo.FilterCommunities(sortBy, isPrivate, page, limit)
	if err != nil {
		log.Printf("[Err] Error filtering communities in CommunityService.FilterCommunities: %v", err)
		return nil, nil, fmt.Errorf("failed to filter communities")
	}

	communityResponses := make([]*response.CommunityListResponse, len(communities))
	for i, community := range communities {
		resp := response.NewCommunityListResponse(community)
		resp.TotalMembers = community.MemberCount
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
		pagination.NextURL = nextURL
	}

	return communityResponses, pagination, nil
}

func (s *CommunityService) GetCommunityMembers(userID, communityID uint64, sortBy, searchName string, page, limit int) ([]*response.MemberListResponse, *response.Pagination, error) {
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

	subscriptions, total, err := s.subscriptionRepo.GetCommunityMembers(communityID, sortBy, searchName, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting community members in CommunityService.GetCommunityMembers: %v", err)
		return nil, nil, fmt.Errorf("failed to get community members")
	}

	memberResponses := make([]*response.MemberListResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		memberResponses[i] = response.NewMemberListResponse(subscription.User, subscription.SubscribedAt)
	}

	// Set pagination
	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		nextURL := fmt.Sprintf("/api/v1/communities/%d/members?sortBy=%s&page=%d&limit=%d", communityID, sortBy, page+1, limit)
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

			action := ""
			if status == constant.POST_STATUS_APPROVED {
				action = constant.NOTIFICATION_ACTION_POST_APPROVED
			} else {
				action = constant.NOTIFICATION_ACTION_POST_REJECTED
			}

			notifPayload := map[string]interface{}{
				"postId": postID,
			}

			s.notificationService.CreateNotification(authorID, action, notifPayload)
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
		if reportMap[report.PostID] == nil {
			reportMap[report.PostID] = &response.PostReportResponse{
				PostID:    report.Post.ID,
				PostTitle: report.Post.Title,
				Author: response.AuthorInfo{
					ID:       report.Post.Author.ID,
					Username: report.Post.Author.Username,
					Avatar:   report.Post.Author.Avatar,
				},
				Reporters:    []response.ReporterInfo{},
				TotalReports: 0,
			}
		}

		reporterInfo := response.ReporterInfo{
			ID:       report.Reporter.ID,
			Username: report.Reporter.Username,
			Avatar:   report.Reporter.Avatar,
			Reasons:  report.Reasons,
			Note:     report.Note,
		}

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
