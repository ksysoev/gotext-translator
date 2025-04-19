package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type LLMConfig struct {
	Provider string            `mapstructure:"provider"`
	APIKey   string            `mapstructure:"api_key"`
	Model    string            `mapstructure:"model"`
	Options  map[string]string `mapstructure:"options"`
}

type Config struct {
	LLM LLMConfig `mapstructure:"llm"`
}

// initConfig initializes the configuration by reading from the specified config file.
func initConfig(arg *args) (*Config, error) {
	v := viper.New()

	if arg.ConfigPath != "" {
		v.SetConfigFile(arg.ConfigPath)

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults if not provided
	if cfg.LLM.Provider == "" {
		cfg.LLM.Provider = "openai"
	}
	if cfg.LLM.Model == "" {
		cfg.LLM.Model = "gpt-3.5-turbo"
	}
	if cfg.LLM.Options == nil {
		cfg.LLM.Options = make(map[string]string)
	}

	return &cfg, nil
}
