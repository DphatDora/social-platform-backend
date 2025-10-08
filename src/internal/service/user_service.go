package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/util"
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
