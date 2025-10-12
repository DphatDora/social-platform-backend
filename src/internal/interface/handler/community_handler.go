package handler

import (
	"log"
	"net/http"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/internal/service"
	"social-platform-backend/package/constant"
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
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Printf("[Err] UserID not found in context in CommunityHandler.CreateCommunity")
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		log.Printf("[Err] Invalid userID type in context in CommunityHandler.CreateCommunity")
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to retrieve user information",
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

	if err := h.communityService.UpdateCommunity(id, &req); err != nil {
		log.Printf("[Err] Error updating community in CommunityHandler.UpdateCommunity: %v", err)
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

	if err := h.communityService.DeleteCommunity(id); err != nil {
		log.Printf("[Err] Error deleting community in CommunityHandler.DeleteCommunity: %v", err)
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
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Printf("[Err] UserID not found in context in CommunityHandler.JoinCommunity")
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		log.Printf("[Err] Invalid userID type in context in CommunityHandler.JoinCommunity")
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to retrieve user information",
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
