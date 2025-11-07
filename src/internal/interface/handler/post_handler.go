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

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in PostHandler.CreatePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in PostHandler.CreatePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.postService.CreatePost(userID, &req); err != nil {
		log.Printf("[Err] Error creating post in PostHandler.CreatePost: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Post created successfully",
	})
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in PostHandler.UpdatePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in PostHandler.UpdatePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	postType := c.Query("type")
	if postType == "" {
		log.Printf("[Err] Post type is required in PostHandler.UpdatePost")
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Post type is required",
		})
		return
	}

	var reqBody interface{}
	var bindErr error

	switch postType {
	case constant.PostTypeText:
		var req request.UpdatePostTextRequest
		bindErr = c.ShouldBindJSON(&req)
		reqBody = &req
	case constant.PostTypeLink:
		var req request.UpdatePostLinkRequest
		bindErr = c.ShouldBindJSON(&req)
		reqBody = &req
	case constant.PostTypeMedia:
		var req request.UpdatePostMediaRequest
		bindErr = c.ShouldBindJSON(&req)
		reqBody = &req
	case constant.PostTypePoll:
		var req request.UpdatePostPollRequest
		bindErr = c.ShouldBindJSON(&req)
		reqBody = &req
	default:
		log.Printf("[Err] Invalid post type in PostHandler.UpdatePost: %s", postType)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post type",
		})
		return
	}

	if bindErr != nil {
		log.Printf("[Err] Error binding JSON in PostHandler.UpdatePost: %v", bindErr)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + bindErr.Error(),
		})
		return
	}

	if err := h.postService.UpdatePost(userID, postID, postType, reqBody); err != nil {
		log.Printf("[Err] Error updating post in PostHandler.UpdatePost: %v", err)

		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to update this post",
			})
			return
		}

		if err.Error() == "post type mismatch" {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Post type mismatch",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to update post",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post updated successfully",
	})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in PostHandler.DeletePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in PostHandler.DeletePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	if err := h.postService.DeletePost(userID, postID); err != nil {
		log.Printf("[Err] Error deleting post in PostHandler.DeletePost: %v", err)

		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "You don't have permission to delete this post",
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

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	// Get userID from context (if exists)
	userID := util.GetOptionalUserIDFromContext(c)

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEW)
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	posts, pagination, err := h.postService.GetAllPosts(sortBy, page, limit, userID)
	if err != nil {
		log.Printf("[Err] Error getting all posts in PostHandler.GetAllPosts: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get posts",
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

func (h *PostHandler) GetPostDetail(c *gin.Context) {
	// Get userID from context (if exists)
	userID := util.GetOptionalUserIDFromContext(c)

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in PostHandler.GetPostDetail: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	post, err := h.postService.GetPostDetailByID(postID, userID)
	if err != nil {
		log.Printf("[Err] Error getting post detail in PostHandler.GetPostDetail: %v", err)

		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get post detail",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post detail retrieved successfully",
		Data:    post,
	})
}

func (h *PostHandler) GetPostsByCommunity(c *gin.Context) {
	// Get userID from context (if exists)
	userID := util.GetOptionalUserIDFromContext(c)

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid community ID in PostHandler.GetPostsByCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEW)
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	posts, pagination, err := h.postService.GetPostsByCommunityID(communityID, sortBy, page, limit, userID)
	if err != nil {
		log.Printf("[Err] Error getting posts by community in PostHandler.GetPostsByCommunity: %v", err)

		if err.Error() == "community not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Community not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get posts",
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

func (h *PostHandler) SearchPosts(c *gin.Context) {
	searchQuery := c.Query("search")
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Search query is required",
		})
		return
	}

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEW)
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

	posts, pagination, err := h.postService.SearchPostsByTitle(searchQuery, sortBy, page, limit, userID)
	if err != nil {
		log.Printf("[Err] Error searching posts in PostHandler.SearchPosts: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to search posts",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Posts searched successfully",
		Data:       posts,
		Pagination: pagination,
	})
}

func (h *PostHandler) VotePost(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in PostHandler.VotePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in PostHandler.VotePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.VotePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in PostHandler.VotePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.postService.VotePost(userID, postID, req.Vote); err != nil {
		log.Printf("[Err] Error voting post in PostHandler.VotePost: %v", err)

		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to vote post",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post voted successfully",
	})
}

func (h *PostHandler) UnvotePost(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in PostHandler.UnvotePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in PostHandler.UnvotePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	if err := h.postService.UnvotePost(userID, postID); err != nil {
		log.Printf("[Err] Error unvoting post in PostHandler.UnvotePost: %v", err)

		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to unvote post",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post unvoted successfully",
	})
}

func (h *PostHandler) VotePoll(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in PostHandler.VotePoll", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in PostHandler.VotePoll: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.VotePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Invalid request in PostHandler.VotePoll: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	if err := h.postService.VotePoll(userID, postID, &req); err != nil {
		log.Printf("[Err] Error voting poll in PostHandler.VotePoll: %v", err)

		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		if err.Error() == "post is not a poll" || err.Error() == "option not found" ||
			err.Error() == "poll has expired" || err.Error() == "already voted for this option" {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to vote poll",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Poll voted successfully",
	})
}

func (h *PostHandler) UnvotePoll(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in PostHandler.UnvotePoll", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in PostHandler.UnvotePoll: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.UnvotePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Invalid request in PostHandler.UnvotePoll: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	if err := h.postService.UnvotePoll(userID, postID, &req); err != nil {
		log.Printf("[Err] Error unvoting poll in PostHandler.UnvotePoll: %v", err)

		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		if err.Error() == "post is not a poll" || err.Error() == "option not found" ||
			err.Error() == "you have not voted for this option" {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to unvote poll",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Poll vote removed successfully",
	})
}

func (h *PostHandler) GetPostsByUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid user ID in PostHandler.GetPostsByUser: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEW)
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}
	if limit < 1 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	posts, pagination, err := h.postService.GetPostsByUserID(userID, sortBy, page, limit)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "User not found",
			})
			return
		}

		log.Printf("[Err] Error getting posts by user in PostHandler.GetPostsByUser: %v", err)
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

func (h *PostHandler) ReportPost(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in PostHandler.ReportPost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in PostHandler.ReportPost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.ReportPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in PostHandler.ReportPost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.postService.ReportPost(userID, postID, &req); err != nil {
		log.Printf("[Err] Error reporting post in PostHandler.ReportPost: %v", err)

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}

		if strings.Contains(err.Error(), "already reported") {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "You have already reported this post",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to report post",
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Post reported successfully",
	})
}
