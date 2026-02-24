package config

import (
	"encoding/json"
	"os"
)

// Settings mirrors the typed fields in ~/.factory/settings.json.
// Unknown fields are preserved via the raw map in ReadSettings/WriteSettings.
type Settings struct {
	Model                    string   `json:"model,omitempty"`
	ReasoningEffort          string   `json:"reasoningEffort,omitempty"`
	AutonomyLevel            string   `json:"autonomyLevel,omitempty"`
	CloudSessionSync         *bool    `json:"cloudSessionSync,omitempty"`
	DiffMode                 string   `json:"diffMode,omitempty"`
	CompletionSound          string   `json:"completionSound,omitempty"`
	AwaitingInputSound       string   `json:"awaitingInputSound,omitempty"`
	SoundFocusMode           string   `json:"soundFocusMode,omitempty"`
	CommandAllowlist         []string `json:"commandAllowlist,omitempty"`
	CommandDenylist          []string `json:"commandDenylist,omitempty"`
	IncludeCoAuthoredByDroid *bool    `json:"includeCoAuthoredByDroid,omitempty"`
	EnableDroidShield        *bool    `json:"enableDroidShield,omitempty"`
	HooksDisabled            *bool    `json:"hooksDisabled,omitempty"`
	IdeAutoConnect           *bool    `json:"ideAutoConnect,omitempty"`
	TodoDisplayMode          string   `json:"todoDisplayMode,omitempty"`
	SpecSaveEnabled          *bool    `json:"specSaveEnabled,omitempty"`
	SpecSaveDir              string   `json:"specSaveDir,omitempty"`
	EnableCustomDroids       *bool    `json:"enableCustomDroids,omitempty"`
	ShowThinkingInMainView   *bool    `json:"showThinkingInMainView,omitempty"`
	AllowBackgroundProcesses *bool    `json:"allowBackgroundProcesses,omitempty"`
	EnableReadinessReport    *bool    `json:"enableReadinessReport,omitempty"`
}

// GetField returns the value of a string/enum setting by its JSON key.
func (s *Settings) GetField(key string) string {
	switch key {
	case "model":              return s.Model
	case "reasoningEffort":    return s.ReasoningEffort
	case "autonomyLevel":      return s.AutonomyLevel
	case "diffMode":           return s.DiffMode
	case "completionSound":    return s.CompletionSound
	case "awaitingInputSound": return s.AwaitingInputSound
	case "soundFocusMode":     return s.SoundFocusMode
	case "todoDisplayMode":    return s.TodoDisplayMode
	case "specSaveDir":        return s.SpecSaveDir
	}
	return ""
}

// SetField sets a string/enum setting by its JSON key.
func (s *Settings) SetField(key, val string) {
	switch key {
	case "model":              s.Model = val
	case "reasoningEffort":    s.ReasoningEffort = val
	case "autonomyLevel":      s.AutonomyLevel = val
	case "diffMode":           s.DiffMode = val
	case "completionSound":    s.CompletionSound = val
	case "awaitingInputSound": s.AwaitingInputSound = val
	case "soundFocusMode":     s.SoundFocusMode = val
	case "todoDisplayMode":    s.TodoDisplayMode = val
	case "specSaveDir":        s.SpecSaveDir = val
	}
}

// GetBool returns a *bool setting by its JSON key.
func (s *Settings) GetBool(key string) *bool {
	switch key {
	case "cloudSessionSync":         return s.CloudSessionSync
	case "includeCoAuthoredByDroid": return s.IncludeCoAuthoredByDroid
	case "enableDroidShield":        return s.EnableDroidShield
	case "hooksDisabled":            return s.HooksDisabled
	case "ideAutoConnect":           return s.IdeAutoConnect
	case "specSaveEnabled":          return s.SpecSaveEnabled
	case "enableCustomDroids":       return s.EnableCustomDroids
	case "showThinkingInMainView":   return s.ShowThinkingInMainView
	case "allowBackgroundProcesses": return s.AllowBackgroundProcesses
	case "enableReadinessReport":    return s.EnableReadinessReport
	}
	return nil
}

// SetBool sets a *bool setting by its JSON key.
func (s *Settings) SetBool(key string, val bool) {
	switch key {
	case "cloudSessionSync":         s.CloudSessionSync = bptr(val)
	case "includeCoAuthoredByDroid": s.IncludeCoAuthoredByDroid = bptr(val)
	case "enableDroidShield":        s.EnableDroidShield = bptr(val)
	case "hooksDisabled":            s.HooksDisabled = bptr(val)
	case "ideAutoConnect":           s.IdeAutoConnect = bptr(val)
	case "specSaveEnabled":          s.SpecSaveEnabled = bptr(val)
	case "enableCustomDroids":       s.EnableCustomDroids = bptr(val)
	case "showThinkingInMainView":   s.ShowThinkingInMainView = bptr(val)
	case "allowBackgroundProcesses": s.AllowBackgroundProcesses = bptr(val)
	case "enableReadinessReport":    s.EnableReadinessReport = bptr(val)
	}
}

func bptr(b bool) *bool { return &b }

// ReadSettings loads Settings and the raw map from settings.json.
// The raw map preserves ALL fields (customModels, hooks, etc.) so they are
// never lost when we write back.
func ReadSettings() (Settings, map[string]any, error) {
	raw := map[string]any{}
	data, err := os.ReadFile(settingsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return Settings{}, raw, nil
		}
		return Settings{}, raw, err
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Settings{}, raw, err
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return Settings{}, raw, err
	}
	return s, raw, nil
}

// WriteSettings merges s into raw (preserving unknown fields) and atomically writes to disk.
func WriteSettings(s Settings, raw map[string]any) error {
	if err := ensureDir(settingsDirPath()); err != nil {
		return err
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	var patch map[string]any
	if err := json.Unmarshal(data, &patch); err != nil {
		return err
	}
	for k, v := range patch {
		raw[k] = v
	}
	return writeJSON(settingsPath(), raw)
}

func settingsDirPath() string {
	return settingsDir()
}
