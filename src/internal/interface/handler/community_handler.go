package handler

import (
	"net/http"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/internal/service"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/logger"
	"social-platform-backend/package/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CommunityHandler struct {
	communityService *service.CommunityService
}

func NewCommunityHandler(communityService *service.CommunityService) *CommunityHandler {
	return &CommunityHandler{
		communityService: communityService,
	}
}

func (h *CommunityHandler) CreateCommunity(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.CreateCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in CommunityHandler.CreateCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.communityService.CreateCommunity(ctx, userID, &req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error creating community in CommunityHandler.CreateCommunity: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create community",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Community created successfully")
	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Community created successfully",
	})
}

func (h *CommunityHandler) GetCommunityByID(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.GetCommunityByID: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	// Get userID from context (if exists)
	userID := util.GetOptionalUserIDFromContext(c)

	community, err := h.communityService.GetCommunityByID(ctx, id, userID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting community in CommunityHandler.GetCommunityByID: %v", err)
		c.JSON(http.StatusNotFound, response.APIResponse{
			Success: false,
			Message: "Community not found",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Community retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Community retrieved successfully",
		Data:    community,
	})
}

func (h *CommunityHandler) UpdateCommunity(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.UpdateCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.UpdateCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	var req request.UpdateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in CommunityHandler.UpdateCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.communityService.UpdateCommunity(ctx, userID, id, &req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error updating community in CommunityHandler.UpdateCommunity: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "Only super admin can update community",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update community",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Community updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Community updated successfully",
	})
}

func (h *CommunityHandler) DeleteCommunity(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.DeleteCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.DeleteCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	if err := h.communityService.DeleteCommunity(ctx, userID, id); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error deleting community in CommunityHandler.DeleteCommunity: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "Only super admin can delete community",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete community",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Community deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Community deleted successfully",
	})
}

func (h *CommunityHandler) JoinCommunity(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.JoinCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.JoinCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	if err := h.communityService.JoinCommunity(ctx, userID, communityID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error joining community in CommunityHandler.JoinCommunity: %v", err)

		if strings.Contains(err.Error(), "already joined") {
			c.JSON(http.StatusConflict, response.APIResponse{
				Success: false,
				Message: "Already joined this community",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to join community",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Successfully joined community")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully joined community",
	})
}

func (h *CommunityHandler) UnjoinCommunity(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.UnjoinCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.UnjoinCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	if err := h.communityService.UnjoinCommunity(ctx, userID, communityID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error leaving community in CommunityHandler.UnjoinCommunity: %v", err)

		if strings.Contains(err.Error(), "not subscribed") {
			c.JSON(http.StatusConflict, response.APIResponse{
				Success: false,
				Message: "Not subscribed to this community",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to leave community",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Successfully left community")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully left community",
	})
}

func (h *CommunityHandler) GetCommunities(c *gin.Context) {
	ctx := c.Request.Context()
	// Get userID from context (if exists)
	userID := util.GetOptionalUserIDFromContext(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	communities, pagination, err := h.communityService.GetCommunities(ctx, page, limit, userID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting communities in CommunityHandler.GetCommunities: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get communities",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Communities retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Communities retrieved successfully",
		Data:       communities,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) SearchCommunities(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Search name is required",
		})
		return
	}

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEWEST)
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	// Get userID from context (set by OptionalAuthMiddleware)
	userID := util.GetOptionalUserIDFromContext(c)

	communities, pagination, err := h.communityService.SearchCommunitiesByName(ctx, name, sortBy, page, limit, userID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error searching communities in CommunityHandler.SearchCommunities: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to search communities",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Communities searched successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Communities searched successfully",
		Data:       communities,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) FilterCommunities(c *gin.Context) {
	ctx := c.Request.Context()
	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEWEST)

	var isPrivate *bool
	if isPrivateStr := c.Query("isPrivate"); isPrivateStr != "" {
		isPrivateVal := isPrivateStr == "true"
		isPrivate = &isPrivateVal
	}

	// Parse topics from query params (comma-separated)
	var topics []string
	if topicsStr := c.Query("topics"); topicsStr != "" {
		topics = strings.Split(topicsStr, ",")
		// Trim spaces from each topic
		for i := range topics {
			topics[i] = strings.TrimSpace(topics[i])
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	// Get userID from context (set by OptionalAuthMiddleware)
	userID := util.GetOptionalUserIDFromContext(c)

	communities, pagination, err := h.communityService.FilterCommunities(ctx, sortBy, isPrivate, topics, page, limit, userID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error filtering communities in CommunityHandler.FilterCommunities: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to filter communities",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Communities filtered successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Communities filtered successfully",
		Data:       communities,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) GetCommunityMembers(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.GetCommunityMembers", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.GetCommunityMembers: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEWEST)
	searchName := c.Query("search")
	status := c.DefaultQuery("status", constant.SUBSCRIPTION_STATUS_APPROVED)
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	members, pagination, err := h.communityService.GetCommunityMembers(ctx, userID, communityID, sortBy, searchName, status, page, limit)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting community members in CommunityHandler.GetCommunityMembers: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to view members",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get community members",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Community members retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Community members retrieved successfully",
		Data:       members,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) UpdateMemberRole(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.UpdateMemberRole", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.UpdateMemberRole: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	userIDParam := c.Param("userId")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid user ID in CommunityHandler.UpdateMemberRole: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var req request.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid request body in CommunityHandler.UpdateMemberRole: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if err := h.communityService.UpdateMemberRole(ctx, userID, communityID, targetUserID, req.Role); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error updating member role in CommunityHandler.UpdateMemberRole: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to update member roles",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community or user not found",
			})
			return
		}

		if strings.Contains(err.Error(), "not a member") {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "User is not a member of this community",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update member role",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Member role updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Member role updated successfully",
	})
}

