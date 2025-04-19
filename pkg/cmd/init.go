package cmd

import (
	"fmt"

	"github.com/ksysoev/gotext-translator/pkg/translator"
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
	cmd.AddCommand(translateDirCommand(args))
	cmd.AddCommand(providersCommand())

	cmd.PersistentFlags().StringVar(&args.ConfigPath, "config", "", "config file path")
	cmd.PersistentFlags().StringVar(&args.LogLevel, "loglevel", "info", "log level (debug, info, warn, error)")
	cmd.PersistentFlags().BoolVar(&args.TextFormat, "logtext", false, "log in text format, otherwise JSON")
	cmd.PersistentFlags().BoolVar(&args.ForceRewrite, "force-rewrite", false, "force rewrite existing translations")

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

// translateDirCommand creates a cobra.Command to translate all files in a directory
func translateDirCommand(args *args) *cobra.Command {
	SetArgs(args) // Store args globally for translation use
	cmd := &cobra.Command{
		Use:   "translate-dir",
		Short: "Translate all files in a directory",
		Long:  "Translate all gotext localization files in a directory structure",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := initLogger(args); err != nil {
				return fmt.Errorf("failed to initialize logger: %w", err)
			}

			if args.SourceDir == "" {
				return fmt.Errorf("source directory path is required")
			}

			if args.TargetLang == "" {
				return fmt.Errorf("target language is required")
			}

			cfg, err := initConfig(args)
			if err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}

			return runDirectoryTranslation(cmd.Context(), cfg)
		},
	}

	cmd.Flags().StringVar(&args.SourceDir, "dir", "", "source directory path containing localization files")
	cmd.Flags().StringVar(&args.TargetLang, "target-lang", "", "target language (e.g., ru-RU)")

	return cmd
}

// providersCommand creates a cobra.Command to list available providers
func providersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "List available translation providers",
		Long:  "List all registered translation providers",
		Run: func(cmd *cobra.Command, args []string) {
			factory := translator.NewFactory()
			translator.RegisterProviders(factory)

			providers := factory.GetProviders()
			if len(providers) == 0 {
				fmt.Println("No translation providers registered")
				return
			}

			fmt.Println("Available translation providers:")
			for _, provider := range providers {
				fmt.Printf("  - %s\n", provider)
			}

			fmt.Println("\nYou can configure the provider in your config file or via environment variables.")
			fmt.Println("See the README for more information.")
		},
	}

	return cmd
}
