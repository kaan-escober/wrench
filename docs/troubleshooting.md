# Troubleshooting

## The binary won't run on Termux

**Symptom:** `exec format error` or `no such file or directory`

Build the binary directly on your device — do not copy a binary compiled on another architecture:

```bash
pkg install golang git
cd byok-cli/droid-cfg
go build -o droid-cfg .
```

---

## The TUI looks garbled or has missing characters

**Symptom:** Unicode box-drawing characters, circles, or arrows render as `?` or blank squares.

1. Make sure your Termux font supports Unicode. Install a Nerd Font or use the default Termux font which already includes most block characters.
2. Set your locale:
   ```bash
   export LANG=en_US.UTF-8
   export LC_ALL=en_US.UTF-8
   ```
3. If you are using tmux, make sure it is configured for UTF-8:
   ```
   set -g utf8 on
   set-window-option -g utf8 on
   ```

---

## Settings are not saving

**Symptom:** Changes made in droid-cfg do not appear in Factory CLI, or the values reset on the next open.

1. Check that `~/.factory/settings.json` exists and is writable:
   ```bash
   ls -la ~/.factory/settings.json
   ```
2. If the file does not exist, create the directory:
   ```bash
   mkdir -p ~/.factory
   ```
3. Make sure the JSON file is valid — an invalid file will be silently treated as empty:
   ```bash
   cat ~/.factory/settings.json | python3 -m json.tool
   ```

---

## BYOK: model fetch fails

**Symptom:** "Could not auto-fetch" screen appears after entering provider details.

Possible causes:

- The provider's `/models` endpoint requires a different path. Some providers use `/api/models` or do not expose a list endpoint at all. Enter the model ID manually when prompted.
- The API key is wrong or expired. Double-check on the provider's dashboard.
- The base URL is incorrect — make sure it ends with `/v1` for OpenAI-compatible providers.
- Network issue on Termux — check that your device has internet access:
  ```bash
  curl -s https://api.openai.com/v1/models -H "Authorization: Bearer $OPENAI_API_KEY"
  ```

---

## BYOK: saved providers not appearing

**Symptom:** The provider list shows no saved entries even though you have added providers before.

Check that `~/.byok-cli/providers.json` exists and contains valid JSON:

```bash
cat ~/.byok-cli/providers.json | python3 -m json.tool
```

If the file is corrupted, remove it and re-add your providers:

```bash
rm ~/.byok-cli/providers.json
```

---

## Footer line appears twice

**Symptom:** The `─` separator line appears on two rows at the bottom of the screen.

This is a known rendering glitch that occurs when the terminal width reported by the OS is different from the actual display width (common in Termux when the keyboard slides up). Resize the terminal or toggle the soft keyboard — the layout will re-render correctly on the next `WindowSizeMsg`.

---

## The cursor row badge text wraps to the next line

**Symptom:** Badges like `BEHV` appear split as `BE` / `HV` on separate rows.

This happens when the terminal font reports a different character width than expected. Try a monospace font with standard glyph widths. In Termux, go to **Settings → Styling → Font** and select the default monospace option.

---

## Custom model does not appear in Droid's model picker

After adding a model via the BYOK wizard:

1. Restart Factory CLI — it reads `settings.json` at startup.
2. Verify the entry was written:
   ```bash
   cat ~/.factory/settings.json | python3 -m json.tool | grep -A 10 '"customModels"'
   ```
3. Make sure the `model` field matches the ID the provider expects, and the `baseUrl` is correct.

---

## droid-cfg panics on launch

Please open an issue at [github.com/kartvya69/byok-cli](https://github.com/kartvya69/byok-cli/issues) and include:

- Go version: `go version`
- OS / architecture: `uname -a`
- The full panic output