func (h *CommunityHandler) RemoveMember(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.RemoveMember", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.RemoveMember: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	memberIDParam := c.Param("memberId")
	memberID, err := strconv.ParseUint(memberIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid member ID in CommunityHandler.RemoveMember: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid member ID",
		})
		return
	}

	if err := h.communityService.RemoveMember(ctx, userID, communityID, memberID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error removing member in CommunityHandler.RemoveMember: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to remove members",
			})
			return
		}

		if strings.Contains(err.Error(), "community not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		if strings.Contains(err.Error(), "member not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Member not found in this community",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to remove member",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Member removed successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Member removed successfully",
	})
}

func (h *CommunityHandler) GetUserRoleInCommunity(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.GetUserRoleInCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.GetUserRoleInCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	role, err := h.communityService.GetUserRoleInCommunity(ctx, userID, communityID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting user role in CommunityHandler.GetUserRoleInCommunity: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get user role",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] User role retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "User role retrieved successfully",
		Data:    gin.H{"role": role},
	})
}

func (h *CommunityHandler) VerifyCommunityName(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.VerifyCommunityNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in CommunityHandler.VerifyCommunityName: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	isUnique, err := h.communityService.VerifyCommunityName(ctx, req.Name)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error verifying community name in CommunityHandler.VerifyCommunityName: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to verify community name",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Community name verified successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Community name verification completed",
		Data:    response.VerifyCommunityNameResponse{IsUnique: isUnique},
	})
}

func (h *CommunityHandler) GetCommunityPostsForModerator(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.GetCommunityPostsForModerator", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.GetCommunityPostsForModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	status := c.Query("status")
	searchTitle := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	posts, pagination, err := h.communityService.GetCommunityPostsForModerator(ctx, userID, communityID, status, searchTitle, page, limit)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting community posts for moderator in CommunityHandler.GetCommunityPostsForModerator: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to view posts",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to retrieve posts",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Posts retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Posts retrieved successfully",
		Data:       posts,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) UpdatePostStatusByModerator(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.UpdatePostStatusByModerator", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.UpdatePostStatusByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	postIDParam := c.Param("postId")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in CommunityHandler.UpdatePostStatusByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.UpdatePostStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in CommunityHandler.UpdatePostStatusByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.communityService.UpdatePostStatusByModerator(ctx, userID, communityID, postID, req.Status); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error updating post status in CommunityHandler.UpdatePostStatusByModerator: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to update post status",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update post status",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Post status updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post status updated successfully",
	})
}

