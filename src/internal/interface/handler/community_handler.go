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

type CommunityHandler struct {
	communityService *service.CommunityService
}

func NewCommunityHandler(communityService *service.CommunityService) *CommunityHandler {
	return &CommunityHandler{
		communityService: communityService,
	}
}

func (h *CommunityHandler) CreateCommunity(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.CreateCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in CommunityHandler.CreateCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.communityService.CreateCommunity(userID, &req); err != nil {
		log.Printf("[Err] Error creating community in CommunityHandler.CreateCommunity: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create community",
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Community created successfully",
	})
}

func (h *CommunityHandler) GetCommunityByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.GetCommunityByID: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	community, err := h.communityService.GetCommunityByID(id)
	if err != nil {
		log.Printf("[Err] Error getting community in CommunityHandler.GetCommunityByID: %v", err)
		c.JSON(http.StatusNotFound, response.APIResponse{
			Success: false,
			Message: "Community not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Community retrieved successfully",
		Data:    community,
	})
}

func (h *CommunityHandler) UpdateCommunity(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.UpdateCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.UpdateCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	var req request.UpdateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in CommunityHandler.UpdateCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.communityService.UpdateCommunity(userID, id, &req); err != nil {
		log.Printf("[Err] Error updating community in CommunityHandler.UpdateCommunity: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Community updated successfully",
	})
}

func (h *CommunityHandler) DeleteCommunity(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.DeleteCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.DeleteCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	if err := h.communityService.DeleteCommunity(userID, id); err != nil {
		log.Printf("[Err] Error deleting community in CommunityHandler.DeleteCommunity: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Community deleted successfully",
	})
}

func (h *CommunityHandler) JoinCommunity(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.JoinCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.JoinCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	if err := h.communityService.JoinCommunity(userID, communityID); err != nil {
		log.Printf("[Err] Error joining community in CommunityHandler.JoinCommunity: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully joined community",
	})
}

func (h *CommunityHandler) GetCommunities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	communities, pagination, err := h.communityService.GetCommunities(page, limit)
	if err != nil {
		log.Printf("[Err] Error getting communities in CommunityHandler.GetCommunities: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get communities",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Communities retrieved successfully",
		Data:       communities,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) SearchCommunities(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Search name is required",
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

	communities, pagination, err := h.communityService.SearchCommunitiesByName(name, page, limit)
	if err != nil {
		log.Printf("[Err] Error searching communities in CommunityHandler.SearchCommunities: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to search communities",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Communities searched successfully",
		Data:       communities,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) FilterCommunities(c *gin.Context) {
	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEWEST)

	var isPrivate *bool
	if isPrivateStr := c.Query("isPrivate"); isPrivateStr != "" {
		isPrivateVal := isPrivateStr == "true"
		isPrivate = &isPrivateVal
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	communities, pagination, err := h.communityService.FilterCommunities(sortBy, isPrivate, page, limit)
	if err != nil {
		log.Printf("[Err] Error filtering communities in CommunityHandler.FilterCommunities: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to filter communities",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Communities filtered successfully",
		Data:       communities,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) GetCommunityMembers(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.GetCommunityMembers", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.GetCommunityMembers: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEWEST)
	searchName := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	members, pagination, err := h.communityService.GetCommunityMembers(userID, communityID, sortBy, searchName, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting community members in CommunityHandler.GetCommunityMembers: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Community members retrieved successfully",
		Data:       members,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) RemoveMember(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.RemoveMember", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.RemoveMember: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	memberIDParam := c.Param("memberId")
	memberID, err := strconv.ParseUint(memberIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid member ID in CommunityHandler.RemoveMember: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid member ID",
		})
		return
	}

	if err := h.communityService.RemoveMember(userID, communityID, memberID); err != nil {
		log.Printf("[Err] Error removing member in CommunityHandler.RemoveMember: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Member removed successfully",
	})
}

func (h *CommunityHandler) GetUserRoleInCommunity(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.GetUserRoleInCommunity", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.GetUserRoleInCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	role, err := h.communityService.GetUserRoleInCommunity(userID, communityID)
	if err != nil {
		log.Printf("[Err] Error getting user role in CommunityHandler.GetUserRoleInCommunity: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get user role",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "User role retrieved successfully",
		Data:    gin.H{"role": role},
	})
}

func (h *CommunityHandler) VerifyCommunityName(c *gin.Context) {
	var req request.VerifyCommunityNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in CommunityHandler.VerifyCommunityName: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	isUnique, err := h.communityService.VerifyCommunityName(req.Name)
	if err != nil {
		log.Printf("[Err] Error verifying community name in CommunityHandler.VerifyCommunityName: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to verify community name",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Community name verification completed",
		Data:    response.VerifyCommunityNameResponse{IsUnique: isUnique},
	})
}

func (h *CommunityHandler) GetCommunityPostsForModerator(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.GetCommunityPostsForModerator", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.GetCommunityPostsForModerator: %v", err)
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

	posts, pagination, err := h.communityService.GetCommunityPostsForModerator(userID, communityID, status, searchTitle, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting community posts for moderator in CommunityHandler.GetCommunityPostsForModerator: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Posts retrieved successfully",
		Data:       posts,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) UpdatePostStatusByModerator(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.UpdatePostStatusByModerator", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.UpdatePostStatusByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	postIDParam := c.Param("postId")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in CommunityHandler.UpdatePostStatusByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.UpdatePostStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in CommunityHandler.UpdatePostStatusByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.communityService.UpdatePostStatusByModerator(userID, communityID, postID, req.Status); err != nil {
		log.Printf("[Err] Error updating post status in CommunityHandler.UpdatePostStatusByModerator: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post status updated successfully",
	})
}

func (h *CommunityHandler) DeletePostByModerator(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.DeletePostByModerator", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.DeletePostByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	postIDParam := c.Param("postId")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in CommunityHandler.DeletePostByModerator: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	if err := h.communityService.DeletePostByModerator(userID, communityID, postID); err != nil {
		log.Printf("[Err] Error deleting post in CommunityHandler.DeletePostByModerator: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post deleted successfully",
	})
}

func (h *CommunityHandler) GetCommunityPostReports(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.GetCommunityPostReports", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.GetCommunityPostReports: %v", err)
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

	reports, pagination, err := h.communityService.GetCommunityPostReports(userID, communityID, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting community post reports in CommunityHandler.GetCommunityPostReports: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Post reports retrieved successfully",
		Data:       reports,
		Pagination: pagination,
	})
}

func (h *CommunityHandler) DeletePostReport(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommunityHandler.DeletePostReport", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	communityIDParam := c.Param("id")
	communityID, err := strconv.ParseUint(communityIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in CommunityHandler.DeletePostReport: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	reportIDParam := c.Param("reportId")
	reportID, err := strconv.ParseUint(reportIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid report ID in CommunityHandler.DeletePostReport: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid report ID",
		})
		return
	}

	if err := h.communityService.DeletePostReport(userID, communityID, reportID); err != nil {
		log.Printf("[Err] Error deleting post report in CommunityHandler.DeletePostReport: %v", err)

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

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post report deleted successfully",
	})
}
