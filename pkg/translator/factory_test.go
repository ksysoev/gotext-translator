package translator_test

import (
	"fmt"
	"testing"

	"github.com/ksysoev/gotext-translator/pkg/translator"
	"github.com/ksysoev/gotext-translator/pkg/translator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefaultFactory_GetProviders(t *testing.T) {
	factory := translator.NewFactory()

	// Register test providers
	provider1 := new(mocks.Provider)
	provider1.On("GetName").Return("provider1")

	provider2 := new(mocks.Provider)
	provider2.On("GetName").Return("provider2")

	// Register providers
	err := factory.RegisterProvider(provider1)
	assert.NoError(t, err)

	err = factory.RegisterProvider(provider2)
	assert.NoError(t, err)

	// Get providers
	providers := factory.GetProviders()
	assert.Len(t, providers, 2)
	assert.Contains(t, providers, "provider1")
	assert.Contains(t, providers, "provider2")
}

func TestDefaultFactory_RegisterProvider(t *testing.T) {
	factory := translator.NewFactory()

	// Register a provider
	provider := new(mocks.Provider)
	provider.On("GetName").Return("test-provider")

	err := factory.RegisterProvider(provider)
	assert.NoError(t, err)

	// Try to register the same provider again
	provider2 := new(mocks.Provider)
	provider2.On("GetName").Return("test-provider")

	err = factory.RegisterProvider(provider2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestDefaultFactory_CreateTranslator(t *testing.T) {
	factory := translator.NewFactory()

	// Mock translator
	mockTranslator := new(mocks.Translator)

	// Register a provider
	provider := new(mocks.Provider)
	provider.On("GetName").Return("test-provider")
	provider.On("CreateTranslator", mock.Anything).Return(mockTranslator, nil)

	err := factory.RegisterProvider(provider)
	assert.NoError(t, err)

	// Create translator
	config := map[string]interface{}{
		"api_key": "test-key",
		"model":   "test-model",
	}

	translatorInstance, err := factory.CreateTranslator("test-provider", config)
	assert.NoError(t, err)
	assert.Equal(t, mockTranslator, translatorInstance)

	// Try to create translator for non-existent provider
	_, err = factory.CreateTranslator("non-existent", config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not registered")

	// Mock error from provider
	errorProvider := new(mocks.Provider)
	errorProvider.On("GetName").Return("error-provider")
	errorProvider.On("CreateTranslator", mock.Anything).Return(nil, fmt.Errorf("test error"))

	err = factory.RegisterProvider(errorProvider)
	assert.NoError(t, err)

	_, err = factory.CreateTranslator("error-provider", config)
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}
