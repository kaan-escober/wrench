# Configuration Files

droid-cfg reads and writes three files. All paths are under your home directory.

## File Locations

| File | Purpose |
|------|---------|
| `~/.factory/settings.json` | Factory CLI settings — the primary config file Droid reads |
| `~/.byok-cli/providers.json` | Saved providers with API keys (managed by droid-cfg) |
| `~/.byok-cli/models.json` | Local record of every custom model you have added |

> **Note:** `~/.factory/settings.json` is also used by the Factory CLI itself. droid-cfg is careful to preserve every field it does not manage, so running droid-cfg will never wipe your hooks, workspace settings, or other Factory configuration.

---

## `~/.factory/settings.json`

The main Factory CLI settings file. droid-cfg reads the full file, makes targeted edits, and writes it back — any fields it does not know about are preserved exactly.

### Full example

```json
{
  "model": "opus",
  "reasoningEffort": "medium",
  "autonomyLevel": "normal",
  "diffMode": "github",
  "todoDisplayMode": "pinned",
  "completionSound": "fx-ok01",
  "awaitingInputSound": "fx-ack01",
  "soundFocusMode": "always",
  "enableDroidShield": true,
  "includeCoAuthoredByDroid": true,
  "allowBackgroundProcesses": false,
  "cloudSessionSync": true,
  "ideAutoConnect": false,
  "enableCustomDroids": true,
  "hooksDisabled": false,
  "specSaveEnabled": false,
  "specSaveDir": ".factory/docs",
  "showThinkingInMainView": false,
  "enableReadinessReport": false,
  "commandAllowlist": [],
  "commandDenylist": [],
  "customModels": [
    {
      "model": "anthropic/claude-3-opus",
      "displayName": "Claude 3 Opus [OpenRouter]",
      "baseUrl": "https://openrouter.ai/api/v1",
      "apiKey": "sk-or-v1-...",
      "provider": "generic-chat-completion-api",
      "maxOutputTokens": 16384,
      "supportsImages": true
    }
  ]
}
```

### `customModels` array — field reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `model` | string | Yes | Model ID as the provider recognises it |
| `displayName` | string | Yes | Human-readable name shown in Droid's model picker |
| `baseUrl` | string | Yes | API endpoint (no trailing slash) |
| `apiKey` | string | No | Raw key or `${ENV_VAR}` reference |
| `provider` | string | Yes | Provider type — see [Providers](./providers.md) |
| `maxOutputTokens` | number | No | Max tokens per response — default 16384 |
| `supportsImages` | boolean | No | `true` if the model accepts image inputs |
| `extraArgs` | object | No | Extra JSON fields passed in the API request body |
| `extraHeaders` | object | No | Extra HTTP headers sent with every request |

---

## `~/.byok-cli/providers.json`

Stores provider configurations saved during BYOK wizard runs. droid-cfg shows these at the top of the provider list for quick reuse.

### Example

```json
[
  {
    "name": "OpenRouter",
    "baseUrl": "https://openrouter.ai/api/v1",
    "providerType": "generic-chat-completion-api",
    "modelsEndpoint": "/models",
    "noAuth": false,
    "apiKey": "sk-or-v1-..."
  },
  {
    "name": "My Local LLM",
    "baseUrl": "http://localhost:8000/v1",
    "providerType": "generic-chat-completion-api",
    "modelsEndpoint": "/models",
    "noAuth": true
  }
]
```

---

## `~/.byok-cli/models.json`

A local index of every custom model that has been added. This is separate from `settings.json` and is used internally by droid-cfg to track what has been added and when.

### Example

```json
[
  {
    "modelId": "anthropic/claude-3-opus",
    "providerName": "OpenRouter",
    "baseUrl": "https://openrouter.ai/api/v1",
    "displayName": "Claude 3 Opus",
    "maxOutputTokens": 16384,
    "supportsImages": true,
    "provider": "generic-chat-completion-api",
    "addedAt": "2026-02-24T10:00:00Z"
  }
]
```

---

## Backup & Restore

### Backup

```bash
cp ~/.factory/settings.json ~/.factory/settings.json.bak
cp ~/.byok-cli/providers.json ~/.byok-cli/providers.json.bak
cp ~/.byok-cli/models.json ~/.byok-cli/models.json.bak
```

### Restore

```bash
cp ~/.factory/settings.json.bak ~/.factory/settings.json
cp ~/.byok-cli/providers.json.bak ~/.byok-cli/providers.json
cp ~/.byok-cli/models.json.bak ~/.byok-cli/models.json
```

---

## Reset

Remove all BYOK CLI data without touching other Factory settings:

```bash
rm -rf ~/.byok-cli
```

Remove all custom models from Factory (edit the file and delete the `customModels` array):

```bash
nano ~/.factory/settings.json
```
