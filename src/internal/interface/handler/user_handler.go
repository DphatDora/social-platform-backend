package handler

import (
	"log"
	"net/http"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/internal/service"
	"strconv"
	"strings"

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

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Printf("[Err] UserID not found in context in UserHandler.ChangePassword")
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		log.Printf("[Err] Invalid userID type in context in UserHandler.ChangePassword")
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to retrieve user information",
		})
		return
	}

	var changePasswordReq request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&changePasswordReq); err != nil {
		log.Printf("[Err] Error binding JSON in UserHandler.ChangePassword: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.userService.ChangePassword(userID, &changePasswordReq); err != nil {
		log.Printf("[Err] Error changing password in UserHandler.ChangePassword: %v", err)

		if strings.Contains(err.Error(), "incorrect") {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Old password is incorrect",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to change password",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

func (h *UserHandler) GetUserConfig(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Printf("[Err] UserID not found in context in UserHandler.GetUserConfig")
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		log.Printf("[Err] Invalid userID type in context in UserHandler.GetUserConfig")
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to retrieve user config",
		})
		return
	}

	// Get user config
	userConfig, err := h.userService.GetUserConfig(userID)
	if err != nil {
		log.Printf("[Err] Error getting user config in UserHandler.GetUserConfig: %v", err)
		c.JSON(http.StatusNotFound, response.APIResponse{
			Success: false,
			Message: "Failed to get user config",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "User config retrieved successfully",
		Data:    userConfig,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid user ID in UserHandler.GetUserByID: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	userProfile, err := h.userService.GetUserProfile(userID)
	if err != nil {
		log.Printf("[Err] Error getting user profile in UserHandler.GetUserByID: %v", err)
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
