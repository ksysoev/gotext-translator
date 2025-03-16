package cmd

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
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

// AnthropicTranslator implements the Translator interface using Anthropic's API
type AnthropicTranslator struct {
	client *anthropic.Client
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

	case "anthropic":
		if cfg.LLM.APIKey == "" {
			return nil, fmt.Errorf("Anthropic API key is required")
		}

		client := anthropic.NewClient(option.WithAPIKey(cfg.LLM.APIKey))
		return &AnthropicTranslator{
			client: client,
			model:  cfg.LLM.Model,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.LLM.Provider)
	}
}

// Translate implements the Translator interface for OpenAI
func (t *OpenAITranslator) Translate(ctx context.Context, text, targetLang string) (string, error) {
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

// Translate implements the Translator interface for Anthropic
func (t *AnthropicTranslator) Translate(ctx context.Context, text, targetLang string) (string, error) {
	msg, err := t.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(t.model),
		MaxTokens: anthropic.F(int64(1024)),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock("You are a professional translator. Your task is to translate text accurately while preserving all formatting, placeholders, and special characters."),
		}),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(fmt.Sprintf("Please translate the following text to %s, preserving all formatting and placeholders:\n\n%s", targetLang, text))),
		}),
	})

	if err != nil {
		return "", fmt.Errorf("failed to get translation from Anthropic: %w", err)
	}

	if len(msg.Content) == 0 {
		return "", fmt.Errorf("empty response from Anthropic API")
	}

	return msg.Content[0].Text, nil
}