func (h *CommunityHandler) DeletePostByModerator(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.DeletePostByModerator", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.DeletePostByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	postIDParam := c.Param("postId")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in CommunityHandler.DeletePostByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	if err := h.communityService.DeletePostByModerator(ctx, userID, communityID, postID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error deleting post in CommunityHandler.DeletePostByModerator: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to delete posts",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete post",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Post deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post deleted successfully",
	})
}

func (h *CommunityHandler) DeleteCommentByModerator(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.DeleteCommentByModerator", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.DeleteCommentByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	commentIDParam := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid comment ID in CommunityHandler.DeleteCommentByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid comment ID",
		})
		return
	}

	if err := h.communityService.DeleteCommentByModerator(ctx, userID, communityID, commentID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error deleting comment in CommunityHandler.DeleteCommentByModerator: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to delete comments",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Comment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete comment",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Comment deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Comment deleted successfully",
	})
}

func (h *CommunityHandler) GetCommunityPostReports(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.GetCommunityPostReports", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.GetCommunityPostReports: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	reports, pagination, err := h.communityService.GetCommunityPostReports(ctx, userID, communityID, page, limit)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting community post reports in CommunityHandler.GetCommunityPostReports: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to view post reports",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get post reports",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Post reports retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Post reports retrieved successfully",
		Data:       reports,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) DeletePostReport(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.DeletePostReport", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.DeletePostReport: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	reportIDParam := c.Param("reportId")
	reportID, err := strconv.ParseUint(reportIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid report ID in CommunityHandler.DeletePostReport: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid report ID",
		})
		return
	}

	if err := h.communityService.DeletePostReport(ctx, userID, communityID, reportID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error deleting post report in CommunityHandler.DeletePostReport: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to delete post reports",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete post report",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Post report deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post report deleted successfully",
	})
}

func (h *CommunityHandler) GetCommunityCommentReports(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.GetCommunityCommentReports", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.GetCommunityCommentReports: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	reports, pagination, err := h.communityService.GetCommunityCommentReports(ctx, userID, communityID, page, limit)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting community comment reports in CommunityHandler.GetCommunityCommentReports: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to view comment reports",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get comment reports",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Comment reports retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Comment reports retrieved successfully",
		Data:       reports,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) DeleteCommentReport(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.DeleteCommentReport", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.DeleteCommentReport: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	reportIDParam := c.Param("reportId")
	reportID, err := strconv.ParseUint(reportIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid report ID in CommunityHandler.DeleteCommentReport: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid report ID",
		})
		return
	}

	if err := h.communityService.DeleteCommentReport(ctx, userID, communityID, reportID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error deleting comment report in CommunityHandler.DeleteCommentReport: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to delete comment reports",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to delete comment report",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Comment report deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Comment report deleted successfully",
	})
}

func (h *CommunityHandler) BanUser(c *gin.Context) {
	ctx := c.Request.Context()
	moderatorID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.BanUser", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.BanUser: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	var req request.BanUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid request body in CommunityHandler.BanUser: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if err := h.communityService.BanUser(ctx, moderatorID, communityID, &req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error banning user in CommunityHandler.BanUser: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to ban users",
			})
			return
		}

		if strings.Contains(err.Error(), "community not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] User banned successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "User banned successfully",
	})
}

