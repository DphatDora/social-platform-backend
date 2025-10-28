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
}

func NewCommunityService(
	communityRepo repository.CommunityRepository,
	subscriptionRepo repository.SubscriptionRepository,
	communityModeratorRepo repository.CommunityModeratorRepository,
) *CommunityService {
	return &CommunityService{
		communityRepo:          communityRepo,
		subscriptionRepo:       subscriptionRepo,
		communityModeratorRepo: communityModeratorRepo,
	}
}

func (s *CommunityService) CreateCommunity(userID uint64, req *request.CreateCommunityRequest) error {
	community := &model.Community{
		Name:             req.Name,
		ShortDescription: req.ShortDescription,
		Description:      req.Description,
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
