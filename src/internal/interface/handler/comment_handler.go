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

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommentHandler.CreateComment", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in CommentHandler.CreateComment: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.commentService.CreateComment(userID, &req); err != nil {
		log.Printf("[Err] Error creating comment in CommentHandler.CreateComment: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Comment created successfully",
	})
}

func (h *CommentHandler) GetCommentsOnPost(c *gin.Context) {
	// Get userID from context (if exists)
	userID := util.GetOptionalUserIDFromContext(c)

	postIDParam := c.Param("id")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in CommentHandler.GetCommentsOnPost: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))
	sortBy := c.DefaultQuery("sortBy", constant.COMMENT_SORT_NEWEST)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 12
	}

	comments, pagination, err := h.commentService.GetCommentsByPostID(postID, sortBy, page, limit, userID)
	if err != nil {
		log.Printf("[Err] Error getting comments in CommentHandler.GetCommentsByPostID: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get comments",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Comments retrieved successfully",
		Data:       comments,
		Pagination: pagination,
	})
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommentHandler.UpdateComment", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	commentIDParam := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid comment ID in CommentHandler.UpdateComment: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid comment ID",
		})
		return
	}

	var req request.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in CommentHandler.UpdateComment: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.commentService.UpdateComment(userID, commentID, &req); err != nil {
		log.Printf("[Err] Error updating comment in CommentHandler.UpdateComment: %v", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "comment not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Comment updated successfully",
	})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommentHandler.DeleteComment", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	commentIDParam := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid comment ID in CommentHandler.DeleteComment: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid comment ID",
		})
		return
	}

	if err := h.commentService.DeleteComment(userID, commentID); err != nil {
		log.Printf("[Err] Error deleting comment in CommentHandler.DeleteComment: %v", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "comment not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Comment deleted successfully",
	})
}

func (h *CommentHandler) VoteComment(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommentHandler.VoteComment", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	commentIDParam := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid comment ID in CommentHandler.VoteComment: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid comment ID",
		})
		return
	}

	var req request.VoteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in CommentHandler.VoteComment: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.commentService.VoteComment(userID, commentID, req.Vote); err != nil {
		log.Printf("[Err] Error voting comment in CommentHandler.VoteComment: %v", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "comment not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Comment voted successfully",
	})
}

func (h *CommentHandler) UnvoteComment(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommentHandler.UnvoteComment", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	commentIDParam := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid comment ID in CommentHandler.UnvoteComment: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid comment ID",
		})
		return
	}

	if err := h.commentService.UnvoteComment(userID, commentID); err != nil {
		log.Printf("[Err] Error unvoting comment in CommentHandler.UnvoteComment: %v", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "comment not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Comment unvoted successfully",
	})
}

func (h *CommentHandler) GetCommentsByUser(c *gin.Context) {
	// Get requestUserID from context (if exists) - this is the user viewing the comments
	requestUserID := util.GetOptionalUserIDFromContext(c)

	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid user ID in CommentHandler.GetCommentsByUser: %v", err)
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

	comments, pagination, err := h.commentService.GetCommentsByUserID(userID, sortBy, page, limit, requestUserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "User not found",
			})
			return
		}

		log.Printf("[Err] Error getting comments by user in CommentHandler.GetCommentsByUser: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to retrieve comments",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Comments retrieved successfully",
		Data:       comments,
		Pagination: pagination,
	})
}

func (h *CommentHandler) ReportComment(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in CommentHandler.ReportComment", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	commentIDParam := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid comment ID in CommentHandler.ReportComment: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid comment ID",
		})
		return
	}

	var req request.ReportCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in CommentHandler.ReportComment: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.commentService.ReportComment(userID, commentID, &req); err != nil {
		log.Printf("[Err] Error reporting comment in CommentHandler.ReportComment: %v", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "comment not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "you have already reported this comment" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Comment reported successfully",
	})
}
