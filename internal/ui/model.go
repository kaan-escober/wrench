package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/kaan-escober/wrench/internal/api"
	"github.com/kaan-escober/wrench/internal/config"
	"github.com/kaan-escober/wrench/internal/theme"
)

// ─────────────────────────────────────────────────────────────────────────────
// BYOK wizard step enum (used when mode == ModeBYOK)
// ─────────────────────────────────────────────────────────────────────────────

type WizStep int

const (
	WizProvider       WizStep = iota
	WizGroupDetail           // view models in a provider group
	WizModelEdit             // edit individual model fields
	WizModelField            // editing a specific model field
	WizURL
	WizTitle
	WizKey
	WizFetching
	WizModels
	WizSettingsTokens
	WizSettingsImages
	WizConfirm
	WizSaving
	WizDone
)

// ─────────────────────────────────────────────────────────────────────────────
// Custom list (shared navigation + multi-select component)
// ─────────────────────────────────────────────────────────────────────────────

type listItem struct {
	label string
	value string
	sub   string // optional muted detail on the right
}

type customList struct {
	items    []listItem
	cursor   int
	selected map[int]bool
	multi    bool
	height   int
	offset   int
}

func newList(items []listItem, multi bool, height int) customList {
	return customList{items: items, selected: make(map[int]bool), multi: multi, height: height}
}

func (l *customList) up() {
	if l.cursor > 0 {
		l.cursor--
		if l.cursor < l.offset {
			l.offset = l.cursor
		}
	}
}

func (l *customList) down() {
	if l.cursor < len(l.items)-1 {
		l.cursor++
		if l.cursor >= l.offset+l.height {
			l.offset = l.cursor - l.height + 1
		}
	}
}

func (l *customList) toggleCurrent() {
	if l.multi {
		l.selected[l.cursor] = !l.selected[l.cursor]
	}
}

func (l *customList) selectedValues() []string {
	out := []string{}
	for i, item := range l.items {
		if l.selected[i] {
			out = append(out, item.value)
		}
	}
	return out
}

// render draws the list into a string.
func (l *customList) render(cursorActive bool) string {
	if len(l.items) == 0 {
		return theme.Muted.Render("  (empty)")
	}
	out := ""
	if l.offset > 0 {
		out += theme.Muted.Render("  ↑ more") + "\n"
	}
	end := l.offset + l.height
	if end > len(l.items) {
		end = len(l.items)
	}
	for i := l.offset; i < end; i++ {
		item := l.items[i]
		isCursor := i == l.cursor

		var prefix, label, sub string

		if l.multi {
			if l.selected[i] {
				prefix = theme.Accent.Render("● ")
			} else {
				prefix = theme.Muted.Render("○ ")
			}
		} else {
			prefix = "  "
		}

		label = theme.Primary.Render(item.label)
		if item.sub != "" {
			sub = "  " + theme.Muted.Render(item.sub)
		}

		var line string
		if isCursor && cursorActive {
			arrow := theme.Accent.Render("> ")
			if l.multi {
				line = arrow + prefix + label + sub
			} else {
				line = arrow + label + sub
			}
		} else if isCursor && !cursorActive {
			arrow := theme.Muted.Render("> ")
			if l.multi {
				line = arrow + prefix + label + sub
			} else {
				line = arrow + label + sub
			}
		} else {
			if l.multi {
				line = "  " + prefix + label + sub
			} else {
				line = "  " + label + sub
			}
		}
		out += line + "\n"
	}
	if end < len(l.items) {
		out += theme.Muted.Render("  ↓ more")
	}
	return out
}

// ─────────────────────────────────────────────────────────────────────────────
// Async message types
// ─────────────────────────────────────────────────────────────────────────────

type groupsLoadedMsg struct{ groups []config.ProviderGroup }
type settingsLoadedMsg struct {
	settings config.Settings
	raw      map[string]any
}
type modelsLoadedMsg struct {
	models       []api.ModelInfo
	displayNames map[string]string
}
type byokSavedMsg struct{ path string }
type settingsSavedMsg struct{}
type clearFlashMsg struct{}
type errMsg struct{ err error }

// ─────────────────────────────────────────────────────────────────────────────
// Top-level model
// ─────────────────────────────────────────────────────────────────────────────

type Model struct {
	width, height int

	// ── App navigation ───────────────────────────────────────────────────────
	mode AppMode

	// ── Loaded settings ──────────────────────────────────────────────────────
	settings config.Settings
	rawCfg   map[string]any // preserves unknown fields

	// ── Main menu ────────────────────────────────────────────────────────────
	menuCursor int

	// ── Category view ────────────────────────────────────────────────────────
	currentCat Category
	catCursor  int

	// ── Option / bool picker ─────────────────────────────────────────────────
	optionList  customList
	customInput bool

	// ── Text input (settings + BYOK shared) ──────────────────────────────────
	textInput textinput.Model

	// ── Command policy editor ─────────────────────────────────────────────────
	allowCmds   []string
	denyCmds    []string
	cmdFocusCol int
	cmdCursor   int
	cmdInput    textinput.Model

	// ── BYOK wizard ──────────────────────────────────────────────────────────
	byokStep        WizStep
	providerList    customList
	providerGroups  []config.ProviderGroup
	currentGroupIdx int
	detailList      customList
	providerKey     string
	providerName    string
	baseURL         string
	providerType    string
	apiKey          string
	extractedName   string
	displayTitle    string
	noAuth          bool
	modelsEndpoint  string

	availableModels   []api.ModelInfo
	modelList         customList
	modelDisplayNames map[string]string
	selectedModels    []string
	maxOutputTokens   int
	supportsImages    bool
	savedPath         string

	// ── Model editor ─────────────────────────────────────────────────────────
	editingModel config.ModelConfig
	editFieldKey string

	// ── Spinner ───────────────────────────────────────────────────────────────
	spinner spinner.Model

	// ── Feedback ─────────────────────────────────────────────────────────────
	err   string
	flash string
}

func initialModel() Model {
	ti := textinput.New()
	ti.PromptStyle = theme.Accent
	ti.TextStyle = theme.Primary
	ti.Prompt = theme.Prompt
	ti.Width = 50

	ci := textinput.New()
	ci.PromptStyle = theme.Accent
	ci.TextStyle = theme.Primary
	ci.Prompt = theme.Prompt
	ci.Width = 40
	ci.Placeholder = "command or pattern"

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = theme.Accent

	return Model{
		mode:            ModeMenu,
		rawCfg:          map[string]any{},
		maxOutputTokens: 16384,
		textInput:       ti,
		cmdInput:        ci,
		spinner:         sp,
	}
}

// currentSettingDef returns the SettingDef at catCursor in currentCat.
func (m Model) currentSettingDef() SettingDef {
	defs := categorySettings[m.currentCat]
	if m.catCursor >= 0 && m.catCursor < len(defs) {
		return defs[m.catCursor]
	}
	return SettingDef{}
}
