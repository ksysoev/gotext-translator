package translator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const openRouterBaseURL = "https://openrouter.ai/api/v1/chat/completions"

// OpenRouterProvider provides translation using OpenRouter
type OpenRouterProvider struct{}

// GetName returns the name of the provider
func (p *OpenRouterProvider) GetName() string {
	return "openrouter"
}

// CreateTranslator creates a translator instance
func (p *OpenRouterProvider) CreateTranslator(config map[string]interface{}) (Translator, error) {
	apiKey, ok := config["api_key"].(string)
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("OpenRouter API key is required")
	}

	model, ok := config["model"].(string)
	if !ok || model == "" {
		model = "openai/gpt-3.5-turbo" // Default model
	}

	return &OpenRouterTranslator{
		apiKey: apiKey,
		model:  model,
	}, nil
}

// OpenRouterTranslator implements the Translator interface using OpenRouter
type OpenRouterTranslator struct {
	apiKey string
	model  string
}

// OpenRouterRequest represents a request to OpenRouter API
type OpenRouterRequest struct {
	Model    string              `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
}

// OpenRouterMessage represents a chat message in OpenRouter API
type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents a response from OpenRouter API
type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Translate translates text to the specified target language
func (t *OpenRouterTranslator) Translate(ctx context.Context, text string, targetLang string) (string, error) {
	client := &http.Client{}

	requestBody := OpenRouterRequest{
		Model: t.model,
		Messages: []OpenRouterMessage{
			{
				Role:    "system",
				Content: "You are a professional translator. Your task is to translate text accurately while preserving all formatting, placeholders, and special characters.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Translate the following text to %s. Preserve any formatting, placeholders, and special characters:\n\n%s", targetLang, text),
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", openRouterBaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/ksysoev/gotext-translator")
	req.Header.Set("X-Title", "Gotext Translator")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response OpenRouterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("OpenRouter API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no translation choices returned from OpenRouter")
	}

	return response.Choices[0].Message.Content, nil
}
