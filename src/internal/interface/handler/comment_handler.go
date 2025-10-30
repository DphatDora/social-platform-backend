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

func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	postIDParam := c.Param("postId")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid post ID in CommentHandler.GetCommentsByPostID: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(constant.DEFAULT_PAGE)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constant.DEFAULT_LIMIT)))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 12
	}

	comments, pagination, err := h.commentService.GetCommentsByPostID(postID, page, limit)
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
