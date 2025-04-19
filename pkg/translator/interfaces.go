package translator

import (
	"context"
)

// Translator defines the interface for translating text
type Translator interface {
	// Translate translates text to the specified target language
	Translate(ctx context.Context, text string, targetLang string) (string, error)
}

// Provider defines the interface for LLM providers
type Provider interface {
	// GetName returns the name of the provider
	GetName() string
	// CreateTranslator creates a translator instance
	CreateTranslator(config map[string]interface{}) (Translator, error)
}

// Factory defines the interface for creating translators
type Factory interface {
	// GetProviders returns all registered providers
	GetProviders() []string
	// RegisterProvider registers a provider with the factory
	RegisterProvider(provider Provider) error
	// CreateTranslator creates a translator for the specified provider
	CreateTranslator(providerName string, config map[string]interface{}) (Translator, error)
}
