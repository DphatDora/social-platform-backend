package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/response"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserProfile(userID uint64) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("[Err] Error getting user by ID in UserService.GetUserProfile: %v", err)
		return nil, fmt.Errorf("user not found")
	}

	userProfile := response.NewUserProfileResponse(user)
	return userProfile, nil
}
