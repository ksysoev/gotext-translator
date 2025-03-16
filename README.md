# Gotext Translator

A CLI utility for automatically translating untranslated strings in gotext localization files using LLM (Language Model) technology. This tool helps streamline the localization process by identifying empty translations and using AI to provide high-quality translations while preserving all formatting and placeholders.

## Features

- Processes gotext JSON format files
- Identifies and translates only untranslated strings (empty translation field)
- Supports multiple LLM providers (OpenAI and Anthropic)
- Preserves JSON structure, placeholders, and special formatting
- Configurable via file or environment variables
- Interactive translation process with progress tracking
- Supports custom output paths

## Installation

Requires Go 1.23 or later.

```bash
go install github.com/ksysoev/gotext-translator/cmd/gotext-translate@latest
```

## Usage

### Basic Command Structure

```bash
gotext-translate translate [flags]
```

### Available Flags

- `--source`: Path to the source gotext JSON file (required)
- `--target-lang`: Target language code (e.g., ru-RU) (required)
- `--output`: Output file path (optional, defaults to out.gotext.json in source directory)
- `--config`: Path to configuration file (optional)
- `--loglevel`: Log level (debug, info, warn, error) (default: info)
- `--logtext`: Use text format for logs instead of JSON (default: false)

### Examples

1. Basic translation:
```bash
gotext-translate translate --source locales/en-US/messages.gotext.json --target-lang ru-RU
```

2. Specify custom output path:
```bash
gotext-translate translate --source locales/en-US/messages.gotext.json --target-lang ru-RU --output locales/ru-RU/messages.gotext.json
```

3. Use configuration file:
```bash
gotext-translate translate --config translator-config.yaml --source locales/en-US/messages.gotext.json --target-lang ru-RU
```

4. Enable debug logging:
```bash
gotext-translate translate --source locales/en-US/messages.gotext.json --target-lang ru-RU --loglevel debug
```

## Configuration

### Configuration File

Create a YAML or JSON configuration file (e.g., translator-config.yaml). You can configure either OpenAI or Anthropic as your LLM provider:

```yaml
# For OpenAI:
llm:
  provider: openai
  api_key: your-openai-api-key
  model: gpt-3.5-turbo  # or gpt-4, etc.
```

```yaml
# For Anthropic:
llm:
  provider: anthropic
  api_key: your-anthropic-api-key
  model: claude-3-haiku  # or claude-3-opus, claude-3-sonnet, etc.
```

### Environment Variables

Instead of using a configuration file, you can set the following environment variables:

- `LLM_PROVIDER`: LLM provider ("openai" or "anthropic")
- `LLM_API_KEY`: API key for the LLM provider
- `LLM_MODEL`: Model name (e.g., "gpt-3.5-turbo" for OpenAI or "claude-3-haiku" for Anthropic)

Example for OpenAI:
```bash
export LLM_PROVIDER=openai
export LLM_API_KEY=your-openai-api-key
export LLM_MODEL=gpt-3.5-turbo
```

Example for Anthropic:
```bash
export LLM_PROVIDER=anthropic
export LLM_API_KEY=your-anthropic-api-key
export LLM_MODEL=claude-3-haiku
```

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

## Error Handling

- The tool validates input files before processing
- Translation errors for individual strings don't stop the entire process
- Progress is logged for monitoring and debugging
- Detailed error messages help identify and resolve issues

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see the [LICENSE](LICENSE) file for details