func (h *CommunityHandler) GetUserRestrictionHistory(c *gin.Context) {
	ctx := c.Request.Context()
	moderatorID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.GetUserRestrictionHistory", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.GetUserRestrictionHistory: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	userIDParam := c.Param("userId")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid user ID in CommunityHandler.GetUserRestrictionHistory: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	restrictions, pagination, err := h.communityService.GetUserRestrictionHistory(ctx, moderatorID, communityID, targetUserID, page, limit)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting user restriction history in CommunityHandler.GetUserRestrictionHistory: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to view restriction history",
			})
			return
		}

		if strings.Contains(err.Error(), "community not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get restriction history",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] User restriction history retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "User restriction history retrieved successfully",
		Data:       restrictions,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) RemoveUserRestriction(c *gin.Context) {
	ctx := c.Request.Context()
	moderatorID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.RemoveUserRestriction", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.RemoveUserRestriction: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	restrictionIDParam := c.Param("restrictionId")
	restrictionID, err := strconv.ParseUint(restrictionIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid restriction ID in CommunityHandler.RemoveUserRestriction: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid restriction ID",
		})
		return
	}

	if err := h.communityService.RemoveUserRestriction(ctx, moderatorID, communityID, restrictionID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error removing user restriction in CommunityHandler.RemoveUserRestriction: %v", err)

		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to remove restrictions",
			})
			return
		}

		if strings.Contains(err.Error(), "community not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to remove restriction",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] User restriction removed successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "User restriction removed successfully",
	})
}

func (h *CommunityHandler) GetAllTopics(c *gin.Context) {
	ctx := c.Request.Context()
	searchQuery := c.Query("search")
	var search *string
	if searchQuery != "" {
		search = &searchQuery
	}

	topics, err := h.communityService.GetAllTopics(ctx, search)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting topics in CommunityHandler.GetAllTopics: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get topics",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Topics retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Topics retrieved successfully",
		Data:    topics,
	})
}

func (h *CommunityHandler) UpdateRequiresPostApproval(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.UpdateRequiresPostApproval", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.UpdateRequiresPostApproval: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	var req request.UpdateRequiresPostApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in CommunityHandler.UpdateRequiresPostApproval: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.communityService.UpdateRequiresPostApproval(ctx, userID, id, &req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error updating requires post approval in CommunityHandler.UpdateRequiresPostApproval: %v", err)

		if err.Error() == "community not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "Permission denied",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update requires post approval",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Post approval requirement updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Requires post approval updated successfully",
	})
}

func (h *CommunityHandler) UpdateRequiresMemberApproval(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.UpdateRequiresMemberApproval", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.UpdateRequiresMemberApproval: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	var req request.UpdateRequiresMemberApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in CommunityHandler.UpdateRequiresMemberApproval: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.communityService.UpdateRequiresMemberApproval(ctx, userID, id, &req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error updating requires member approval in CommunityHandler.UpdateRequiresMemberApproval: %v", err)

		if err.Error() == "community not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "Permission denied",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update requires member approval",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Member approval requirement updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Requires member approval updated successfully",
	})
}

func (h *CommunityHandler) UpdateSubscriptionStatus(c *gin.Context) {
	ctx := c.Request.Context()
	moderatorUserID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in CommunityHandler.UpdateSubscriptionStatus", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in CommunityHandler.UpdateSubscriptionStatus: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	userIDParam := c.Param("userId")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid user ID in CommunityHandler.UpdateSubscriptionStatus: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var req request.UpdateSubscriptionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid request body in CommunityHandler.UpdateSubscriptionStatus: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if err := h.communityService.UpdateSubscriptionStatus(ctx, moderatorUserID, communityID, targetUserID, req.Status); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error updating subscription status in CommunityHandler.UpdateSubscriptionStatus: %v", err)

		if err.Error() == "community not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "Permission denied",
			})
			return
		}

		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Subscription not found",
			})
			return
		}

		if err.Error() == "invalid status" {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Invalid status",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update subscription status",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Subscription status updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Subscription status updated successfully",
	})
}
