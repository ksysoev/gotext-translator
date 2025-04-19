package translator

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider provides translation using OpenAI
type OpenAIProvider struct{}

// GetName returns the name of the provider
func (p *OpenAIProvider) GetName() string {
	return "openai"
}

// CreateTranslator creates a translator instance
func (p *OpenAIProvider) CreateTranslator(config map[string]interface{}) (Translator, error) {
	apiKey, ok := config["api_key"].(string)
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	model, ok := config["model"].(string)
	if !ok || model == "" {
		model = "gpt-3.5-turbo" // Default model
	}

	client := openai.NewClient(apiKey)
	return &OpenAITranslator{
		client: client,
		model:  model,
	}, nil
}

// OpenAITranslator implements the Translator interface using OpenAI
type OpenAITranslator struct {
	client *openai.Client
	model  string
}

// Translate translates text to the specified target language
func (t *OpenAITranslator) Translate(ctx context.Context, text string, targetLang string) (string, error) {
	resp, err := t.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: t.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a professional translator. Your task is to translate text accurately while preserving all formatting, placeholders, and special characters.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Translate the following text to %s. Preserve any formatting, placeholders, and special characters:\n\n%s", targetLang, text),
				},
			},
			Temperature: 0.3, // Lower temperature for more consistent translations
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to get translation from OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no translation choices returned from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
