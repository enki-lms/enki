package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/enki/daemon/internal/config"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a request to the AI
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse represents a response from the AI
type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// Client handles communication with OAI-compatible API
type Client struct {
	cfg        *config.AIConfig
	httpClient *http.Client
}

// NewClient creates a new AI client
func NewClient(cfg *config.AIConfig) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Chat sends a message to the AI and returns the response
func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	// Prepend system message if configured
	allMessages := messages
	if c.cfg.SystemPrompt != "" {
		allMessages = append([]Message{{
			Role:    "system",
			Content: c.cfg.SystemPrompt,
		}}, messages...)
	}

	reqBody := ChatRequest{
		Model:    c.cfg.Model,
		Messages: allMessages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.cfg.BaseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// IsEnabled returns whether AI is enabled
func (c *Client) IsEnabled() bool {
	return c.cfg.Enabled && c.cfg.APIKey != ""
}

// GradeImage grades a student's image against an ideal image
func (c *Client) GradeImage(ctx context.Context, studentImageBase64, idealImageUrl string) (int, string, error) {
	if !c.IsEnabled() {
		return 0, "AI grading is not enabled", fmt.Errorf("AI grading is not enabled")
	}

	prompt := `Compare the student's drawing (first image) with the ideal solution (second image).
	Rate the student's drawing on a scale of 0 to 100 based on how closely it matches the ideal solution.
	Ignore minor differences in stroke width or exact positioning if the shape and relative proportions are correct.
	Provide helpful feedback for the student.
	
	Return your response in JSON format:
	{
		"score": <0-100>,
		"feedback": "<feedback string>"
	}`

	// Construct multimodal request manually since Message struct is simple
	reqBody := map[string]interface{}{
		"model": c.cfg.Model,
		"messages": []interface{}{
			map[string]interface{}{
				"role": "user",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": prompt,
					},
					map[string]interface{}{
						"type": "image_url",
						"image_url": map[string]string{
							"url": "data:image/svg+xml;base64," + studentImageBase64,
						},
					},
					map[string]interface{}{
						"type": "image_url",
						"image_url": map[string]string{
							"url": idealImageUrl,
						},
					},
				},
			},
		},
		"response_format": map[string]string{"type": "json_object"},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return 0, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.cfg.BaseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return 0, "", fmt.Errorf("failed to parse response: %w", err)
	}

	if chatResp.Error != nil {
		return 0, "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return 0, "", fmt.Errorf("no response from AI")
	}

	content := chatResp.Choices[0].Message.Content

	var result struct {
		Score    int    `json:"score"`
		Feedback string `json:"feedback"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return 0, "", fmt.Errorf("failed to parse AI response json: %w", err)
	}

	return result.Score, result.Feedback, nil
}
