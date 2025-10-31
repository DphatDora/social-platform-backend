package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/internal/service"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/util"
	"strconv"
	"time"

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

// handles SSE streaming for conversations (GET)
func (h *MessageHandler) StreamConversations(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in MessageHandler.StreamConversations", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Register SSE client
	client := h.messageService.RegisterSSEClient(userID)
	defer h.messageService.UnregisterSSEClient(userID, client)

	log.Printf("[Info] SSE stream started for user %d", userID)

	// Send initial ping
	c.SSEvent("ping", "connected")
	c.Writer.Flush()

	// Stream events
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-client.Channel:
			if !ok {
				log.Printf("[Info] SSE channel closed for user %d", userID)
				return
			}

			eventData, err := json.Marshal(event.Data)
			if err != nil {
				log.Printf("[Err] Error marshaling event data: %v", err)
				continue
			}

			// Send SSE event
			c.SSEvent(event.Event, string(eventData))
			c.Writer.Flush()

		case <-ticker.C:
			// Send keepalive ping
			c.SSEvent("ping", "keepalive")
			c.Writer.Flush()

		case <-c.Request.Context().Done():
			log.Printf("[Info] SSE client disconnected for user %d", userID)
			return
		}
	}
}

// handles SSE streaming for messages in a conversation (GET)
func (h *MessageHandler) StreamMessages(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in MessageHandler.StreamMessages", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	conversationIDParam := c.Param("conversationId")
	conversationID, err := strconv.ParseUint(conversationIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid conversation ID in MessageHandler.StreamMessages: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid conversation ID",
		})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Register SSE client
	client := h.messageService.RegisterSSEClient(userID)
	defer h.messageService.UnregisterSSEClient(userID, client)

	log.Printf("[Info] SSE message stream started for user %d in conversation %d", userID, conversationID)

	// Send initial ping
	c.SSEvent("ping", fmt.Sprintf("connected to conversation %d", conversationID))
	c.Writer.Flush()

	// Stream events
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-client.Channel:
			if !ok {
				log.Printf("[Info] SSE channel closed for user %d", userID)
				return
			}

			// Only send events related to this conversation
			if event.Event == "new_message" {
				if msgEvent, ok := event.Data.(response.NewMessageEvent); ok {
					if msgEvent.ConversationID != conversationID {
						continue // Skip messages from other conversations
					}
				}
			}

			eventData, err := json.Marshal(event.Data)
			if err != nil {
				log.Printf("[Err] Error marshaling event data: %v", err)
				continue
			}

			// Send SSE event
			c.SSEvent(event.Event, string(eventData))
			c.Writer.Flush()

		case <-ticker.C:
			// Send keepalive ping
			c.SSEvent("ping", "keepalive")
			c.Writer.Flush()

		case <-c.Request.Context().Done():
			log.Printf("[Info] SSE client disconnected for user %d", userID)
			return
		}
	}
}
