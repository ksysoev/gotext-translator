package translator_test

import (
	"testing"

	"github.com/ksysoev/gotext-translator/pkg/translator"
	"github.com/stretchr/testify/assert"
)

func TestAnthropicProvider_GetName(t *testing.T) {
	provider := &translator.AnthropicProvider{}
	assert.Equal(t, "anthropic", provider.GetName())
}

func TestAnthropicProvider_CreateTranslator(t *testing.T) {
	provider := &translator.AnthropicProvider{}

	// Test with valid config
	config := map[string]interface{}{
		"api_key": "test-key",
		"model":   "test-model",
	}

	translator, err := provider.CreateTranslator(config)
	assert.NoError(t, err)
	assert.NotNil(t, translator)

	// Test with missing API key
	invalidConfig := map[string]interface{}{
		"model": "test-model",
	}

	translator, err = provider.CreateTranslator(invalidConfig)
	assert.Error(t, err)
	assert.Nil(t, translator)
	assert.Contains(t, err.Error(), "API key is required")

	// Test with empty API key
	emptyKeyConfig := map[string]interface{}{
		"api_key": "",
		"model":   "test-model",
	}

	translator, err = provider.CreateTranslator(emptyKeyConfig)
	assert.Error(t, err)
	assert.Nil(t, translator)
	assert.Contains(t, err.Error(), "API key is required")

	// Test with missing model (should use default)
	noModelConfig := map[string]interface{}{
		"api_key": "test-key",
	}

	translator, err = provider.CreateTranslator(noModelConfig)
	assert.NoError(t, err)
	assert.NotNil(t, translator)
}
