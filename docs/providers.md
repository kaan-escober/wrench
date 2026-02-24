# Supported Providers

droid-cfg supports 12 built-in providers plus any custom OpenAI-compatible or Anthropic-compatible endpoint.

## Built-in Providers

| Provider | Type | Base URL | Auth |
|----------|------|----------|------|
| OpenRouter | `generic-chat-completion-api` | `https://openrouter.ai/api/v1` | API Key |
| OpenAI | `openai` | `https://api.openai.com/v1` | API Key |
| Anthropic | `anthropic` | `https://api.anthropic.com` | API Key |
| Groq | `generic-chat-completion-api` | `https://api.groq.com/openai/v1` | API Key |
| Google Gemini | `generic-chat-completion-api` | `https://generativelanguage.googleapis.com/v1beta/` | API Key |
| DeepInfra | `generic-chat-completion-api` | `https://api.deepinfra.com/v1/openai` | API Key |
| Fireworks AI | `generic-chat-completion-api` | `https://api.fireworks.ai/inference/v1` | API Key |
| Hugging Face | `generic-chat-completion-api` | `https://router.huggingface.co/v1` | API Key |
| Ollama (Local) | `generic-chat-completion-api` | `http://localhost:11434/v1` | None |
| OpenAI Compatible | `generic-chat-completion-api` | Custom URL | API Key |
| Anthropic Compatible | `anthropic` | Custom URL | API Key |
| Custom Provider | `generic-chat-completion-api` | Custom URL | Optional |

## Provider Types

| Type | Description |
|------|-------------|
| `generic-chat-completion-api` | Any OpenAI-compatible `/v1/chat/completions` endpoint |
| `openai` | Native OpenAI API (Responses API) |
| `anthropic` | Anthropic Messages API format |

## Getting API Keys

### OpenRouter
1. Go to [openrouter.ai/keys](https://openrouter.ai/keys)
2. Create a new API key
3. Paste it into the BYOK wizard when prompted

### OpenAI
1. Go to [platform.openai.com/api-keys](https://platform.openai.com/api-keys)
2. Create a new key
3. Paste it into the BYOK wizard

### Anthropic
1. Go to [console.anthropic.com/settings/keys](https://console.anthropic.com/settings/keys)
2. Create a new key
3. Paste it into the BYOK wizard

### Groq
1. Go to [console.groq.com/keys](https://console.groq.com/keys)
2. Create a new key

### Google Gemini
1. Go to [aistudio.google.com/app/apikey](https://aistudio.google.com/app/apikey)
2. Create a new key

### DeepInfra
1. Go to [deepinfra.com/dash/api_keys](https://deepinfra.com/dash/api_keys)
2. Create a new key

### Fireworks AI
1. Go to [fireworks.ai/api-keys](https://fireworks.ai/api-keys)
2. Create a new key

### Hugging Face
1. Go to [huggingface.co/settings/tokens](https://huggingface.co/settings/tokens)
2. Create a token with **Read** access

### Ollama (Local)
No API key needed. Make sure Ollama is running before launching the wizard:

```bash
ollama serve
```

## Using Environment Variables

Instead of pasting a raw key you can use an environment variable reference:

```
${OPENROUTER_API_KEY}
```

This is stored as-is in `~/.factory/settings.json`. Factory CLI expands it at runtime using your shell environment.

## Custom Providers

Any provider that implements the OpenAI `/v1/chat/completions` API can be added as a custom provider. Examples:

- [LM Studio](https://lmstudio.ai/) — `http://localhost:1234/v1`
- [LocalAI](https://localai.io/) — `http://localhost:8080/v1`
- [vLLM](https://docs.vllm.ai/) — `http://localhost:8000/v1`
- [Text Generation Inference](https://huggingface.co/docs/text-generation-inference) — `http://localhost:3000/v1`
- Any self-hosted proxy or gateway

Select **OpenAI Compatible (Custom URL)** or **Custom Provider** from the provider list, enter your base URL, and continue through the wizard.
