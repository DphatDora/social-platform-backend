package handler

import (
	"log"
	"net/http"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/internal/service"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/util"
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
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in UserHandler.GetCurrentUser", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
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
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in UserHandler.UpdateUserProfile", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var updateReq request.UpdateUserProfileRequest
	if err = c.ShouldBindJSON(&updateReq); err != nil {
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
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in UserHandler.ChangePassword", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var changePasswordReq request.ChangePasswordRequest
	if err = c.ShouldBindJSON(&changePasswordReq); err != nil {
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
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in UserHandler.GetUserConfig", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
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

func (h *UserHandler) GetUserBadgeHistory(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid user ID in UserHandler.GetUserBadgeHistory: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	badgeHistory, err := h.userService.GetUserBadgeHistory(userID)
	if err != nil {
		log.Printf("[Err] Error getting user badge history in UserHandler.GetUserBadgeHistory: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get user badge history",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "User badge history retrieved successfully",
		Data:    badgeHistory,
	})
}

func (h *UserHandler) GetUserSavedPosts(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in UserHandler.GetUserSavedPosts", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	searchTitle := c.Query("search")
	var isFollowed *bool
	if followedParam := c.Query("isFollowed"); followedParam != "" {
		followed := followedParam == "true"
		isFollowed = &followed
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	savedPosts, pagination, err := h.userService.GetUserSavedPosts(userID, searchTitle, isFollowed, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting user saved posts in UserHandler.GetUserSavedPosts: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get saved posts",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Saved posts retrieved successfully",
		Data:       savedPosts,
		Pagination: pagination,
	})
}

func (h *UserHandler) CreateUserSavedPost(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in UserHandler.CreateUserSavedPost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var savedPostReq request.UserSavedPostRequest
	if err = c.ShouldBindJSON(&savedPostReq); err != nil {
		log.Printf("[Err] Error binding JSON in UserHandler.CreateUserSavedPost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	if err := h.userService.CreateUserSavedPost(userID, &savedPostReq); err != nil {
		log.Printf("[Err] Error creating user saved post in UserHandler.CreateUserSavedPost: %v", err)

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		if strings.Contains(err.Error(), "already saved") {
			c.JSON(http.StatusConflict, response.APIResponse{
				Success: false,
				Message: "Post already saved",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to save post",
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Post saved successfully",
	})
}

func (h *UserHandler) UpdateUserSavedPostFollowStatus(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in UserHandler.UpdateUserSavedPostFollowStatus", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	postIDParam := c.Param("postId")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in UserHandler.UpdateUserSavedPostFollowStatus: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var updateReq request.UpdateUserSavedPostRequest
	if err = c.ShouldBindJSON(&updateReq); err != nil {
		log.Printf("[Err] Error binding JSON in UserHandler.UpdateUserSavedPostFollowStatus: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	if err := h.userService.UpdateUserSavedPostFollowStatus(userID, postID, &updateReq); err != nil {
		log.Printf("[Err] Error updating user saved post follow status in UserHandler.UpdateUserSavedPostFollowStatus: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update follow status",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Follow status updated successfully",
	})
}

func (h *UserHandler) DeleteUserSavedPost(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in UserHandler.DeleteUserSavedPost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	postIDParam := c.Param("postId")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in UserHandler.DeleteUserSavedPost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	if err := h.userService.DeleteUserSavedPost(userID, postID); err != nil {
		log.Printf("[Err] Error deleting user saved post in UserHandler.DeleteUserSavedPost: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete saved post",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Saved post deleted successfully",
	})
}
