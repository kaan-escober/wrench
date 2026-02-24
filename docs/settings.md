# Settings Reference

All settings are stored in `~/.factory/settings.json`. droid-cfg reads and writes this file while preserving every field it does not manage (hooks, custom droids config, etc.).

---

## Model & Reasoning  `MOD`

### Default Model

The AI model Droid uses for all tasks unless overridden per-session.

| Value | Description |
|-------|-------------|
| `opus` | Claude Opus 4.5 — default |
| `opus-4-6` | Claude Opus 4.6 — max reasoning |
| `opus-4-6-fast` | Opus 4.6 tuned for speed |
| `sonnet` | Claude Sonnet 4.5 — balanced |
| `gpt-5.1` | OpenAI GPT-5.1 |
| `gpt-5.1-codex` | Advanced coding focus |
| `gpt-5.1-codex-max` | Extra high reasoning |
| `gpt-5.2` | OpenAI GPT-5.2 |
| `gpt-5.2-codex` | GPT-5.2 coding |
| `gpt-5.3-codex` | Latest OpenAI coding model |
| `haiku` | Claude Haiku 4.5 — fast & cheap |
| `gemini-3-pro` | Google Gemini 3 Pro |
| `droid-core` | GLM-4.7 open-source |
| `kimi-k2.5` | Kimi K2.5 with image support |
| `minimax-m2.5` | MiniMax M2.5 — 0.12× cost |
| `custom-model` | Your BYOK configured model |

**JSON key:** `model`

### Reasoning Effort

Controls how much structured deliberation the model applies before responding.

| Value | Description |
|-------|-------------|
| `off` / `none` | No reasoning — fastest |
| `low` | Light deliberation |
| `medium` | Balanced thinking — GPT-5 default |
| `high` | Maximum deliberation — slowest |

**JSON key:** `reasoningEffort`

---

## Autonomy  `AUTO`

### Autonomy Level

Controls how proactively Droid executes commands without asking.

| Value | Description |
|-------|-------------|
| `normal` | Ask before every tool use — default |
| `spec` | Plan first, then execute |
| `auto-low` | Auto-approve low-risk actions |
| `auto-medium` | Auto-approve medium-risk actions |
| `auto-high` | Auto-approve most actions |

**JSON key:** `autonomyLevel`

---

## Display  `DISP`

### Diff Mode

How Droid renders file diffs in the TUI.

| Value | Description |
|-------|-------------|
| `github` | Side-by-side GitHub-style — default |
| `unified` | Traditional single-column |

**JSON key:** `diffMode`

### Todo Display

Where Droid's todo list appears.

| Value | Description |
|-------|-------------|
| `pinned` | Pinned above the input area — default |
| `inline` | Inline within the message flow |

**JSON key:** `todoDisplayMode`

### Show AI Thinking

Show the model's internal reasoning steps in the main view.

| Value | Description |
|-------|-------------|
| `off` | Hidden — default |
| `on` | Visible in main view |

**JSON key:** `showThinkingInMainView`

---

## Sound  `SND`

### Completion Sound

Sound played when Droid finishes a task.

| Value | Description |
|-------|-------------|
| `fx-ok01` | Soft success bloop — default |
| `fx-ack01` | Tactile ripple feedback |
| `bell` | System terminal bell |
| `off` | No sound |
| *(custom path)* | Path to your own `.wav`, `.mp3`, or `.ogg` |

**JSON key:** `completionSound`

### Input Awaiting Sound

Sound played when Droid is waiting for your input.

| Value | Description |
|-------|-------------|
| `fx-ok01` | Soft success bloop |
| `fx-ack01` | Tactile ripple — default |
| `bell` | System terminal bell |
| `off` | No sound |
| *(custom path)* | Path to your own audio file |

**JSON key:** `awaitingInputSound`

### Sound Focus Mode

When sounds are played relative to terminal focus.

| Value | Description |
|-------|-------------|
| `always` | Play regardless of focus — default |
| `focused` | Only when terminal is focused |
| `unfocused` | Only when terminal is not focused |

**JSON key:** `soundFocusMode`

---

## Security  `SEC`

### Droid Shield

Secret scanning and git guardrails to prevent accidental credential leaks.

**Default:** `on`  
**JSON key:** `enableDroidShield`

### Co-authored Commits

Appends a `Co-authored-by: Droid` trailer to commits made by Droid.

**Default:** `on`  
**JSON key:** `includeCoAuthoredByDroid`

### Background Processes

Allow Droid to spawn background processes (e.g. dev servers) without asking.

**Default:** `off`  
**JSON key:** `allowBackgroundProcesses`

---

## Agent Behavior  `BEHV`

### Cloud Session Sync

Mirror CLI sessions to the Factory web dashboard.

**Default:** `on`  
**JSON key:** `cloudSessionSync`

### IDE Auto-Connect

Automatically connect to an IDE extension from any terminal.

**Default:** `off`  
**JSON key:** `ideAutoConnect`

### Custom Droids

Enable the Custom Droids feature (persona-based agent configurations).

**Default:** `on`  
**JSON key:** `enableCustomDroids`

### Hooks Disabled

Globally disable all hooks execution.

**Default:** `off`  
**JSON key:** `hooksDisabled`

### Spec Save

Persist spec outputs to disk after each spec-mode run.

**Default:** `off`  
**JSON key:** `specSaveEnabled`

### Spec Save Dir

Directory where spec outputs are saved when Spec Save is enabled.

**Default:** `.factory/docs`  
**JSON key:** `specSaveDir`

### Readiness Report

Enable the `/readiness-report` slash command.

**Default:** `off`  
**JSON key:** `enableReadinessReport`

---

## Command Policies  `CMD`

### Command Allowlist

Commands Droid is always allowed to run, bypassing the autonomy level check.

**JSON key:** `commandAllowlist`  
**Type:** array of strings

```json
"commandAllowlist": ["npm test", "go build ./..."]
```

### Command Denylist

Commands Droid is never allowed to run, regardless of autonomy level.

**JSON key:** `commandDenylist`  
**Type:** array of strings

```json
"commandDenylist": ["rm -rf", "sudo"]
```

Use `Tab` in the Command Policies screen to switch between the allowlist and denylist columns. Press `a` to add a command, `d` to delete the selected one.
