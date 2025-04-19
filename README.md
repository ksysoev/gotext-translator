# Gotext Translator

A CLI utility for automatically translating untranslated strings in gotext localization files using LLM (Language Model) technology. This tool helps streamline the localization process by identifying empty translations and using AI to provide high-quality translations while preserving all formatting and placeholders.

## Features

- Processes gotext JSON format files
- Identifies and translates only untranslated strings (empty translation field)
- Model-agnostic architecture with support for multiple LLM providers
- Current providers: OpenAI and OpenRouter (with more planned)
- Preserves JSON structure, placeholders, and special formatting
- Configurable via file or environment variables
- Interactive translation process with progress tracking
- Supports custom output paths
- Supports batch translation of entire directory structures

## Installation

Requires Go 1.23 or later.

```bash
go install github.com/ksysoev/gotext-translator/cmd/gotext-translate@latest
```

## Usage

### Basic Command Structure

```bash
# Translate a single file
gotext-translate translate [flags]

# Translate all files in a directory structure
gotext-translate translate-dir [flags]
```

### Available Flags

Common flags:
- `--config`: Path to configuration file (optional)
- `--loglevel`: Log level (debug, info, warn, error) (default: info)
- `--logtext`: Use text format for logs instead of JSON (default: false)
- `--force-rewrite`: Force rewrite existing translations (default: false)

Translate command flags:
- `--source`: Path to the source gotext JSON file (required)
- `--target-lang`: Target language code (e.g., ru-RU) (required)
- `--output`: Output file path (optional, defaults to out.gotext.json in source directory)

Translate-dir command flags:
- `--dir`: Path to the source directory containing localization files (required)
- `--target-lang`: Target language code (e.g., ru-RU) (required)

### Examples

1. Basic single file translation:
```bash
gotext-translate translate --source locales/en-US/messages.gotext.json --target-lang ru-RU
```

2. Specify custom output path:
```bash
gotext-translate translate --source locales/en-US/messages.gotext.json --target-lang ru-RU --output locales/ru-RU/messages.gotext.json
```

3. Translate all files in a directory structure:
```bash
gotext-translate translate-dir --dir samples --target-lang fr-FR
```

4. Force rewrite existing translations:
```bash
gotext-translate translate-dir --dir samples --target-lang ru-RU --force-rewrite
```

5. Use configuration file:
```bash
gotext-translate translate-dir --config translator-config.yaml --dir samples --target-lang ru-RU
```

## Configuration

### Configuration File

Create a YAML or JSON configuration file (e.g., translator-config.yaml). You can configure different LLM providers:

```yaml
# For OpenAI:
llm:
  provider: openai
  api_key: your-openai-api-key
  model: gpt-3.5-turbo  # or gpt-4, etc.
```

```yaml
# For OpenRouter:
llm:
  provider: openrouter
  api_key: your-openrouter-api-key
  model: openai/gpt-3.5-turbo  # or anthropic/claude-3-opus, mistralai/mistral-tiny, etc.
  options:
    route_prefix: gotext-translator
```

### Environment Variables

Instead of using a configuration file, you can set the following environment variables:

- `LLM_PROVIDER`: LLM provider ("openai" or "openrouter")
- `LLM_API_KEY`: API key for the LLM provider
- `LLM_MODEL`: Model name (e.g., "gpt-3.5-turbo" for OpenAI or "anthropic/claude-3-haiku" for OpenRouter)

Example for OpenAI:
```bash
export LLM_PROVIDER=openai
export LLM_API_KEY=your-openai-api-key
export LLM_MODEL=gpt-3.5-turbo
```

Example for OpenRouter:
```bash
export LLM_PROVIDER=openrouter
export LLM_API_KEY=your-openrouter-api-key
export LLM_MODEL=anthropic/claude-3-haiku
```

## Directory Structure

When using the `translate-dir` command, the tool expects a directory structure like this:

```
samples/
└── locales/
    ├── en-GB/
    │   └── messages.gotext.json
    └── ru-RU/
        └── messages.gotext.json
```

The tool will:
1. Look for the first non-target language directory as the source (e.g., en-GB)
2. Find all .gotext.json files in the source directory
3. Create or update corresponding files in the target language directory
4. Translate all untranslated strings

## Input File Format

The tool expects gotext JSON files in the following format:

```json
{
  "language": "en-US",
  "messages": [
    {
      "id": "greeting",
      "message": "Hello, World!",
      "translation": "",
      "placeholders": []
    }
  ]
}
```

## Output

The tool generates a new JSON file with translations added:

```json
{
  "language": "ru-RU",
  "messages": [
    {
      "id": "greeting",
      "message": "Hello, World!",
      "translation": "Привет, Мир!",
      "placeholders": []
    }
  ]
}
```

## Architecture

The tool is built using SOLID principles and follows a modular design:
- Factory pattern for creating translator instances
- Strategy pattern for different translation providers
- Interface-based design for easy extension

## Error Handling

- The tool validates input files before processing
- Translation errors for individual strings don't stop the entire process
- Progress is logged for monitoring and debugging
- Detailed error messages help identify and resolve issues

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development

1. Clone the repository
2. Install dependencies: `go mod download`
3. Build the project: `go build ./cmd/gotext-translate`
4. Run tests: `go test ./...`

## License

MIT License - see the [LICENSE](LICENSE) file for details