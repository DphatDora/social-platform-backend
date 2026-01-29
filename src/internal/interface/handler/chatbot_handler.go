package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/internal/service"
	"social-platform-backend/package/util"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatbotHandler struct {
	chatbotService *service.ChatbotService
}

func NewChatbotHandler(chatbotService *service.ChatbotService) *ChatbotHandler {
	return &ChatbotHandler{
		chatbotService: chatbotService,
	}
}

func (h *ChatbotHandler) StreamChat(c *gin.Context) {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("[Err] Invalid user ID in ChatbotHandler.StreamChat, %v", err)
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req request.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in ChatbotHandler.StreamChat, %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request payload" + err.Error(),
		})
		return
	}

	// set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	log.Printf("[Info] Chatbot stream started for user %d", userID)

	c.SSEvent("connected", "stream started")
	c.Writer.Flush()

	// create channels for streaming
	chunkChan := make(chan response.ChatStreamChunk, 10)
	errorChan := make(chan error, 1)

	go h.chatbotService.StreamChat(c.Request.Context(), &req, chunkChan, errorChan)

	for {
		select {
		case <-c.Request.Context().Done():
			log.Printf("[Info] Client disconnected for user %d", userID)
			return

		case err := <-errorChan:
			if err != nil {
				log.Printf("[Err] Chatbot error for user %d: %v", userID, err)

				errorEvent := response.SSEEvent{
					Event: "error",
					Data: response.ChatErrorEvent{
						Error:     err.Error(),
						Code:      "CHATBOT_ERROR",
						Timestamp: time.Now(),
					},
				}

				eventData, _ := json.Marshal(errorEvent.Data)
				c.SSEvent(errorEvent.Event, string(eventData))
				c.Writer.Flush()
				return
			}

		case chunk, ok := <-chunkChan:
			if !ok {
				completeEvent := response.SSEEvent{
					Event: "complete",
					Data: response.ChatCompleteEvent{
						Model:     "ollama",
						Timestamp: time.Now(),
					},
				}

				eventData, _ := json.Marshal(completeEvent.Data)
				c.SSEvent(completeEvent.Event, string(eventData))
				c.Writer.Flush()

				log.Printf("[Info] Chatbot stream completed for user %d", userID)
				return
			}

			// Send chunk event
			chunkEvent := response.SSEEvent{
				Event: "chunk",
				Data: response.ChatStreamEvent{
					Chunk: chunk,
				},
			}

			eventData, _ := json.Marshal(chunkEvent.Data)
			c.SSEvent(chunkEvent.Event, string(eventData))
			c.Writer.Flush()
		}
	}
}
