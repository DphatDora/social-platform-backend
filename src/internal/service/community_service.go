package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"time"
)

type CommunityService struct {
	communityRepo repository.CommunityRepository
}

func NewCommunityService(communityRepo repository.CommunityRepository) *CommunityService {
	return &CommunityService{
		communityRepo: communityRepo,
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

	return nil
}

func (s *CommunityService) GetCommunityByID(id uint64) (*response.CommunityDetailResponse, error) {
	community, err := s.communityRepo.GetCommunityByID(id)
	if err != nil {
		log.Printf("[Err] Error getting community by ID in CommunityService.GetCommunityByID: %v", err)
		return nil, fmt.Errorf("community not found")
	}

	return response.NewCommunityDetailResponse(community), nil
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
