package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"social-platform-backend/config"
)

type AIServiceClient struct {
	conf       config.AIService
	httpClient *http.Client
}

type aiServiceCheckContentRequest struct {
	Content   string   `json:"content,omitempty"`
	ImageURLs []string `json:"imageUrls,omitempty"`
}

type AIServiceModerationResult struct {
	IsViolation bool   `json:"isViolation"`
	Reason      string `json:"reason,omitempty"`
	Category    string `json:"category,omitempty"`
}

type aiServiceAPIResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    AIServiceModerationResult `json:"data"`
}

func NewAIServiceClient(conf *config.Config) *AIServiceClient {
	timeout := time.Duration(conf.AIService.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &AIServiceClient{
		conf: conf.AIService,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *AIServiceClient) CheckContent(ctx context.Context, content string, imageURLs []string) (*AIServiceModerationResult, error) {
	if strings.TrimSpace(c.conf.BaseURL) == "" {
		return nil, fmt.Errorf("ai service base URL is not configured")
	}

	requestBody := aiServiceCheckContentRequest{
		Content:   content,
		ImageURLs: imageURLs,
	}

	data, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal ai service request: %w", err)
	}

	endpoint := strings.TrimRight(c.conf.BaseURL, "/") + "/api/v1/moderation/check-content"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("create ai service request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(c.conf.APIKey) != "" {
		req.Header.Set("X-API-Key", c.conf.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ai service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read ai service response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ai service error: status %d, body: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var apiResp aiServiceAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("decode ai service response: %w", err)
	}
	if !apiResp.Success {
		return nil, fmt.Errorf("ai service rejected request: %s", apiResp.Message)
	}

	return &apiResp.Data, nil
}
