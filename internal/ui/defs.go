package ui

// AppMode is the top-level navigation state of the application.
type AppMode int

const (
	ModeMenu        AppMode = iota // main dashboard
	ModeCategory                   // inside a settings category list
	ModeOptionPick                 // picking an enum value from a list
	ModeBoolPick                   // picking yes/no for a boolean
	ModeTextInput                  // free-text entry for a setting
	ModeCommandEdit                // command allow/deny list editor
	ModeCommandAdd                 // text input inside command editor
	ModeBYOK                       // full BYOK wizard
)

// Category identifies a settings group.
type Category int

const (
	CatBYOK     Category = iota
	CatModel
	CatAutonomy
	CatDisplay
	CatSound
	CatSecurity
	CatBehavior
	CatCommands
)

// SettingKind is the input type for a setting.
type SettingKind int

const (
	KindEnum SettingKind = iota
	KindBool
	KindText
)

// OptionDef is one selectable value in an enum picker.
type OptionDef struct {
	Value string
	Label string
	Desc  string
}

// SettingDef describes one configurable setting within a category.
type SettingDef struct {
	Label   string
	Key     string
	Kind    SettingKind
	Options []OptionDef
	Default string // displayed when value is unset
}

// ─────────────────────────────────────────────────────────────────────────────
// All category definitions
// ─────────────────────────────────────────────────────────────────────────────

