package translator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	anthropicAPIURL = "https://api.anthropic.com/v1/messages"
)

// AnthropicProvider provides translation using Anthropic Claude API
type AnthropicProvider struct{}

// GetName returns the name of the provider
func (p *AnthropicProvider) GetName() string {
	return "anthropic"
}

// CreateTranslator creates a translator instance
func (p *AnthropicProvider) CreateTranslator(config map[string]interface{}) (Translator, error) {
	apiKey, ok := config["api_key"].(string)
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required")
	}

	model, ok := config["model"].(string)
	if !ok || model == "" {
		model = "claude-3-haiku-20240307" // Default model
	}

	return &AnthropicTranslator{
		apiKey: apiKey,
		model:  model,
	}, nil
}

// AnthropicTranslator implements the Translator interface using Anthropic API
type AnthropicTranslator struct {
	apiKey string
	model  string
}

// AnthropicRequest represents a request to the Anthropic API
type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []AnthropicMessage `json:"messages"`
}

// AnthropicMessage represents a message in the Anthropic API
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnthropicResponse represents a response from the Anthropic API
type AnthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Translate translates text to the specified target language
func (t *AnthropicTranslator) Translate(ctx context.Context, text string, targetLang string) (string, error) {
	systemPrompt := "You are a professional translator. Your task is to translate text accurately while preserving all formatting, placeholders, and special characters."
	userPrompt := fmt.Sprintf("Please translate the following text to %s, preserving all formatting and placeholders:\n\n%s", targetLang, text)

	requestBody := AnthropicRequest{
		Model:     t.model,
		MaxTokens: 1024,
		System:    systemPrompt,
		Messages: []AnthropicMessage{
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", anthropicAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", t.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response AnthropicResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("Anthropic API error: %s", response.Error.Message)
	}

	// Extract text from the response
	var translation string
	for _, content := range response.Content {
		if content.Type == "text" {
			translation = content.Text
			break
		}
	}

	if translation == "" {
		return "", fmt.Errorf("empty or invalid response from Anthropic API")
	}

	// Sometimes the model might include extraneous text about the translation,
	// so we try to extract just the translated content.
	if strings.Contains(translation, "\n\n") {
		parts := strings.SplitN(translation, "\n\n", 2)
		if len(parts) > 1 && (strings.Contains(parts[0], "translation") || strings.Contains(parts[0], "Translation")) {
			// If the first part looks like an explanation, return the second part
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return translation, nil
}
