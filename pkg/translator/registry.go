package translator

import "log/slog"

// RegisterProviders registers all available providers with the factory
func RegisterProviders(factory Factory) {
	providers := []Provider{
		&OpenAIProvider{},
		&OpenRouterProvider{},
		// Future providers to be added:
		// &LangChainProvider{},
		// &AnthropicProvider{},
	}

	for _, provider := range providers {
		if err := factory.RegisterProvider(provider); err != nil {
			slog.Error("failed to register provider",
				slog.String("provider", provider.GetName()),
				slog.Any("error", err),
			)
		} else {
			slog.Debug("registered provider", slog.String("provider", provider.GetName()))
		}
	}
}
