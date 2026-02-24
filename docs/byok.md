# BYOK Wizard

BYOK (Bring Your Own Key) lets you connect any external AI model to Factory/Droid. Once added, the model appears in Droid's model selector just like a built-in model.

## Starting the Wizard

From the main menu, select **Custom Models** (`BYOK`) and press `Enter`.

## Wizard Steps

### 1. Select Provider

Choose from saved providers (shown at the top if you have any) or pick a built-in one from the list.

- Saved providers are listed first for quick access
- Scroll with `↑↓`, confirm with `Enter`

If you select a saved provider, you skip directly to step 5 (Add Models or Edit).

### 2. Base URL *(custom providers only)*

Enter the API base URL for your provider, e.g.:

```
https://my-llm-api.example.com/v1
```

For OpenAI-compatible endpoints the URL must include the `/v1` path segment.

### 3. Display Name *(custom providers only)*

The name shown next to your models in Droid's model selector. droid-cfg auto-suggests a name extracted from the URL — press `Enter` to accept it, or type a new one.

### 4. API Key

Enter your API key. Supports:

- Raw keys: `sk-abc123...`
- Environment variable references: `${OPENROUTER_API_KEY}`

For providers that require no auth (e.g. Ollama), this step is skipped automatically.

> Keys are stored in `~/.byok-cli/providers.json` and written to `~/.factory/settings.json` inside the model config.

### 5. Fetch Models

droid-cfg calls the provider's `/models` endpoint and shows the full list. If the provider does not expose a models endpoint, you will be prompted to enter a model ID manually.

### 6. Select Models

Browse the fetched model list and press `Space` to toggle each one. You can select as many as you like. Press `Enter` when done.

```
  ○  gpt-4o
> ●  claude-opus-4-5       ← selected
  ●  qwen3:32b             ← selected
  ○  gemini-2.5-pro
```

### 7. Max Output Tokens

Set the maximum tokens the model can return per response. Common values:

| Range | Typical use |
|-------|-------------|
| 8192 | Standard models |
| 16384 | Default suggestion |
| 32768–131072 | Long-context models |

Press `Enter` to confirm.

### 8. Image Support

Select `Yes` if the model can process image inputs (vision models), `No` otherwise.

### 9. Confirm & Save

Review the summary before saving:

```
  Provider:      OpenRouter
  Base URL:      https://openrouter.ai/api/v1
  Type:          generic-chat-completion-api
  Models:        claude-opus-4-5, qwen3:32b
  Max tokens:    16384
  Images:        No
```

Select **Save** to write to `~/.factory/settings.json`, or **Cancel** to discard.

### 10. Done

A success screen confirms how many models were added and where the config was saved. From here you can:

- **Add more models** to the same provider
- **Add models for another provider** — starts the wizard fresh
- **Back to main menu**

## Managing Saved Providers

Providers are saved automatically after the first successful run. On the next run they appear at the top of the provider list.

To update a saved provider's URL or key, select it from the list and choose **Edit configuration**.

## Removing a Custom Model

droid-cfg does not have a delete UI yet. To remove a custom model, edit `~/.factory/settings.json` directly and delete the entry from the `customModels` array:

```bash
nano ~/.factory/settings.json
```
