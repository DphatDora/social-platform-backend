package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"social-platform-backend/config"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
	"time"
)

type ChatbotService struct {
	config     *config.Config
	httpClient *http.Client
}

func NewChatbotService(conf *config.Config) *ChatbotService {
	timeout := time.Duration(conf.Ollama.Timeout) * time.Second
	if timeout == 0 {
		timeout = constant.DEFAUTL_TIMEOUT * time.Second
	}

	return &ChatbotService{
		config: conf,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *ChatbotService) StreamChat(ctx context.Context, req *request.ChatRequest, chunkChan chan<- response.ChatStreamChunk, errorChan chan<- error) {
	defer close(chunkChan)
	defer close(errorChan)

	// Build messages array for chat API
	messages := []map[string]string{}

	// Add system message
	messages = append(messages, map[string]string{
		"role":    "system",
		"content": constant.SYSTEM_PROMPT,
	})

	// Add conversation history if provided
	if len(req.ConversationHistory) > 0 {
		for _, msg := range req.ConversationHistory {
			messages = append(messages, map[string]string{
				"role":    msg.Role,
				"content": msg.Content,
			})
		}
	}

	// Add current user message
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": req.Message,
	})

	options := &request.OllamaOptions{
		NumCtx: 4096,
	}

	if req.Temperature != nil {
		options.Temperature = *req.Temperature
	}

	if req.MaxTokens != nil {
		options.NumPredict = *req.MaxTokens
	} else {
		options.NumPredict = constant.DEFAULT_MAX_TOKENS
	}

	ollamaReq := request.OllamaRequest{
		Model:     s.config.Ollama.Model,
		Messasges: messages,
		Stream:    true,
		Options:   options,
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		errorChan <- fmt.Errorf("failed to marshal request: %w", err)
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/api/chat", s.config.Ollama.BaseURL),
		bytes.NewBuffer(reqBody))
	if err != nil {
		errorChan <- fmt.Errorf("failed to create request: %w", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		errorChan <- fmt.Errorf("failed to send request to Ollama: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errorChan <- fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
		return
	}

	// Stream response line by line
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			errorChan <- ctx.Err()
			return
		default:
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var ollamaResp response.OllamaResponse

			if err := json.Unmarshal(line, &ollamaResp); err != nil {
				log.Printf("[Warn] Failed to parse Ollama response chunk: %v", err)
				continue
			}

			// Send chunk to channel
			chunk := response.ChatStreamChunk{
				Content:   ollamaResp.Message.Content,
				Done:      ollamaResp.Done,
				Timestamp: time.Now(),
			}

			select {
			case chunkChan <- chunk:
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			}

			if ollamaResp.Done {
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		errorChan <- fmt.Errorf("error reading response stream: %w", err)
	}
}
