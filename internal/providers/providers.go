package providers

// Provider holds static config for a known provider.
type Provider struct {
	Name            string
	BaseURL         string
	Type            string // generic-chat-completion-api | anthropic | openai
	ModelsEndpoint  string // "" means no auto-fetch
	RequiresBaseURL bool
	NoAuth          bool
}

// All known providers, in display order.
var All = []struct {
	Key      string
	Provider Provider
}{
	{"openrouter", Provider{
		Name:           "OpenRouter",
		BaseURL:        "https://openrouter.ai/api/v1",
		Type:           "generic-chat-completion-api",
		ModelsEndpoint: "/models",
	}},
	{"openai", Provider{
		Name:           "OpenAI",
		BaseURL:        "https://api.openai.com/v1",
		Type:           "openai",
		ModelsEndpoint: "/models",
	}},
	{"anthropic", Provider{
		Name:    "Anthropic",
		BaseURL: "https://api.anthropic.com",
		Type:    "anthropic",
	}},
	{"groq", Provider{
		Name:           "Groq",
		BaseURL:        "https://api.groq.com/openai/v1",
		Type:           "generic-chat-completion-api",
		ModelsEndpoint: "/models",
	}},
	{"gemini", Provider{
		Name:    "Google Gemini",
		BaseURL: "https://generativelanguage.googleapis.com/v1beta/",
		Type:    "generic-chat-completion-api",
	}},
	{"deepinfra", Provider{
		Name:           "DeepInfra",
		BaseURL:        "https://api.deepinfra.com/v1/openai",
		Type:           "generic-chat-completion-api",
		ModelsEndpoint: "/models",
	}},
	{"fireworks", Provider{
		Name:           "Fireworks AI",
		BaseURL:        "https://api.fireworks.ai/inference/v1",
		Type:           "generic-chat-completion-api",
		ModelsEndpoint: "/models",
	}},
	{"huggingface", Provider{
		Name:           "Hugging Face",
		BaseURL:        "https://router.huggingface.co/v1",
		Type:           "generic-chat-completion-api",
		ModelsEndpoint: "/models",
	}},
	{"ollama", Provider{
		Name:           "Ollama (Local)",
		BaseURL:        "http://localhost:11434/v1",
		Type:           "generic-chat-completion-api",
		ModelsEndpoint: "/models",
		NoAuth:         true,
	}},
	{"openai-compatible", Provider{
		Name:            "OpenAI Compatible (Custom URL)",
		Type:            "generic-chat-completion-api",
		ModelsEndpoint:  "/models",
		RequiresBaseURL: true,
	}},
	{"anthropic-compatible", Provider{
		Name:            "Anthropic Compatible (Custom URL)",
		Type:            "anthropic",
		RequiresBaseURL: true,
	}},
	{"custom", Provider{
		Name:            "Custom",
		Type:            "generic-chat-completion-api",
		ModelsEndpoint:  "/models",
		RequiresBaseURL: true,
	}},
}

// ProviderTypes for display in settings
var ProviderTypes = []struct {
	Value string
	Label string
}{
	{"generic-chat-completion-api", "OpenAI-compatible (Chat Completions)"},
	{"openai", "OpenAI (Responses API)"},
	{"anthropic", "Anthropic (Messages API)"},
}

// Get returns a provider by key, nil if not found.
func Get(key string) *Provider {
	for _, p := range All {
		if p.Key == key {
			cp := p.Provider
			return &cp
		}
	}
	return nil
}
