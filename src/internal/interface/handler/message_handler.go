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

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in MessageHandler.SendMessage", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Invalid request in MessageHandler.SendMessage: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Validate user is not sending message to themselves
	if req.RecipientID == userID {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Cannot send message to yourself",
		})
		return
	}

	message, err := h.messageService.SendMessage(userID, &req)
	if err != nil {
		log.Printf("[Err] Error sending message in MessageHandler.SendMessage: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Message sent successfully",
		Data:    message,
	})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in MessageHandler.MarkAsRead", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	messageIDParam := c.Param("messageId")
	messageID, err := strconv.ParseUint(messageIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid message ID in MessageHandler.MarkAsRead: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid message ID",
		})
		return
	}

	if err := h.messageService.MarkMessageAsRead(userID, messageID); err != nil {
		log.Printf("[Err] Error marking message as read in MessageHandler.MarkAsRead: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Message marked as read",
	})
}

func (h *MessageHandler) MarkConversationAsRead(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in MessageHandler.MarkConversationAsRead", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	conversationIDParam := c.Param("conversationId")
	conversationID, err := strconv.ParseUint(conversationIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid conversation ID in MessageHandler.MarkConversationAsRead: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid conversation ID",
		})
		return
	}

	if err := h.messageService.MarkConversationAsRead(userID, conversationID); err != nil {
		log.Printf("[Err] Error marking conversation as read in MessageHandler.MarkConversationAsRead: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Conversation marked as read",
	})
}

func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in MessageHandler.GetConversations", err.Error())
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

	conversations, pagination, err := h.messageService.GetConversations(userID, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting conversations in MessageHandler.GetConversations: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get conversations",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Conversations retrieved successfully",
		Data:       conversations,
		Pagination: pagination,
	})
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in MessageHandler.GetMessages", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	conversationIDParam := c.Param("conversationId")
	conversationID, err := strconv.ParseUint(conversationIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid conversation ID in MessageHandler.GetMessages: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid conversation ID",
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

	messages, pagination, err := h.messageService.GetMessages(userID, conversationID, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting messages in MessageHandler.GetMessages: %v", err)
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success:    true,
		Message:    "Messages retrieved successfully",
		Data:       messages,
		Pagination: pagination,
	})
}
