package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ksysoev/gotext-translator/pkg/translator"
)

type GotextMessage struct {
	ID           string `json:"id"`
	Message      string `json:"message"`
	Translation  string `json:"translation"`
	Placeholders []struct {
		ID             string `json:"id"`
		String         string `json:"string"`
		Type           string `json:"type"`
		UnderlyingType string `json:"underlyingType"`
		Expr           string `json:"expr"`
		ArgNum         int    `json:"argNum"`
	} `json:"placeholders,omitempty"`
}

type GotextFile struct {
	Language string          `json:"language"`
	Messages []GotextMessage `json:"messages"`
}

func runTranslation(ctx context.Context, cfg *Config) error {
	sourceData, err := os.ReadFile(globalArgs.SourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	var gotextFile GotextFile
	if err := json.Unmarshal(sourceData, &gotextFile); err != nil {
		return fmt.Errorf("failed to parse source file: %w", err)
	}

	gotextFile.Language = globalArgs.TargetLang

	// Initialize translator factory
	factory := translator.NewFactory()
	translator.RegisterProviders(factory)

	// Create translator instance
	config := map[string]interface{}{
		"api_key": cfg.LLM.APIKey,
		"model":   cfg.LLM.Model,
	}

	// Add any additional options from config
	for k, v := range cfg.LLM.Options {
		config[k] = v
	}

	trans, err := factory.CreateTranslator(cfg.LLM.Provider, config)
	if err != nil {
		return fmt.Errorf("failed to initialize translator: %w", err)
	}

	// Process each message
	slog.Info("starting translation", slog.Int("total_messages", len(gotextFile.Messages)))
	for i := range gotextFile.Messages {
		msg := &gotextFile.Messages[i]
		if msg.Translation != "" {
			slog.Debug("skipping translated message", slog.String("id", msg.ID))
			continue
		}

		translation, err := trans.Translate(ctx, msg.Message, gotextFile.Language)
		if err != nil {
			slog.Error("failed to translate message",
				slog.String("id", msg.ID),
				slog.String("error", err.Error()))
			continue
		}

		msg.Translation = translation
		slog.Info("translated message",
			slog.String("id", msg.ID),
			slog.String("original", msg.Message),
			slog.String("translation", translation))
	}

	// Determine output path
	outputPath := globalArgs.OutputPath
	if outputPath == "" {
		dir := filepath.Dir(globalArgs.SourcePath)
		outputPath = filepath.Join(dir, "out.gotext.json")
	}

	// Save result
	output, err := json.MarshalIndent(gotextFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	slog.Info("translation completed", slog.String("output", outputPath))
	return nil
}

var globalArgs *args // Store args globally for use in translation

// SetArgs stores the command arguments for use in translation
func SetArgs(a *args) {
	globalArgs = a
}