var categorySettings = map[Category][]SettingDef{
	CatModel: {
		{
			Label: "Default Model", Key: "model", Kind: KindEnum, Default: "opus",
			Options: []OptionDef{
				{Value: "opus", Label: "opus", Desc: "Claude Opus 4.5  (default)"},
				{Value: "opus-4-6", Label: "opus-4-6", Desc: "Claude Opus 4.6 · Max reasoning"},
				{Value: "opus-4-6-fast", Label: "opus-4-6-fast", Desc: "Opus 4.6 Fast · tuned for speed"},
				{Value: "sonnet", Label: "sonnet", Desc: "Claude Sonnet 4.5 · balanced"},
				{Value: "gpt-5.1", Label: "gpt-5.1", Desc: "OpenAI GPT-5.1"},
				{Value: "gpt-5.1-codex", Label: "gpt-5.1-codex", Desc: "Advanced coding focus"},
				{Value: "gpt-5.1-codex-max", Label: "gpt-5.1-codex-max", Desc: "Extra High reasoning"},
				{Value: "gpt-5.2", Label: "gpt-5.2", Desc: "OpenAI GPT-5.2"},
				{Value: "gpt-5.2-codex", Label: "gpt-5.2-codex", Desc: "GPT-5.2 coding · Extra High"},
				{Value: "gpt-5.3-codex", Label: "gpt-5.3-codex", Desc: "Latest OpenAI coding model"},
				{Value: "haiku", Label: "haiku", Desc: "Claude Haiku 4.5 · fast & cheap"},
				{Value: "gemini-3-pro", Label: "gemini-3-pro", Desc: "Google Gemini 3 Pro"},
				{Value: "droid-core", Label: "droid-core", Desc: "GLM-4.7 open-source"},
				{Value: "kimi-k2.5", Label: "kimi-k2.5", Desc: "Kimi K2.5 · image support"},
				{Value: "minimax-m2.5", Label: "minimax-m2.5", Desc: "MiniMax M2.5 · 0.12× cost"},
				{Value: "custom-model", Label: "custom-model", Desc: "Your BYOK configured model"},
			},
		},
		{
			Label: "Reasoning Effort", Key: "reasoningEffort", Kind: KindEnum, Default: "model default",
			Options: []OptionDef{
				{Value: "off", Label: "off", Desc: "No structured reasoning · fastest"},
				{Value: "none", Label: "none", Desc: "Alias for off"},
				{Value: "low", Label: "low", Desc: "Light deliberation"},
				{Value: "medium", Label: "medium", Desc: "Balanced thinking · GPT-5 default"},
				{Value: "high", Label: "high", Desc: "Maximum deliberation · slowest"},
			},
		},
	},

	CatAutonomy: {
		{
			Label: "Autonomy Level", Key: "autonomyLevel", Kind: KindEnum, Default: "normal",
			Options: []OptionDef{
				{Value: "normal", Label: "normal", Desc: "Ask before every tool use  (default)"},
				{Value: "spec", Label: "spec", Desc: "Plan first, then execute"},
				{Value: "auto-low", Label: "auto-low", Desc: "Auto-approve low-risk actions"},
				{Value: "auto-medium", Label: "auto-medium", Desc: "Auto-approve medium-risk actions"},
				{Value: "auto-high", Label: "auto-high", Desc: "Auto-approve most actions"},
			},
		},
	},

	CatDisplay: {
		{
			Label: "Diff Mode", Key: "diffMode", Kind: KindEnum, Default: "github",
			Options: []OptionDef{
				{Value: "github", Label: "github", Desc: "Side-by-side GitHub-style  (default)"},
				{Value: "unified", Label: "unified", Desc: "Traditional single-column"},
			},
		},
		{
			Label: "Todo Display", Key: "todoDisplayMode", Kind: KindEnum, Default: "pinned",
			Options: []OptionDef{
				{Value: "pinned", Label: "pinned", Desc: "Pinned above input area  (default)"},
				{Value: "inline", Label: "inline", Desc: "Inline within message flow"},
			},
		},
		{
			Label: "Show AI Thinking", Key: "showThinkingInMainView", Kind: KindBool, Default: "off",
		},
	},

	CatSound: {
		{
			Label: "Completion Sound", Key: "completionSound", Kind: KindEnum, Default: "fx-ok01",
			Options: []OptionDef{
				{Value: "fx-ok01", Label: "fx-ok01", Desc: "Soft success bloop  (default)"},
				{Value: "fx-ack01", Label: "fx-ack01", Desc: "Tactile ripple feedback"},
				{Value: "bell", Label: "bell", Desc: "System terminal bell"},
				{Value: "off", Label: "off", Desc: "No sound"},
				{Value: "__custom__", Label: "Custom file path...", Desc: "Provide your own .wav / .mp3"},
			},
		},
		{
			Label: "Input Awaiting Sound", Key: "awaitingInputSound", Kind: KindEnum, Default: "fx-ack01",
			Options: []OptionDef{
				{Value: "fx-ok01", Label: "fx-ok01", Desc: "Soft success bloop"},
				{Value: "fx-ack01", Label: "fx-ack01", Desc: "Tactile ripple  (default)"},
				{Value: "bell", Label: "bell", Desc: "System terminal bell"},
				{Value: "off", Label: "off", Desc: "No sound"},
				{Value: "__custom__", Label: "Custom file path...", Desc: "Provide your own .wav / .mp3"},
			},
		},
		{
			Label: "Sound Focus Mode", Key: "soundFocusMode", Kind: KindEnum, Default: "always",
			Options: []OptionDef{
				{Value: "always", Label: "always", Desc: "Play regardless of focus  (default)"},
				{Value: "focused", Label: "focused", Desc: "Only when terminal is focused"},
				{Value: "unfocused", Label: "unfocused", Desc: "Only when terminal is not focused"},
			},
		},
	},

	CatSecurity: {
		{Label: "Droid Shield", Key: "enableDroidShield", Kind: KindBool, Default: "on",
			Options: []OptionDef{{Value: "", Desc: "Secret scanning & git guardrails"}},
		},
		{Label: "Co-authored Commits", Key: "includeCoAuthoredByDroid", Kind: KindBool, Default: "on",
			Options: []OptionDef{{Value: "", Desc: "Append co-author trailer to commits"}},
		},
		{Label: "Background Processes", Key: "allowBackgroundProcesses", Kind: KindBool, Default: "off",
			Options: []OptionDef{{Value: "", Desc: "Allow Droid to spawn background procs"}},
		},
	},

	CatBehavior: {
		{Label: "Cloud Session Sync", Key: "cloudSessionSync", Kind: KindBool, Default: "on",
			Options: []OptionDef{{Value: "", Desc: "Mirror CLI sessions to Factory web"}},
		},
		{Label: "IDE Auto-Connect", Key: "ideAutoConnect", Kind: KindBool, Default: "off",
			Options: []OptionDef{{Value: "", Desc: "Auto-connect to IDE from any terminal"}},
		},
		{Label: "Custom Droids", Key: "enableCustomDroids", Kind: KindBool, Default: "on",
			Options: []OptionDef{{Value: "", Desc: "Enable the Custom Droids feature"}},
		},
		{Label: "Hooks Disabled", Key: "hooksDisabled", Kind: KindBool, Default: "off",
			Options: []OptionDef{{Value: "", Desc: "Globally disable all hooks execution"}},
		},
		{Label: "Spec Save", Key: "specSaveEnabled", Kind: KindBool, Default: "off",
			Options: []OptionDef{{Value: "", Desc: "Persist spec outputs to disk"}},
		},
		{Label: "Spec Save Dir", Key: "specSaveDir", Kind: KindText, Default: ".factory/docs"},
		{Label: "Readiness Report", Key: "enableReadinessReport", Kind: KindBool, Default: "off",
			Options: []OptionDef{{Value: "", Desc: "Enable /readiness-report command"}},
		},
	},
}

// menuEntry is one row in the main dashboard menu.
type menuEntry struct {
	cat   Category
	badge string
	label string
}

var menuEntries = []menuEntry{
	{CatBYOK, "BYOK", "Custom Models"},
	{CatModel, "MOD", "Model & Reasoning"},
	{CatAutonomy, "AUTO", "Autonomy"},
	{CatDisplay, "DISP", "Display"},
	{CatSound, "SND", "Sound"},
	{CatSecurity, "SEC", "Security"},
	{CatBehavior, "BEHV", "Agent Behavior"},
	{CatCommands, "CMD", "Command Policies"},
}
