package cmd

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Translator interface defines methods for translating text
type Translator interface {
	Translate(ctx context.Context, text, targetLang string) (string, error)
}

// OpenAITranslator implements the Translator interface using OpenAI's API
type OpenAITranslator struct {
	client *openai.Client
	model  string
}

func initTranslator(cfg *Config) (Translator, error) {
	switch cfg.LLM.Provider {
	case "openai":
		if cfg.LLM.APIKey == "" {
			return nil, fmt.Errorf("OpenAI API key is required")
		}

		client := openai.NewClient(cfg.LLM.APIKey)
		return &OpenAITranslator{
			client: client,
			model:  cfg.LLM.Model,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.LLM.Provider)
	}
}

// Translate implements the Translator interface for OpenAI
func (t *OpenAITranslator) Translate(ctx context.Context, text, targetLang string) (string, error) {
	prompt := fmt.Sprintf("Translate the following text to %s. Preserve any formatting, placeholders, and special characters:\n\n%s", targetLang, text)

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
					Content: prompt,
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
