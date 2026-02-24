package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var client = &http.Client{Timeout: 15 * time.Second}

// ModelInfo is a model returned from a provider's /models endpoint.
type ModelInfo struct {
	ID   string
	Name string
}

// FetchModels calls the provider's models endpoint and returns available models.
func FetchModels(baseURL, apiKey, modelsEndpoint, providerType string, noAuth bool) ([]ModelInfo, error) {
	if modelsEndpoint == "" {
		return nil, nil
	}

	endpoint := strings.TrimRight(baseURL, "/") + modelsEndpoint
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if !noAuth && apiKey != "" {
		if providerType == "anthropic" {
			req.Header.Set("x-api-key", apiKey)
			req.Header.Set("anthropic-dangerous-direct-browser-access", "true")
		} else {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, endpoint)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseModelsResponse(body)
}

func parseModelsResponse(body []byte) ([]ModelInfo, error) {
	// Try { "data": [ { "id": "..." } ] } — OpenAI format
	var openaiResp struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &openaiResp); err == nil && len(openaiResp.Data) > 0 {
		out := make([]ModelInfo, 0, len(openaiResp.Data))
		for _, m := range openaiResp.Data {
			if m.ID == "" {
				continue
			}
			name := m.Name
			if name == "" {
				name = m.ID
			}
			out = append(out, ModelInfo{ID: m.ID, Name: name})
		}
		return out, nil
	}

	// Try { "models": [ ... ] }
	var modelsResp struct {
		Models []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.Unmarshal(body, &modelsResp); err == nil && len(modelsResp.Models) > 0 {
		out := make([]ModelInfo, 0, len(modelsResp.Models))
		for _, m := range modelsResp.Models {
			if m.ID == "" {
				continue
			}
			name := m.Name
			if name == "" {
				name = m.ID
			}
			out = append(out, ModelInfo{ID: m.ID, Name: name})
		}
		return out, nil
	}

	// Try direct array [ "model-id" ] or [ { "id": "..." } ]
	var rawArr []json.RawMessage
	if err := json.Unmarshal(body, &rawArr); err == nil {
		out := make([]ModelInfo, 0, len(rawArr))
		for _, raw := range rawArr {
			var s string
			if json.Unmarshal(raw, &s) == nil {
				out = append(out, ModelInfo{ID: s, Name: s})
				continue
			}
			var obj struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			if json.Unmarshal(raw, &obj) == nil && obj.ID != "" {
				name := obj.Name
				if name == "" {
					name = obj.ID
				}
				out = append(out, ModelInfo{ID: obj.ID, Name: name})
			}
		}
		return out, nil
	}

	return nil, fmt.Errorf("unrecognised models response format")
}

// ───────────────────────────────────────────────
// models.dev enrichment
// ───────────────────────────────────────────────

type modelsDevEntry struct {
	ID   string
	Name string
}

var (
	modelsDevCache map[string]modelsDevEntry
	modelsDevOnce  sync.Once
)

func fetchModelsDevData() map[string]modelsDevEntry {
	modelsDevOnce.Do(func() {
		resp, err := client.Get("https://models.dev/api.json")
		if err != nil {
			modelsDevCache = map[string]modelsDevEntry{}
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			modelsDevCache = map[string]modelsDevEntry{}
			return
		}

		// Structure: { providerId: { models: { modelId: { name: "..." } } } }
		var raw map[string]struct {
			Models map[string]struct {
				Name string `json:"name"`
			} `json:"models"`
		}
		if err := json.Unmarshal(body, &raw); err != nil {
			modelsDevCache = map[string]modelsDevEntry{}
			return
		}

		cache := make(map[string]modelsDevEntry)
		for _, provider := range raw {
			for id, m := range provider.Models {
				name := m.Name
				if name == "" {
					name = id
				}
				cache[id] = modelsDevEntry{ID: id, Name: name}
			}
		}
		modelsDevCache = cache
	})
	return modelsDevCache
}

// GetDisplayName returns a human-friendly model name, falling back to normalising the ID.
func GetDisplayName(modelID string) string {
	cache := fetchModelsDevData()
	if entry, ok := cache[modelID]; ok && entry.Name != "" {
		return entry.Name
	}
	return normalizeID(modelID)
}

func normalizeID(id string) string {
	// Strip provider prefix (e.g. "openai/gpt-4" → "gpt-4")
	if i := strings.Index(id, "/"); i != -1 {
		id = id[i+1:]
	}
	id = strings.ReplaceAll(id, "-", " ")
	id = strings.ReplaceAll(id, "_", " ")
	words := strings.Fields(id)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, " ")
}

// ───────────────────────────────────────────────
// URL helpers
// ───────────────────────────────────────────────

// NormalizeURL ensures the URL has a scheme and no trailing slash mess.
func NormalizeURL(raw string) (string, error) {
	if raw == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}
	if u.Host == "" {
		return "", fmt.Errorf("URL must have a host")
	}
	return u.String(), nil
}

// ExtractProviderName pulls a readable name from a URL (e.g. "api.openrouter.ai" → "openrouter").
func ExtractProviderName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	host := u.Hostname()
	parts := strings.Split(host, ".")
	// Remove "api", "www", etc.
	meaningful := []string{}
	skip := map[string]bool{"api": true, "www": true, "app": true, "v1": true, "v2": true}
	for _, p := range parts {
		if !skip[p] && p != "" {
			meaningful = append(meaningful, p)
		}
	}
	if len(meaningful) > 0 {
		name := meaningful[0]
		return strings.ToUpper(name[:1]) + name[1:]
	}
	return host
}
