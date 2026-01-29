package response

import "time"

// represents a chunk of streamed response from the chatbot
type ChatStreamChunk struct {
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	Timestamp time.Time `json:"timestamp"`
}

// SSE event for chatbot streaming
type ChatStreamEvent struct {
	Chunk ChatStreamChunk `json:"chunk"`
}

// SSE event for chatbot errors
type ChatErrorEvent struct {
	Error     string    `json:"error"`
	Code      string    `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}

// SSE event for chatbot completion
type ChatCompleteEvent struct {
	TotalTokens int       `json:"totalTokens,omitempty"`
	Model       string    `json:"model"`
	Timestamp   time.Time `json:"timestamp"`
}

// response payload for Ollama
type OllamaGenerateResponse struct {
	Model         string `json:"model"`
	Response      string `json:"response"`
	Done          bool   `json:"done"`
	Context       []int  `json:"context,omitempty"`
	TotalDuration int64  `json:"totalDuration,omitempty"`
}

type OllamaResponse struct {
	Model   string `json:"model"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}
