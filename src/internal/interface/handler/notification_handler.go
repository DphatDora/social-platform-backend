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

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in NotificationHandler.GetNotifications", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
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

	notifications, pagination, err := h.notificationService.GetUserNotifications(userID, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting notifications in NotificationHandler.GetNotifications: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get notifications",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Notifications retrieved successfully",
		Data:       notifications,
		Pagination: pagination,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in NotificationHandler.MarkAsRead", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	notificationIDParam := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid notification ID in NotificationHandler.MarkAsRead: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid notification ID",
		})
		return
	}

	if err := h.notificationService.MarkAsRead(userID, notificationID); err != nil {
		log.Printf("[Err] Error marking notification as read in NotificationHandler.MarkAsRead: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Notification marked as read",
	})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in NotificationHandler.MarkAllAsRead", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	if err := h.notificationService.MarkAllAsRead(userID); err != nil {
		log.Printf("[Err] Error marking all notifications as read in NotificationHandler.MarkAllAsRead: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to mark all notifications as read",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "All notifications marked as read",
	})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in NotificationHandler.DeleteNotification", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	notificationIDParam := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid notification ID in NotificationHandler.DeleteNotification: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid notification ID",
		})
		return
	}

	if err := h.notificationService.DeleteNotification(userID, notificationID); err != nil {
		log.Printf("[Err] Error deleting notification in NotificationHandler.DeleteNotification: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Notification deleted successfully",
	})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in NotificationHandler.GetUnreadCount", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		log.Printf("[Err] Error getting unread count in NotificationHandler.GetUnreadCount: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get unread count",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Unread count retrieved successfully",
		Data: gin.H{
			"unreadCount": count,
		},
	})
}

func (h *NotificationHandler) GetNotificationSettings(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in NotificationHandler.GetNotificationSettings", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	settings, err := h.notificationService.GetUserNotificationSettings(userID)
	if err != nil {
		log.Printf("[Err] Error getting notification settings in NotificationHandler.GetNotificationSettings: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get notification settings",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Notification settings retrieved successfully",
		Data:    settings,
	})
}

func (h *NotificationHandler) UpdateNotificationSetting(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in NotificationHandler.UpdateNotificationSetting", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.UpdateNotificationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Invalid request body in NotificationHandler.UpdateNotificationSetting: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if err := h.notificationService.UpdateNotificationSetting(userID, req.Action, req.IsPush, req.IsSendMail); err != nil {
		log.Printf("[Err] Error updating notification setting in NotificationHandler.UpdateNotificationSetting: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Notification setting updated successfully",
	})
}
