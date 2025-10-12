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

	communityResponse := response.NewCommunityDetailResponse(community)
	communityResponse.TotalMembers = memberCount

	return communityResponse, nil
}

func (s *CommunityService) UpdateCommunity(id uint64, req *request.UpdateCommunityRequest) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(id)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.UpdateCommunity: %v", err)
		return fmt.Errorf("community not found")
	}

	if err := s.communityRepo.UpdateCommunity(id, req); err != nil {
		log.Printf("[Err] Error updating community in CommunityService.UpdateCommunity: %v", err)
		return fmt.Errorf("failed to update community")
	}

	return nil
}

func (s *CommunityService) DeleteCommunity(id uint64) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(id)
	if err != nil {
		log.Printf("[Err] Community not found in CommunityService.DeleteCommunity: %v", err)
		return fmt.Errorf("community not found")
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
