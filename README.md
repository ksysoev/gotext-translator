# Gotext Translator

A CLI utility for automatically translating untranslated strings in gotext localization files using LLM (Language Model) technology. This tool helps streamline the localization process by identifying empty translations and using AI to provide high-quality translations while preserving all formatting and placeholders.

## Features

- Processes gotext JSON format files
- Identifies and translates only untranslated strings (empty translation field)
- Uses OpenAI GPT models for translation
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

Create a YAML or JSON configuration file (e.g., translator-config.yaml):

```yaml
llm:
  provider: openai
  api_key: your-api-key-here
  model: gpt-3.5-turbo
```

### Environment Variables

Instead of using a configuration file, you can set the following environment variables:

- `LLM_PROVIDER`: LLM provider (currently only "openai" is supported)
- `LLM_API_KEY`: API key for the LLM provider
- `LLM_MODEL`: Model name (e.g., "gpt-3.5-turbo")

Example:
```bash
export LLM_PROVIDER=openai
export LLM_API_KEY=your-api-key-here
export LLM_MODEL=gpt-3.5-turbo
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
