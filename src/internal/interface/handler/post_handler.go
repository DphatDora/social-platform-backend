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
