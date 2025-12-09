package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"social-platform-backend/config"
)

type SensitiveKeywords struct {
	Profanity     KeywordLanguage `json:"profanity"`
	HateSpeech    KeywordLanguage `json:"hate_speech"`
	Politics      KeywordLanguage `json:"politics"`
	SexualContent KeywordLanguage `json:"sexual_content"`
	Violence      KeywordLanguage `json:"violence"`
	Drugs         KeywordLanguage `json:"drugs"`
	Scam          KeywordLanguage `json:"scam"`
}

type KeywordLanguage struct {
	Vi []string `json:"vi"`
	En []string `json:"en"`
}

type ContentViolation struct {
	IsViolation bool
	Reason      string
	Category    string
}

var keywordsCache *SensitiveKeywords

func LoadSensitiveKeywords() (*SensitiveKeywords, error) {
	if keywordsCache != nil {
		return keywordsCache, nil
	}

	filePath := filepath.Join("package", "data", "sensitive_keywords.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keywords file: %v", err)
	}

	var keywords SensitiveKeywords
	if err := json.Unmarshal(data, &keywords); err != nil {
		return nil, fmt.Errorf("failed to parse keywords: %v", err)
	}

	keywordsCache = &keywords
	return keywordsCache, nil
}

func stripHTML(htmlContent string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(htmlContent, " ")

	// Decode common HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&apos;", "'")

	// Remove extra whitespace
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text
}

func CheckTextContent(text string) (*ContentViolation, error) {
	if text == "" {
		return &ContentViolation{IsViolation: false}, nil
	}

	keywords, err := LoadSensitiveKeywords()
	if err != nil {
		log.Printf("[Err] Failed to load sensitive keywords: %v", err)
		return nil, err
	}

	plainText := stripHTML(text)
	textLower := strings.ToLower(plainText)

	// Check each category
	categories := map[string]struct {
		keywords []string
		name     string
	}{
		"profanity":      {append(keywords.Profanity.Vi, keywords.Profanity.En...), "Inappropriate or offensive content"},
		"hate_speech":    {append(keywords.HateSpeech.Vi, keywords.HateSpeech.En...), "Hate speech or discriminatory content"},
		"politics":       {append(keywords.Politics.Vi, keywords.Politics.En...), "Sensitive political content"},
		"sexual_content": {append(keywords.SexualContent.Vi, keywords.SexualContent.En...), "Sexual or pornographic content"},
		"violence":       {append(keywords.Violence.Vi, keywords.Violence.En...), "Violent or threatening content"},
		"drugs":          {append(keywords.Drugs.Vi, keywords.Drugs.En...), "Drug-related content"},
		"scam":           {append(keywords.Scam.Vi, keywords.Scam.En...), "Fraudulent or scam content"},
	}

	for category, data := range categories {
		for _, keyword := range data.keywords {
			keywordLower := strings.ToLower(keyword)
			if containsWordBoundary(textLower, keywordLower) {
				return &ContentViolation{
					IsViolation: true,
					Reason:      data.name,
					Category:    category,
				}, nil
			}
		}
	}

	return &ContentViolation{IsViolation: false}, nil
}

func containsWordBoundary(text, keyword string) bool {
	if !strings.Contains(text, keyword) {
		return false
	}

	// Check if keyword appears as a complete word
	textWithBoundaries := " " + text + " "
	keywordWithBoundaries := " " + keyword + " "

	if strings.Contains(textWithBoundaries, keywordWithBoundaries) {
		return true
	}

	// Check with common punctuation as boundaries
	punctuations := []string{".", ",", "!", "?", ":", ";", "\"", "'", "\n", "\t", "(", ")", "[", "]"}
	for _, p := range punctuations {
		if strings.Contains(text, p+keyword+p) ||
			strings.Contains(text, p+keyword+" ") ||
			strings.Contains(text, " "+keyword+p) {
			return true
		}
	}

	if text == keyword {
		return true
	}

	return false
}

