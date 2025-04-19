package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ksysoev/gotext-translator/pkg/translator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessFile(t *testing.T) {
	// Initialize globalArgs
	globalArgs = &args{ForceRewrite: false}

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "gotext-translator-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock translator
	mockTranslator := new(mocks.Translator)
	mockTranslator.On("Translate", mock.Anything, "Hello, World!", "ru-RU").
		Return("Привет, Мир!", nil)
	mockTranslator.On("Translate", mock.Anything, "Welcome to the app!", "ru-RU").
		Return("Добро пожаловать в приложение!", nil)

	// Create a source file
	sourceFile := GotextFile{
		Language: "en-US",
		Messages: []GotextMessage{
			{
				ID:      "greeting",
				Message: "Hello, World!",
			},
			{
				ID:      "welcome",
				Message: "Welcome to the app!",
			},
		},
	}

	sourcePath := filepath.Join(tempDir, "messages.gotext.json")
	sourceData, err := json.MarshalIndent(sourceFile, "", "  ")
	assert.NoError(t, err)
	err = os.WriteFile(sourcePath, sourceData, 0644)
	assert.NoError(t, err)

	// Process the file
	targetPath := filepath.Join(tempDir, "out.gotext.json")
	count, err := processFile(context.Background(), mockTranslator, sourcePath, targetPath, "ru-RU")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	// Verify the output file
	targetData, err := os.ReadFile(targetPath)
	assert.NoError(t, err)
	var targetFile GotextFile
	err = json.Unmarshal(targetData, &targetFile)
	assert.NoError(t, err)

	assert.Equal(t, "ru-RU", targetFile.Language)
	assert.Len(t, targetFile.Messages, 2)
	assert.Equal(t, "greeting", targetFile.Messages[0].ID)
	assert.Equal(t, "Hello, World!", targetFile.Messages[0].Message)
	assert.Equal(t, "Привет, Мир!", targetFile.Messages[0].Translation)
	assert.Equal(t, "welcome", targetFile.Messages[1].ID)
	assert.Equal(t, "Welcome to the app!", targetFile.Messages[1].Message)
	assert.Equal(t, "Добро пожаловать в приложение!", targetFile.Messages[1].Translation)

	// Test with an existing file - only translate missing translations
	existingFile := GotextFile{
		Language: "ru-RU",
		Messages: []GotextMessage{
			{
				ID:          "greeting",
				Message:     "Hello, World!",
				Translation: "Привет, Мир!", // Already translated
			},
			{
				ID:      "welcome",
				Message: "Welcome to the app!",
				// Missing translation
			},
		},
	}

	existingPath := filepath.Join(tempDir, "existing.gotext.json")
	existingData, err := json.MarshalIndent(existingFile, "", "  ")
	assert.NoError(t, err)
	err = os.WriteFile(existingPath, existingData, 0644)
	assert.NoError(t, err)

	// Process the file
	count, err = processFile(context.Background(), mockTranslator, sourcePath, existingPath, "ru-RU")
	assert.NoError(t, err)
	assert.Equal(t, 1, count) // Only one message should be translated

	// Verify the output file
	existingData, err = os.ReadFile(existingPath)
	assert.NoError(t, err)
	var updatedFile GotextFile
	err = json.Unmarshal(existingData, &updatedFile)
	assert.NoError(t, err)

	assert.Equal(t, "ru-RU", updatedFile.Language)
	assert.Len(t, updatedFile.Messages, 2)
	assert.Equal(t, "greeting", updatedFile.Messages[0].ID)
	assert.Equal(t, "Hello, World!", updatedFile.Messages[0].Message)
	assert.Equal(t, "Привет, Мир!", updatedFile.Messages[0].Translation)
	assert.Equal(t, "welcome", updatedFile.Messages[1].ID)
	assert.Equal(t, "Welcome to the app!", updatedFile.Messages[1].Message)
	assert.Equal(t, "Добро пожаловать в приложение!", updatedFile.Messages[1].Translation)

	// Test with force rewrite
	globalArgs.ForceRewrite = true

	// Create updated expectations for force rewrite
	mockTranslator = new(mocks.Translator)
	mockTranslator.On("Translate", mock.Anything, "Hello, World!", "ru-RU").
		Return("Привет, Мир! (updated)", nil)
	mockTranslator.On("Translate", mock.Anything, "Welcome to the app!", "ru-RU").
		Return("Добро пожаловать в приложение! (updated)", nil)

	count, err = processFile(context.Background(), mockTranslator, sourcePath, existingPath, "ru-RU")
	assert.NoError(t, err)
	assert.Equal(t, 2, count) // Both messages should be translated

	// Verify the output file
	existingData, err = os.ReadFile(existingPath)
	assert.NoError(t, err)
	var forceRewriteFile GotextFile
	err = json.Unmarshal(existingData, &forceRewriteFile)
	assert.NoError(t, err)

	assert.Equal(t, "ru-RU", forceRewriteFile.Language)
	assert.Len(t, forceRewriteFile.Messages, 2)
	assert.Equal(t, "greeting", forceRewriteFile.Messages[0].ID)
	assert.Equal(t, "Hello, World!", forceRewriteFile.Messages[0].Message)
	assert.Equal(t, "Привет, Мир! (updated)", forceRewriteFile.Messages[0].Translation)
	assert.Equal(t, "Machine translated", forceRewriteFile.Messages[0].TranslatorComment)
	assert.Equal(t, "welcome", forceRewriteFile.Messages[1].ID)
	assert.Equal(t, "Welcome to the app!", forceRewriteFile.Messages[1].Message)
	assert.Equal(t, "Добро пожаловать в приложение! (updated)", forceRewriteFile.Messages[1].Translation)
	assert.Equal(t, "Machine translated", forceRewriteFile.Messages[1].TranslatorComment)
}

// Reset globalArgs after tests
func TestMain(m *testing.M) {
	// Setup
	originalArgs := globalArgs

	// Run tests
	code := m.Run()

	// Teardown
	globalArgs = originalArgs

	os.Exit(code)
}
