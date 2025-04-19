package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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
	TranslatorComment string `json:"translatorComment,omitempty"`
	Fuzzy             bool   `json:"fuzzy,omitempty"`
}

type GotextFile struct {
	Language string          `json:"language"`
	Messages []GotextMessage `json:"messages"`
}

// runTranslation handles translation of a single file
func runTranslation(ctx context.Context, cfg *Config) error {
	// Prepare the translator
	trans, err := prepareTranslator(ctx, cfg)
	if err != nil {
		return err
	}

	// Process the file
	sourceData, err := os.ReadFile(globalArgs.SourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	var gotextFile GotextFile
	if err := json.Unmarshal(sourceData, &gotextFile); err != nil {
		return fmt.Errorf("failed to parse source file: %w", err)
	}

	gotextFile.Language = globalArgs.TargetLang

	// Process each message
	slog.Info("starting translation",
		slog.String("file", globalArgs.SourcePath),
		slog.String("target_lang", globalArgs.TargetLang),
		slog.Int("total_messages", len(gotextFile.Messages)),
	)

	processedCount := 0
	for i := range gotextFile.Messages {
		msg := &gotextFile.Messages[i]
		if msg.Translation != "" && !globalArgs.ForceRewrite {
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
		// If this was a forced rewrite, add a comment
		if globalArgs.ForceRewrite && msg.Translation != "" {
			msg.TranslatorComment = "Machine translated"
		}

		processedCount++
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

	slog.Info("translation completed",
		slog.String("file", globalArgs.SourcePath),
		slog.String("output", outputPath),
		slog.Int("processed", processedCount),
	)

	return nil
}

// runDirectoryTranslation handles translation of all files in a directory
func runDirectoryTranslation(ctx context.Context, cfg *Config) error {
	// Prepare the translator
	trans, err := prepareTranslator(ctx, cfg)
	if err != nil {
		return err
	}

	// Find base language directory (usually en-US, en-GB, etc.)
	baseDir := filepath.Join(globalArgs.SourceDir, "locales")
	targetDir := filepath.Join(baseDir, globalArgs.TargetLang)

	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Find all source language directories
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return fmt.Errorf("failed to read base directory: %w", err)
	}

	var sourceEntries []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != globalArgs.TargetLang {
			sourceEntries = append(sourceEntries, entry)
		}
	}

	if len(sourceEntries) == 0 {
		return fmt.Errorf("no source language directories found in %s", baseDir)
	}

	// Automatically choose the first source directory
	sourceSubdir := sourceEntries[0].Name()
	sourceLangDir := filepath.Join(baseDir, sourceSubdir)

	slog.Info("starting directory translation",
		slog.String("source_dir", sourceLangDir),
		slog.String("target_dir", targetDir),
		slog.String("target_lang", globalArgs.TargetLang),
	)

	// Find all .gotext.json files in the source language directory
	var sourceFiles []string
	if err := filepath.Walk(sourceLangDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".gotext.json") && !strings.HasSuffix(path, "out.gotext.json") {
			sourceFiles = append(sourceFiles, path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to find source files: %w", err)
	}

	slog.Info("found source files", slog.Int("count", len(sourceFiles)))

	// Process each source file
	totalProcessedFiles := 0
	totalProcessedMessages := 0

	for _, sourceFile := range sourceFiles {
		// Determine target file path
		relPath, err := filepath.Rel(sourceLangDir, sourceFile)
		if err != nil {
			slog.Error("failed to get relative path", slog.String("file", sourceFile), slog.String("error", err.Error()))
			continue
		}

		targetFile := filepath.Join(targetDir, relPath)

		// Create parent directories if they don't exist
		if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
			slog.Error("failed to create target directory", slog.String("dir", filepath.Dir(targetFile)), slog.String("error", err.Error()))
			continue
		}

		// Process the file
		processedCount, err := processFile(ctx, trans, sourceFile, targetFile, globalArgs.TargetLang)
		if err != nil {
			slog.Error("failed to process file", slog.String("file", sourceFile), slog.String("error", err.Error()))
			continue
		}

		totalProcessedFiles++
		totalProcessedMessages += processedCount
	}

	slog.Info("directory translation completed",
		slog.String("target_lang", globalArgs.TargetLang),
		slog.Int("processed_files", totalProcessedFiles),
		slog.Int("processed_messages", totalProcessedMessages),
	)

	return nil
}

// processFile processes a single gotext file
func processFile(ctx context.Context, trans translator.Translator, sourcePath, targetPath, targetLang string) (int, error) {
	// Read source file
	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read source file: %w", err)
	}

	var sourceFile GotextFile
	if err := json.Unmarshal(sourceData, &sourceFile); err != nil {
		return 0, fmt.Errorf("failed to parse source file: %w", err)
	}

	// Create target file or read existing one
	var targetFile GotextFile
	targetExists := false

	if _, err := os.Stat(targetPath); err == nil {
		// Target file exists, read it
		targetData, err := os.ReadFile(targetPath)
		if err != nil {
			return 0, fmt.Errorf("failed to read target file: %w", err)
		}

		if err := json.Unmarshal(targetData, &targetFile); err != nil {
			return 0, fmt.Errorf("failed to parse target file: %w", err)
		}

		targetExists = true
	} else {
		// Create new target file structure
		targetFile = GotextFile{
			Language: targetLang,
			Messages: make([]GotextMessage, len(sourceFile.Messages)),
		}

		// Copy message structure from source
		for i, msg := range sourceFile.Messages {
			targetFile.Messages[i] = GotextMessage{
				ID:           msg.ID,
				Message:      msg.Message,
				Placeholders: msg.Placeholders,
			}
		}
	}

	// Create map from message ID to index for quick lookup
	targetMsgMap := make(map[string]int)
	for i, msg := range targetFile.Messages {
		targetMsgMap[msg.ID] = i
	}

	// Process each message
	slog.Info("processing file",
		slog.String("source", sourcePath),
		slog.String("target", targetPath),
		slog.Int("total_messages", len(sourceFile.Messages)),
	)

	processedCount := 0

	for _, srcMsg := range sourceFile.Messages {
		// Find or create target message
		targetIdx, exists := targetMsgMap[srcMsg.ID]
		if !exists {
			// Message doesn't exist in target file, add it
			targetIdx = len(targetFile.Messages)
			targetFile.Messages = append(targetFile.Messages, GotextMessage{
				ID:           srcMsg.ID,
				Message:      srcMsg.Message,
				Placeholders: srcMsg.Placeholders,
			})
			targetMsgMap[srcMsg.ID] = targetIdx
		} else {
			// Update message text from source if it changed
			targetMsg := &targetFile.Messages[targetIdx]
			targetMsg.Message = srcMsg.Message
			targetMsg.Placeholders = srcMsg.Placeholders
		}

		targetMsg := &targetFile.Messages[targetIdx]

		// Skip if already translated and not forced to rewrite
		if targetMsg.Translation != "" && !globalArgs.ForceRewrite {
			slog.Debug("skipping translated message", slog.String("id", targetMsg.ID))
			continue
		}

		// Translate the message
		translation, err := trans.Translate(ctx, targetMsg.Message, targetLang)
		if err != nil {
			slog.Error("failed to translate message",
				slog.String("id", targetMsg.ID),
				slog.String("error", err.Error()))
			continue
		}

		targetMsg.Translation = translation
		// If this was a forced rewrite, add a comment
		if globalArgs.ForceRewrite && targetMsg.Translation != "" {
			targetMsg.TranslatorComment = "Machine translated"
		}

		processedCount++
		slog.Info("translated message",
			slog.String("id", targetMsg.ID),
			slog.String("original", targetMsg.Message),
			slog.String("translation", translation))
	}

	// Save the target file
	output, err := json.MarshalIndent(targetFile, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("failed to marshal output: %w", err)
	}

	if err := os.WriteFile(targetPath, output, 0644); err != nil {
		return 0, fmt.Errorf("failed to write output file: %w", err)
	}

	slog.Info("file processing completed",
		slog.String("file", targetPath),
		slog.Bool("new_file", !targetExists),
		slog.Int("processed", processedCount),
	)

	return processedCount, nil
}

// prepareTranslator creates and initializes a translator
func prepareTranslator(ctx context.Context, cfg *Config) (translator.Translator, error) {
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
		return nil, fmt.Errorf("failed to initialize translator: %w", err)
	}

	return trans, nil
}

var globalArgs *args // Store args globally for translation use

// SetArgs stores the command arguments for use in translation
func SetArgs(a *args) {
	globalArgs = a
}
