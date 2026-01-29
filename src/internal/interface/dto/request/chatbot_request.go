package request

type ChatMessage struct {
	Role    string `json:"role" binding:"required,oneof=user assistant system"`
	Content string `json:"content" binding:"required"`
}

type ChatRequest struct {
	Message             string        `json:"message" binding:"required,max=2000"`
	ConversationHistory []ChatMessage `json:"conversationHistory,omitempty"`
	Temperature         *float64      `json:"temperature,omitempty" binding:"omitempty,min=0,max=2"`
	MaxTokens           *int          `json:"maxTokens,omitempty" binding:"omitempty,min=1,max=4096"`
}

// request payload for Ollama
type OllamaRequest struct {
	Model     string              `json:"model"`
	Messasges []map[string]string `json:"messages"`
	Stream    bool                `json:"stream"`
	Options   *OllamaOptions      `json:"options,omitempty"`
}

type OllamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
	NumCtx      int     `json:"num_ctx,omitempty"`
}
