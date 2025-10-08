package handler

import (
	"log"
	"net/http"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Printf("[Err] UserID not found in context in UserHandler.GetCurrentUser")
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		log.Printf("[Err] Invalid userID type in context in UserHandler.GetCurrentUser")
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to retrieve user information",
		})
		return
	}

	// Get user profile
	userProfile, err := h.userService.GetUserProfile(userID)
	if err != nil {
		log.Printf("[Err] Error getting user profile in UserHandler.GetCurrentUser: %v", err)
		c.JSON(http.StatusNotFound, response.APIResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		Data:    userProfile,
	})
}

func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Printf("[Err] UserID not found in context in UserHandler.UpdateUserProfile")
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		log.Printf("[Err] Invalid userID type in context in UserHandler.UpdateUserProfile")
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to retrieve user information",
		})
		return
	}

	var updateReq request.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		log.Printf("[Err] Error binding JSON in UserHandler.UpdateUserProfile: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	if err := h.userService.UpdateUserProfile(userID, &updateReq); err != nil {
		log.Printf("[Err] Error updating user profile in UserHandler.UpdateUserProfile: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update user profile",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "User profile updated successfully",
	})
}