func CheckImageContent(imageURL string) (*ContentViolation, error) {
	if imageURL == "" {
		return &ContentViolation{IsViolation: false}, nil
	}

	// Skip video files
	lowerURL := strings.ToLower(imageURL)
	videoExtensions := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv", ".webm", ".m4v"}
	for _, ext := range videoExtensions {
		if strings.Contains(lowerURL, ext) {
			log.Printf("[Info] Skipping video file: %s", imageURL)
			return &ContentViolation{IsViolation: false}, nil
		}
	}

	conf := config.GetConfig()
	if conf.Gemini.APIKey == "" {
		log.Printf("[Warning] Gemini API key not configured, skipping image moderation")
		return &ContentViolation{IsViolation: false}, nil
	}

	// Download and encode image to base64
	base64Image, mimeType, err := downloadImageAsBase64(imageURL)
	if err != nil {
		log.Printf("[Err] Failed to download image: %v", err)
		return &ContentViolation{IsViolation: false}, nil
	}

	// Call Gemini API with base64 image
	result, err := callGeminiVisionAPI(conf.Gemini.APIKey, base64Image, mimeType)
	if err != nil {
		log.Printf("[Err] Failed to call Gemini API: %v", err)
		return &ContentViolation{IsViolation: false}, nil
	}

	if result.IsViolation {
		return &ContentViolation{
			IsViolation: true,
			Reason:      result.Reason,
			Category:    "inappropriate_image",
		}, nil
	}

	return &ContentViolation{IsViolation: false}, nil
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text       string            `json:"text,omitempty"`
	InlineData *GeminiInlineData `json:"inlineData,omitempty"`
	FileData   *GeminiFileData   `json:"fileData,omitempty"`
}

type GeminiInlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type GeminiFileData struct {
	FileUri  string `json:"fileUri"`
	MimeType string `json:"mimeType"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func downloadImageAsBase64(imageURL string) (string, string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Detect mime type from content type header or URL
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		lowerURL := strings.ToLower(imageURL)
		if strings.Contains(lowerURL, ".png") {
			mimeType = "image/png"
		} else if strings.Contains(lowerURL, ".webp") {
			mimeType = "image/webp"
		} else if strings.Contains(lowerURL, ".gif") {
			mimeType = "image/gif"
		} else {
			mimeType = "image/jpeg"
		}
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)
	return base64Image, mimeType, nil
}

func callGeminiVisionAPI(apiKey, base64Image, mimeType string) (*ContentViolation, error) {
	prompt := `Analyze this image and determine if it contains any inappropriate content including:
- Violence, gore, or graphic content
- Sexual or pornographic content
- Hate symbols or extremist content
- Disturbing or shocking imagery

Respond in JSON format:
{
  "is_violation": true/false,
  "reason": "brief description if violation found"
}`

	requestBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
					{
						InlineData: &GeminiInlineData{
							MimeType: mimeType,
							Data:     base64Image,
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return &ContentViolation{IsViolation: false}, nil
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

	var result struct {
		IsViolation bool   `json:"is_violation"`
		Reason      string `json:"reason"`
	}

	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := responseText[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			log.Printf("[Warning] Failed to parse Gemini response as JSON: %v", err)
			lowerText := strings.ToLower(responseText)
			if strings.Contains(lowerText, "violation") || strings.Contains(lowerText, "inappropriate") ||
				strings.Contains(lowerText, "sexual") || strings.Contains(lowerText, "violence") {
				return &ContentViolation{
					IsViolation: true,
					Reason:      "Image contains inappropriate content",
				}, nil
			}
		}
	}

	if result.IsViolation {
		return &ContentViolation{
			IsViolation: true,
			Reason:      fmt.Sprintf("Image violation: %s", result.Reason),
		}, nil
	}

	return &ContentViolation{IsViolation: false}, nil
}
