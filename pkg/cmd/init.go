package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type args struct {
	version      string
	LogLevel     string
	ConfigPath   string
	SourcePath   string
	SourceDir    string
	TargetLang   string
	OutputPath   string
	TextFormat   bool
	ForceRewrite bool
}

// InitCommands initializes and returns the root command for the application.
func InitCommands(version string) (*cobra.Command, error) {
	args := &args{
		version: version,
	}

	cmd := &cobra.Command{
		Use:   "gotext-translate",
		Short: "Translate untranslated strings in gotext localization files",
		Long:  "A CLI utility to translate untranslated strings in gotext localization files using LLM",
	}

	cmd.AddCommand(translateCommand(args))

	cmd.PersistentFlags().StringVar(&args.ConfigPath, "config", "", "config file path")
	cmd.PersistentFlags().StringVar(&args.LogLevel, "loglevel", "info", "log level (debug, info, warn, error)")
	cmd.PersistentFlags().BoolVar(&args.TextFormat, "logtext", false, "log in text format, otherwise JSON")

	return cmd, nil
}

// translateCommand creates a cobra.Command to translate a single file
func translateCommand(args *args) *cobra.Command {
	SetArgs(args) // Store args globally for translation use
	cmd := &cobra.Command{
		Use:   "translate",
		Short: "Translate a single file",
		Long:  "Translate a single gotext localization file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := initLogger(args); err != nil {
				return fmt.Errorf("failed to initialize logger: %w", err)
			}

			if args.SourcePath == "" {
				return fmt.Errorf("source file path is required")
			}

			if args.TargetLang == "" {
				return fmt.Errorf("target language is required")
			}

			cfg, err := initConfig(args)
			if err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}

			return runTranslation(cmd.Context(), cfg)
		},
	}

	cmd.Flags().StringVar(&args.SourcePath, "source", "", "source file path")
	cmd.Flags().StringVar(&args.TargetLang, "target-lang", "", "target language (e.g., ru-RU)")
	cmd.Flags().StringVar(&args.OutputPath, "output", "", "output file path (optional)")

	return cmd
}
