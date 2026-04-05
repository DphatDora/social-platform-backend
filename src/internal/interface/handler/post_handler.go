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

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in PostHandler.CreatePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in PostHandler.CreatePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.postService.CreatePost(ctx, userID, &req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error creating post in PostHandler.CreatePost: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Post created successfully")
	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Post created successfully",
	})
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in PostHandler.UpdatePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in PostHandler.UpdatePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	postType := c.Query("type")
	if postType == "" {
		logger.ErrorfWithCtx(ctx, "[Err] Post type is required in PostHandler.UpdatePost")
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
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post type in PostHandler.UpdatePost: %s", postType)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post type",
		})
		return
	}

	if bindErr != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in PostHandler.UpdatePost: %v", bindErr)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + bindErr.Error(),
		})
		return
	}

	if err := h.postService.UpdatePost(ctx, userID, postID, postType, reqBody); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error updating post in PostHandler.UpdatePost: %v", err)

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

	logger.InfofWithCtx(ctx, "[Info] Post updated successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post updated successfully",
	})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in PostHandler.DeletePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in PostHandler.DeletePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	if err := h.postService.DeletePost(ctx, userID, postID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error deleting post in PostHandler.DeletePost: %v", err)

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

	logger.InfofWithCtx(ctx, "[Info] Post deleted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post deleted successfully",
	})
}

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	ctx := c.Request.Context()
	// Get userID from context (if exists)
	userID := util.GetOptionalUserIDFromContext(c)

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEW)

	// Parse tags from query params
	var tags []string
	if tagsStr := c.Query("tags"); tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		// Trim spaces from each tag
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
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

	posts, pagination, err := h.postService.GetAllPosts(ctx, sortBy, page, limit, tags, userID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting all posts in PostHandler.GetAllPosts: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get posts",
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

func (h *PostHandler) GetPostDetail(c *gin.Context) {
	ctx := c.Request.Context()
	// Get userID from context (if exists)
	userID := util.GetOptionalUserIDFromContext(c)

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in PostHandler.GetPostDetail: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	post, err := h.postService.GetPostDetailByID(ctx, postID, userID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting post detail in PostHandler.GetPostDetail: %v", err)

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

	logger.InfofWithCtx(ctx, "[Info] Post detail retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post detail retrieved successfully",
		Data:    post,
	})
}

func (h *PostHandler) GetPostsByCommunity(c *gin.Context) {
	ctx := c.Request.Context()
	// Get userID from context (if exists)
	userID := util.GetOptionalUserIDFromContext(c)

	idParam := c.Param("id")
	communityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid community ID in PostHandler.GetPostsByCommunity: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid community ID",
		})
		return
	}

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEW)

	// Parse tags from query params
	var tags []string
	if tagsStr := c.Query("tags"); tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		// Trim spaces from each tag
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
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

	posts, pagination, err := h.postService.GetPostsByCommunityID(ctx, communityID, sortBy, page, limit, tags, userID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting posts by community in PostHandler.GetPostsByCommunity: %v", err)

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

	logger.InfofWithCtx(ctx, "[Info] Community posts retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Posts retrieved successfully",
		Data:       posts,
		Pagination: pagination,
	})
}

func (h *PostHandler) SearchPosts(c *gin.Context) {
	ctx := c.Request.Context()
	searchQuery := c.Query("search")
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Search query is required",
		})
		return
	}

	sortBy := c.DefaultQuery("sortBy", constant.SORT_NEW)

	// Parse tags from query params
	var tags []string
	if tagsStr := c.Query("tags"); tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		// Trim spaces from each tag
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
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

	posts, pagination, err := h.postService.SearchPostsByTitle(ctx, searchQuery, sortBy, page, limit, tags, userID)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error searching posts in PostHandler.SearchPosts: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to search posts",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Posts searched successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Posts searched successfully",
		Data:       posts,
		Pagination: pagination,
	})
}

func (h *PostHandler) VotePost(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in PostHandler.VotePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in PostHandler.VotePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.VotePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in PostHandler.VotePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.postService.VotePost(ctx, userID, postID, req.Vote); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error voting post in PostHandler.VotePost: %v", err)

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

	logger.InfofWithCtx(ctx, "[Info] Post voted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post voted successfully",
	})
}

func (h *PostHandler) UnvotePost(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in PostHandler.UnvotePost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in PostHandler.UnvotePost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	if err := h.postService.UnvotePost(ctx, userID, postID); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error unvoting post in PostHandler.UnvotePost: %v", err)

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

	logger.InfofWithCtx(ctx, "[Info] Post unvoted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Post unvoted successfully",
	})
}

func (h *PostHandler) VotePoll(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in PostHandler.VotePoll", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in PostHandler.VotePoll: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.VotePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid request in PostHandler.VotePoll: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	if err := h.postService.VotePoll(ctx, userID, postID, &req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error voting poll in PostHandler.VotePoll: %v", err)

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

	logger.InfofWithCtx(ctx, "[Info] Poll voted successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Poll voted successfully",
	})
}

func (h *PostHandler) UnvotePoll(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in PostHandler.UnvotePoll", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in PostHandler.UnvotePoll: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.UnvotePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid request in PostHandler.UnvotePoll: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	if err := h.postService.UnvotePoll(ctx, userID, postID, &req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error unvoting poll in PostHandler.UnvotePoll: %v", err)

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

	logger.InfofWithCtx(ctx, "[Info] Poll vote removed successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Poll vote removed successfully",
	})
}

func (h *PostHandler) GetPostsByUser(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid user ID in PostHandler.GetPostsByUser: %v", err)
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

	posts, pagination, err := h.postService.GetPostsByUserID(ctx, userID, sortBy, page, limit)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "User not found",
			})
			return
		}

		logger.ErrorfWithCtx(ctx, "[Err] Error getting posts by user in PostHandler.GetPostsByUser: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to retrieve posts",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] User posts retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Posts retrieved successfully",
		Data:       posts,
		Pagination: pagination,
	})
}

func (h *PostHandler) ReportPost(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] %s in PostHandler.ReportPost", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Invalid post ID in PostHandler.ReportPost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req request.ReportPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error binding JSON in PostHandler.ReportPost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.postService.ReportPost(ctx, userID, postID, &req); err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error reporting post in PostHandler.ReportPost: %v", err)

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

	logger.InfofWithCtx(ctx, "[Info] Post reported successfully")
	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Post reported successfully",
	})
}

func (h *PostHandler) GetAllTags(c *gin.Context) {
	ctx := c.Request.Context()
	searchQuery := c.Query("search")
	var search *string
	if searchQuery != "" {
		search = &searchQuery
	}

	tags, err := h.postService.GetAllTags(ctx, search)
	if err != nil {
		logger.ErrorfWithCtx(ctx, "[Err] Error getting tags in PostHandler.GetAllTags: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get tags",
		})
		return
	}

	logger.InfofWithCtx(ctx, "[Info] Tags retrieved successfully")
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Tags retrieved successfully",
		Data:    tags,
	})
}
