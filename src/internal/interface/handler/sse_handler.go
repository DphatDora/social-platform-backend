package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/internal/service"
	"social-platform-backend/package/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SSEHandler struct {
	sseService *service.SSEService
}

func NewSSEHandler(sseService *service.SSEService) *SSEHandler {
	return &SSEHandler{
		sseService: sseService,
	}
}

func (h *SSEHandler) Stream(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in SSEHandler.Stream", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// Register SSE client
	client := h.sseService.RegisterClient(userID)
	defer h.sseService.UnregisterClient(userID, client)

	log.Printf("[Info] SSE stream started for user %d", userID)

	// Send initial connection event
	c.SSEvent("ping", "connected")
	c.Writer.Flush()

	// Keep connection alive and listen for events
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			log.Printf("[Info] SSE client disconnected for user %d", userID)
			return

		case event, ok := <-client.Channel:
			if !ok {
				log.Printf("[Info] SSE channel closed for user %d", userID)
				return
			}

			// Marshal event data to JSON
			eventData, err := json.Marshal(event.Data)
			if err != nil {
				log.Printf("[Err] Failed to marshal SSE event data: %v", err)
				continue
			}

			// Send SSE event (conversation_update, new_notification, ...)
			c.SSEvent(event.Event, string(eventData))
			c.Writer.Flush()

		case <-ticker.C:
			// Send keepalive ping
			c.SSEvent("ping", "keepalive")
			c.Writer.Flush()
		}
	}
}

func (h *SSEHandler) StreamConversationMessages(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] %s in SSEHandler.StreamConversationMessages", err.Error())
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	conversationIDParam := c.Param("conversationId")
	conversationID, err := strconv.ParseUint(conversationIDParam, 10, 64)
	if err != nil {
		log.Printf("[Err] Invalid conversation ID in SSEHandler.StreamConversationMessages: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid conversation ID",
		})
		return
	}

	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// Register SSE client
	client := h.sseService.RegisterClient(userID)
	defer h.sseService.UnregisterClient(userID, client)

	log.Printf("[Info] SSE conversation stream started for user %d in conversation %d", userID, conversationID)

	// Send initial ping
	c.SSEvent("ping", fmt.Sprintf("connected to conversation %d", conversationID))
	c.Writer.Flush()

	// Stream events
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			log.Printf("[Info] SSE client disconnected for user %d", userID)
			return

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
		}
	}
}
