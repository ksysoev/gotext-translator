package translator

import (
	"fmt"
	"sync"
)

// DefaultFactory is the standard implementation of Factory
type DefaultFactory struct {
	providers map[string]Provider
	mu        sync.RWMutex
}

// NewFactory creates a new factory instance
func NewFactory() Factory {
	return &DefaultFactory{
		providers: make(map[string]Provider),
	}
}

// GetProviders returns all registered providers
func (f *DefaultFactory) GetProviders() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	providers := make([]string, 0, len(f.providers))
	for name := range f.providers {
		providers = append(providers, name)
	}
	return providers
}

// RegisterProvider registers a provider with the factory
func (f *DefaultFactory) RegisterProvider(provider Provider) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	name := provider.GetName()
	if _, exists := f.providers[name]; exists {
		return fmt.Errorf("provider %s already registered", name)
	}

	f.providers[name] = provider
	return nil
}

// CreateTranslator creates a translator for the specified provider
func (f *DefaultFactory) CreateTranslator(providerName string, config map[string]interface{}) (Translator, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	provider, exists := f.providers[providerName]
	if !exists {
		return nil, fmt.Errorf("provider %s not registered", providerName)
	}

	return provider.CreateTranslator(config)
}
