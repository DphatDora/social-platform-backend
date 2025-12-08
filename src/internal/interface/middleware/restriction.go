package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/util"

	"github.com/gin-gonic/gin"
)

func CheckUserRestrictionForPostMiddleware(restrictionRepo repository.UserRestrictionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetUserIDFromContext(c)
		if err != nil {
			c.Next()
			return
		}

		communityID := getCommunityIDFromPostRequest(c)
		if communityID != nil {
			communityRestriction, err := restrictionRepo.GetActiveRestrictionByUserAndCommunity(userID, *communityID)
			if err == nil && communityRestriction != nil {
				if shouldBlockAction(communityRestriction.RestrictionType) {
					handleRestriction(c, communityRestriction)
					return
				}
			}
		}

		c.Next()
	}
}

func CheckUserRestrictionForCommentMiddleware(restrictionRepo repository.UserRestrictionRepository, postRepo repository.PostRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetUserIDFromContext(c)
		if err != nil {
			c.Next()
			return
		}

		postID := getPostIDFromCommentRequest(c)
		if postID != nil {
			post, err := postRepo.GetPostByID(*postID)
			if err == nil && post != nil {
				communityRestriction, err := restrictionRepo.GetActiveRestrictionByUserAndCommunity(userID, post.CommunityID)
				if err == nil && communityRestriction != nil {
					if shouldBlockAction(communityRestriction.RestrictionType) {
						handleRestriction(c, communityRestriction)
						return
					}
				}
			}
		}

		c.Next()
	}
}

// Blocking temporary and permanent bans
func shouldBlockAction(restrictionType string) bool {
	return restrictionType == constant.RESTRICTION_TEMPORARY_BAN ||
		restrictionType == constant.RESTRICTION_PERMANENT_BAN
}

func handleRestriction(c *gin.Context, restriction *model.UserRestriction) {
	message := "You are restricted from performing this action"
	var details map[string]interface{}

	switch restriction.RestrictionType {
	case constant.RESTRICTION_TEMPORARY_BAN:
		message = "You are temporarily banned from this community"
		details = map[string]interface{}{
			"type":      "temporary_ban",
			"reason":    restriction.Reason,
			"expiresAt": restriction.ExpiresAt,
		}
	case constant.RESTRICTION_PERMANENT_BAN:
		message = "You are permanently banned from this community"
		details = map[string]interface{}{
			"type":   "permanent_ban",
			"reason": restriction.Reason,
		}
	}

	c.JSON(http.StatusForbidden, response.APIResponse{
		Success: false,
		Message: message,
		Data:    details,
	})
	c.Abort()
}

func getCommunityIDFromPostRequest(c *gin.Context) *uint64 {
	type PostRequest struct {
		CommunityID uint64 `json:"communityId"`
	}

	// Read body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}

	// Restore body for next handler
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req PostRequest
	if err := json.Unmarshal(bodyBytes, &req); err == nil && req.CommunityID > 0 {
		return &req.CommunityID
	}

	return nil
}

func getPostIDFromCommentRequest(c *gin.Context) *uint64 {
	type CommentRequest struct {
		PostID uint64 `json:"postId"`
	}

	// Read body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}

	// Restore body for next handler
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req CommentRequest
	if err := json.Unmarshal(bodyBytes, &req); err == nil && req.PostID > 0 {
		return &req.PostID
	}

	return nil
}
