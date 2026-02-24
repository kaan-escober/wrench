package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var mu sync.Mutex

// Paths

func settingsPath() string {
	return filepath.Join(home(), ".factory", "settings.json")
}

func home() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return h
}

func settingsDir() string {
	return filepath.Join(home(), ".factory")
}

// SettingsPath returns the path to Factory's settings.json for display.
func SettingsPath() string {
	return settingsPath()
}

// ───────────────────────────────────────────────
// Custom model config (stored in settings.json → customModels)
// ───────────────────────────────────────────────

type ModelConfig struct {
	ID              string         `json:"id,omitempty"`
	Index           int            `json:"index,omitempty"`
	Model           string         `json:"model"`
	DisplayName     string         `json:"displayName"`
	BaseURL         string         `json:"baseUrl"`
	APIKey          string         `json:"apiKey,omitempty"`
	Provider        string         `json:"provider"`
	MaxOutputTokens int            `json:"maxOutputTokens"`
	SupportsImages  bool           `json:"supportsImages,omitempty"`
	ExtraArgs       map[string]any `json:"extraArgs,omitempty"`
	ExtraHeaders    map[string]any `json:"extraHeaders,omitempty"`
}

// ProviderGroup holds models sharing the same ID prefix.
type ProviderGroup struct {
	Prefix string
	Models []ModelConfig
}

// ReadCustomModels reads the customModels array from settings.json.
func ReadCustomModels() ([]ModelConfig, error) {
	data, err := os.ReadFile(settingsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	arr, ok := raw["customModels"]
	if !ok {
		return nil, nil
	}
	b, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	var models []ModelConfig
	if err := json.Unmarshal(b, &models); err != nil {
		return nil, err
	}
	return models, nil
}

// ReadProviderGroups reads customModels and groups them by ID prefix.
func ReadProviderGroups() ([]ProviderGroup, error) {
	models, err := ReadCustomModels()
	if err != nil {
		return nil, err
	}
	groups := map[string][]ModelConfig{}
	order := []string{}
	for _, m := range models {
		prefix := IDPrefix(m.ID)
		if _, exists := groups[prefix]; !exists {
			order = append(order, prefix)
		}
		groups[prefix] = append(groups[prefix], m)
	}
	out := make([]ProviderGroup, len(order))
	for i, prefix := range order {
		out[i] = ProviderGroup{Prefix: prefix, Models: groups[prefix]}
	}
	return out, nil
}

// GetNextModelIndex returns the next available index for a given prefix.
func GetNextModelIndex(prefix string) (int, error) {
	models, err := ReadCustomModels()
	if err != nil {
		return 0, err
	}
	max := -1
	for _, m := range models {
		if IDPrefix(m.ID) == prefix && m.Index > max {
			max = m.Index
		}
	}
	return max + 1, nil
}

// GenerateModelID creates an ID in the form "prefix:index".
func GenerateModelID(prefix string, index int) string {
	return fmt.Sprintf("%s:%d", prefix, index)
}

// AddModelToSettings adds or updates a model in settings.json.
// If cfg.ID is empty, an ID and index are auto-generated.
func AddModelToSettings(cfg ModelConfig) error {
	mu.Lock()
	defer mu.Unlock()

	path := settingsPath()
	if err := ensureDir(filepath.Dir(path)); err != nil {
		return err
	}

	raw := map[string]any{}
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &raw) //nolint
	}

	var models []any
	if v, ok := raw["customModels"]; ok {
		if arr, ok := v.([]any); ok {
			models = arr
		}
	}

	// Auto-generate ID if not set (caller should set ID with the provider key prefix)
	if cfg.ID == "" {
		prefix := "custom"
		nextIdx := nextIndexFromRaw(models, prefix)
		cfg.Index = nextIdx
		cfg.ID = GenerateModelID(prefix, nextIdx)
	}

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	var cfgMap map[string]any
	json.Unmarshal(cfgBytes, &cfgMap) //nolint

	replaced := false
	for i, m := range models {
		if existing, ok := m.(map[string]any); ok {
			if existing["id"] == cfg.ID {
				models[i] = cfgMap
				replaced = true
				break
			}
		}
	}
	if !replaced {
		models = append(models, cfgMap)
	}

	raw["customModels"] = models
	return writeJSON(path, raw)
}

// DeleteModelFromSettings removes a model by ID from settings.json.
func DeleteModelFromSettings(id string) error {
	mu.Lock()
	defer mu.Unlock()

	path := settingsPath()
	raw := map[string]any{}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	arr, ok := raw["customModels"]
	if !ok {
		return nil
	}
	models, ok := arr.([]any)
	if !ok {
		return nil
	}

	filtered := make([]any, 0, len(models))
	for _, m := range models {
		if existing, ok := m.(map[string]any); ok {
			if existing["id"] == id {
				continue
			}
		}
		filtered = append(filtered, m)
	}

	raw["customModels"] = filtered
	return writeJSON(path, raw)
}

// ───────────────────────────────────────────────
// Helpers
// ───────────────────────────────────────────────

// IDPrefix extracts the prefix from an ID like "openrouter:0" → "openrouter".
func IDPrefix(id string) string {
	if i := strings.LastIndex(id, ":"); i >= 0 {
		return id[:i]
	}
	if id == "" {
		return "custom"
	}
	return id
}

func nextIndexFromRaw(models []any, prefix string) int {
	max := -1
	for _, m := range models {
		if existing, ok := m.(map[string]any); ok {
			if eid, ok := existing["id"].(string); ok && IDPrefix(eid) == prefix {
				if idx, ok := existing["index"].(float64); ok && int(idx) > max {
					max = int(idx)
				}
			}
		}
	}
	return max + 1
}

// GroupDisplayName returns the prefix as the display name for a provider group.
func GroupDisplayName(g ProviderGroup) string {
	return g.Prefix
}

// SortProviderGroups sorts groups alphabetically by prefix.
func SortProviderGroups(groups []ProviderGroup) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0o700)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".wrench-tmp-*")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	if err := os.Chmod(tmpName, 0o600); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}
