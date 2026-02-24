package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kaan-escober/wrench/internal/api"
	"github.com/kaan-escober/wrench/internal/config"
	"github.com/kaan-escober/wrench/internal/providers"
)

// ─────────────────────────────────────────────────────────────────────────────
// Init
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadAllSettings(),
		m.spinner.Tick,
	)
}

// ─────────────────────────────────────────────────────────────────────────────
// Update
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.providerList.height = listHeight(m.height)
		m.modelList.height = listHeight(m.height)
		m.optionList.height = listHeight(m.height)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m.handleKey(msg)

	case settingsLoadedMsg:
		m.settings = msg.settings
		m.rawCfg = msg.raw
		return m, nil

	case groupsLoadedMsg:
		m.providerGroups = msg.groups
		m.providerList = buildProviderList(msg.groups)
		m.providerList.height = listHeight(m.height)
		return m, nil

	case modelDeletedMsg:
		m.byokStep = WizProvider
		m.flash = "  ✓ Model deleted"
		return m, tea.Batch(loadProviderGroups(), clearFlashAfter())

	case modelsLoadedMsg:
		m.availableModels = msg.models
		m.modelDisplayNames = msg.displayNames
		m.modelList = buildModelList(msg.models)
		m.modelList.height = listHeight(m.height)
		m.byokStep = WizModels
		return m, nil

	case byokSavedMsg:
		m.savedPath = msg.path
		m.byokStep = WizDone
		m.detailList = buildDoneList()
		return m, loadAllSettings()

	case settingsSavedMsg:
		m.flash = "  ✓ Saved"
		return m, clearFlashAfter()

	case clearFlashMsg:
		m.flash = ""
		return m, nil

	case errMsg:
		m.err = msg.err.Error()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmds []tea.Cmd
	if m.textInput.Focused() {
		var c tea.Cmd
		m.textInput, c = m.textInput.Update(msg)
		cmds = append(cmds, c)
	}
	if m.cmdInput.Focused() {
		var c tea.Cmd
		m.cmdInput, c = m.cmdInput.Update(msg)
		cmds = append(cmds, c)
	}
	return m, tea.Batch(cmds...)
}

// ─────────────────────────────────────────────────────────────────────────────
// Key routing
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.err = ""

	switch m.mode {
	case ModeMenu:
		return m.handleMenuKey(msg)
	case ModeCategory:
		return m.handleCategoryKey(msg)
	case ModeOptionPick:
		return m.handleOptionPickKey(msg)
	case ModeBoolPick:
		return m.handleBoolPickKey(msg)
	case ModeTextInput:
		return m.handleTextInputKey(msg)
	case ModeCommandEdit:
		return m.handleCommandEditKey(msg)
	case ModeCommandAdd:
		return m.handleCommandAddKey(msg)
	case ModeBYOK:
		return m.handleBYOKKey(msg)
	}
	return m, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Menu
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	n := len(menuEntries)
	switch msg.String() {
	case "up", "k":
		if m.menuCursor > 0 {
			m.menuCursor--
		}
	case "down", "j":
		if m.menuCursor < n-1 {
			m.menuCursor++
		}
	case "enter":
		entry := menuEntries[m.menuCursor]
		return m.enterCategory(entry.cat)
	}
	return m, nil
}

func (m Model) enterCategory(cat Category) (tea.Model, tea.Cmd) {
	m.currentCat = cat
	m.catCursor = 0
	m.err = ""
	m.flash = ""

	switch cat {
	case CatBYOK:
		m.mode = ModeBYOK
		m.byokStep = WizProvider
		return m, loadProviderGroups()

	case CatCommands:
		m.mode = ModeCommandEdit
		m.cmdFocusCol = 0
		m.cmdCursor = 0
		m.allowCmds = append([]string{}, m.settings.CommandAllowlist...)
		m.denyCmds = append([]string{}, m.settings.CommandDenylist...)
		return m, nil

	default:
		m.mode = ModeCategory
		return m, nil
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Category list
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) handleCategoryKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	defs := categorySettings[m.currentCat]

	switch msg.String() {
	case "up", "k":
		if m.catCursor > 0 {
			m.catCursor--
		}
	case "down", "j":
		if m.catCursor < len(defs)-1 {
			m.catCursor++
		}
	case "enter":
		if len(defs) == 0 {
			break
		}
		def := defs[m.catCursor]
		return m.enterSettingEdit(def)
	case "esc":
		m.mode = ModeMenu
	}
	return m, nil
}

func (m Model) enterSettingEdit(def SettingDef) (tea.Model, tea.Cmd) {
	switch def.Kind {
	case KindEnum:
		current := m.settings.GetField(def.Key)
		items := make([]listItem, len(def.Options))
		cursor := 0
		for i, o := range def.Options {
			items[i] = listItem{label: o.Label, value: o.Value, sub: o.Desc}
			if o.Value == current {
				cursor = i
			}
		}
		m.optionList = newList(items, false, listHeight(m.height))
		m.optionList.cursor = cursor
		m.mode = ModeOptionPick

	case KindBool:
		m.optionList = newList([]listItem{
			{label: "Enable", value: "true"},
			{label: "Disable", value: "false"},
		}, false, 4)
		b := m.settings.GetBool(def.Key)
		if b != nil && !*b {
			m.optionList.cursor = 1
		}
		m.mode = ModeBoolPick

	case KindText:
		m.textInput.Reset()
		m.textInput.Placeholder = def.Default
		m.textInput.SetValue(m.settings.GetField(def.Key))
		m.textInput.Focus()
		m.mode = ModeTextInput
	}
	return m, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Option picker (enum)
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) handleOptionPickKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.optionList.up()
	case "down", "j":
		m.optionList.down()
	case "enter":
		val := m.optionList.items[m.optionList.cursor].value

		if val == "__custom__" {
			def := m.currentSettingDef()
			m.textInput.Reset()
			m.textInput.Placeholder = "/path/to/sound.wav"
			m.textInput.SetValue("")
			m.textInput.Focus()
			m.customInput = true
			_ = def
			m.mode = ModeTextInput
			return m, nil
		}

		m.customInput = false
		def := m.currentSettingDef()
		m.settings.SetField(def.Key, val)
		m.mode = ModeCategory
		return m, saveSettings(m.settings, m.rawCfg)

	case "esc":
		m.mode = ModeCategory
	}
	return m, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Bool picker
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) handleBoolPickKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.optionList.up()
	case "down", "j":
		m.optionList.down()
	case "enter":
		val := m.optionList.items[m.optionList.cursor].value == "true"
		def := m.currentSettingDef()
		m.settings.SetBool(def.Key, val)
		m.mode = ModeCategory
		return m, saveSettings(m.settings, m.rawCfg)
	case "esc":
		m.mode = ModeCategory
	}
	return m, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Text input (settings)
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) handleTextInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		val := strings.TrimSpace(m.textInput.Value())
		def := m.currentSettingDef()
		if val == "" {
			val = def.Default
		}
		m.settings.SetField(def.Key, val)
		m.textInput.Blur()
		m.customInput = false
		m.mode = ModeCategory
		return m, saveSettings(m.settings, m.rawCfg)
	case "esc":
		m.textInput.Blur()
		m.customInput = false
		m.mode = ModeCategory
	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Command policy editor
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) handleCommandEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	list := m.activeCommandList()
	n := len(list)

	switch msg.String() {
	case "tab":
		m.cmdFocusCol = 1 - m.cmdFocusCol
		m.cmdCursor = 0
	case "up", "k":
		if m.cmdCursor > 0 {
			m.cmdCursor--
		}
	case "down", "j":
		if m.cmdCursor < n-1 {
			m.cmdCursor++
		}
	case "a":
		m.cmdInput.Reset()
		m.cmdInput.Focus()
		m.mode = ModeCommandAdd
		return m, nil
	case "d", "backspace":
		if n > 0 && m.cmdCursor < n {
			if m.cmdFocusCol == 0 {
				m.allowCmds = append(m.allowCmds[:m.cmdCursor], m.allowCmds[m.cmdCursor+1:]...)
			} else {
				m.denyCmds = append(m.denyCmds[:m.cmdCursor], m.denyCmds[m.cmdCursor+1:]...)
			}
			if m.cmdCursor > 0 {
				m.cmdCursor--
			}
		}
	case "esc":
		m.settings.CommandAllowlist = m.allowCmds
		m.settings.CommandDenylist = m.denyCmds
		m.mode = ModeMenu
		return m, saveSettings(m.settings, m.rawCfg)
	}
	return m, nil
}

func (m Model) handleCommandAddKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		val := strings.TrimSpace(m.cmdInput.Value())
		if val != "" {
			if m.cmdFocusCol == 0 {
				m.allowCmds = append(m.allowCmds, val)
			} else {
				m.denyCmds = append(m.denyCmds, val)
			}
		}
		m.cmdInput.Blur()
		m.mode = ModeCommandEdit
	case "esc":
		m.cmdInput.Blur()
		m.mode = ModeCommandEdit
	default:
		var cmd tea.Cmd
		m.cmdInput, cmd = m.cmdInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) activeCommandList() []string {
	if m.cmdFocusCol == 0 {
		return m.allowCmds
	}
	return m.denyCmds
}

// ─────────────────────────────────────────────────────────────────────────────
// BYOK wizard key handling
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) handleBYOKKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.byokStep {

	case WizProvider:
		switch msg.String() {
		case "esc":
			m.mode = ModeMenu
		case "up", "k":
			m.providerList.up()
		case "down", "j":
			m.providerList.down()
		case "enter":
			if len(m.providerList.items) == 0 {
				break
			}
			return m.wizHandleProviderSelect(m.providerList.items[m.providerList.cursor].value)
		}

	case WizGroupDetail:
		switch msg.String() {
		case "esc":
			m.byokStep = WizProvider
		case "up", "k":
			m.detailList.up()
		case "down", "j":
			m.detailList.down()
		case "enter":
			return m.wizHandleGroupAction(m.detailList.items[m.detailList.cursor].value)
		}

	case WizModelEdit:
		switch msg.String() {
		case "esc":
			m.byokStep = WizGroupDetail
			g := m.providerGroups[m.currentGroupIdx]
			m.detailList = buildGroupDetailList(g)
		case "up", "k":
			m.detailList.up()
		case "down", "j":
			m.detailList.down()
		case "enter":
			field := m.detailList.items[m.detailList.cursor].value
			return m.wizEnterModelField(field)
		}

	case WizModelField:
		return m.handleModelFieldKey(msg)

	case WizURL:
		switch msg.String() {
		case "esc":
			m.byokStep = WizProvider
			m.textInput.Blur()
		case "enter":
			return m.wizSubmitURL()
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	case WizTitle:
		switch msg.String() {
		case "esc":
			m.byokStep = WizURL
			m.focusInput("", m.baseURL)
		case "enter":
			return m.wizSubmitTitle()
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	case WizKey:
		switch msg.String() {
		case "esc":
			m.byokStep = WizTitle
			m.focusInput("", m.displayTitle)
		case "enter":
			return m.wizSubmitKey()
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	case WizModels:
		switch msg.String() {
		case "esc":
			m.byokStep = WizKey
			m.focusInput("sk-... or ${ENV_VAR}", "")
		case "up", "k":
			m.modelList.up()
		case "down", "j":
			m.modelList.down()
		case " ":
			m.modelList.toggleCurrent()
		case "enter":
			if m.modelList.multi {
				sel := m.modelList.selectedValues()
				if len(sel) == 0 {
					m.err = "select at least one model  (space to toggle)"
					break
				}
				m.err = ""
				m.selectedModels = sel
				m.byokStep = WizSettingsTokens
				m.focusInput("16384", strconv.Itoa(m.maxOutputTokens))
			} else {
				val := strings.TrimSpace(m.textInput.Value())
				if val == "" {
					m.err = "enter a model ID"
					break
				}
				m.selectedModels = []string{val}
				m.byokStep = WizSettingsTokens
				m.focusInput("16384", strconv.Itoa(m.maxOutputTokens))
			}
		default:
			if len(m.availableModels) == 0 {
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		}

	case WizSettingsTokens:
		switch msg.String() {
		case "esc":
			m.byokStep = WizModels
		case "enter":
			val := strings.TrimSpace(m.textInput.Value())
			if val == "" {
				val = "16384"
			}
			n, err := strconv.Atoi(val)
			if err != nil || n < 1 {
				m.err = "enter a valid number (e.g. 16384)"
				break
			}
			m.maxOutputTokens = n
			m.byokStep = WizSettingsImages
			m.detailList = buildImagesList()
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	case WizSettingsImages:
		switch msg.String() {
		case "esc":
			m.byokStep = WizSettingsTokens
			m.focusInput("16384", strconv.Itoa(m.maxOutputTokens))
		case "up", "k":
			m.detailList.up()
		case "down", "j":
			m.detailList.down()
		case "enter":
			m.supportsImages = m.detailList.items[m.detailList.cursor].value == "yes"
			m.byokStep = WizConfirm
			m.detailList = buildConfirmList()
		}

	case WizConfirm:
		switch msg.String() {
		case "esc":
			m.byokStep = WizSettingsImages
			m.detailList = buildImagesList()
		case "up", "k":
			m.detailList.up()
		case "down", "j":
			m.detailList.down()
		case "enter":
			if m.detailList.items[m.detailList.cursor].value == "yes" {
				m.byokStep = WizSaving
				return m, wizSaveAll(m)
			}
			m.mode = ModeMenu
		}

	case WizDone:
		switch msg.String() {
		case "up", "k":
			m.detailList.up()
		case "down", "j":
			m.detailList.down()
		case "enter":
			switch m.detailList.items[m.detailList.cursor].value {
			case "same":
				m.selectedModels = nil
				m.byokStep = WizFetching
				return m, cmdFetchModels(m)
			case "other":
				m.byokStep = WizProvider
				return m, loadProviderGroups()
			case "exit":
				m.mode = ModeMenu
			}
		}
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// BYOK wizard submission helpers
// ─────────────────────────────────────────────────────────────────────────────

func (m Model) wizHandleProviderSelect(value string) (tea.Model, tea.Cmd) {
	// Existing provider group
	if strings.HasPrefix(value, "group:") {
		prefix := strings.TrimPrefix(value, "group:")
		for i, g := range m.providerGroups {
			if g.Prefix == prefix {
				m.currentGroupIdx = i
				m.providerKey = g.Prefix
				if len(g.Models) > 0 {
					first := g.Models[0]
					m.baseURL = first.BaseURL
					m.providerType = first.Provider
					m.apiKey = first.APIKey
					m.displayTitle = groupDisplayName(g)
				}
				// Look up models endpoint from known providers
				if p := providers.Get(g.Prefix); p != nil {
					m.modelsEndpoint = p.ModelsEndpoint
					m.noAuth = p.NoAuth
				} else {
					m.modelsEndpoint = "/models"
					m.noAuth = false
				}
				m.byokStep = WizGroupDetail
				m.detailList = buildGroupDetailList(g)
				return m, nil
			}
		}
	}

	// Known provider template
	p := providers.Get(value)
	if p == nil {
		return m, nil
	}
	m.providerKey = value
	m.providerName = p.Name
	m.baseURL = p.BaseURL
	m.providerType = p.Type
	m.noAuth = p.NoAuth
	m.modelsEndpoint = p.ModelsEndpoint

	if p.RequiresBaseURL {
		m.byokStep = WizURL
		m.focusInput("https://", "")
		return m, nil
	}
	m.extractedName = p.Name
	m.displayTitle = p.Name
	if p.NoAuth {
		m.apiKey = "not-needed"
		m.byokStep = WizFetching
		return m, cmdFetchModels(m)
	}
	m.byokStep = WizKey
	m.focusInput("sk-... or ${ENV_VAR}", "")
	return m, nil
}

func (m Model) wizHandleGroupAction(action string) (tea.Model, tea.Cmd) {
	switch {
	case strings.HasPrefix(action, "model:"):
		modelID := strings.TrimPrefix(action, "model:")
		g := m.providerGroups[m.currentGroupIdx]
		for _, model := range g.Models {
			if model.ID == modelID {
				m.editingModel = model
				m.byokStep = WizModelEdit
				m.detailList = buildModelEditList(model)
				return m, nil
			}
		}
	case action == "add-models":
		m.selectedModels = nil
		m.byokStep = WizFetching
		return m, cmdFetchModels(m)
	case action == "back":
		m.byokStep = WizProvider
	}
	return m, nil
}

func (m Model) wizEnterModelField(field string) (tea.Model, tea.Cmd) {
	m.editFieldKey = field

	switch field {
	case "displayName":
		m.byokStep = WizModelField
		m.focusInput("Display name", m.editingModel.DisplayName)
	case "model":
		m.byokStep = WizModelField
		m.focusInput("Model ID", m.editingModel.Model)
	case "baseUrl":
		m.byokStep = WizModelField
		m.focusInput("Base URL", m.editingModel.BaseURL)
	case "apiKey":
		m.byokStep = WizModelField
		m.focusInput("API Key", m.editingModel.APIKey)
	case "maxOutputTokens":
		m.byokStep = WizModelField
		m.focusInput("Max tokens", strconv.Itoa(m.editingModel.MaxOutputTokens))
	case "provider":
		m.byokStep = WizModelField
		items := make([]listItem, len(providers.ProviderTypes))
		cursor := 0
		for i, pt := range providers.ProviderTypes {
			items[i] = listItem{label: pt.Label, value: pt.Value}
			if pt.Value == m.editingModel.Provider {
				cursor = i
			}
		}
		m.detailList = newList(items, false, 5)
		m.detailList.cursor = cursor
	case "supportsImages":
		m.byokStep = WizModelField
		m.detailList = buildImagesList()
		if m.editingModel.SupportsImages {
			m.detailList.cursor = 1
		}
	case "delete":
		m.byokStep = WizModelField
		m.detailList = newList([]listItem{
			{label: "Yes, delete", value: "yes"},
			{label: "Cancel", value: "no"},
		}, false, 3)
	case "back":
		m.byokStep = WizGroupDetail
		g := m.providerGroups[m.currentGroupIdx]
		m.detailList = buildGroupDetailList(g)
	}
	return m, nil
}

func (m Model) handleModelFieldKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.editFieldKey {
	case "displayName", "model", "baseUrl", "apiKey", "maxOutputTokens":
		switch msg.String() {
		case "esc":
			m.textInput.Blur()
			m.byokStep = WizModelEdit
			m.detailList = buildModelEditList(m.editingModel)
		case "enter":
			return m.wizSaveModelField()
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	case "provider", "supportsImages", "delete":
		switch msg.String() {
		case "esc":
			m.byokStep = WizModelEdit
			m.detailList = buildModelEditList(m.editingModel)
		case "up", "k":
			m.detailList.up()
		case "down", "j":
			m.detailList.down()
		case "enter":
			return m.wizSaveModelField()
		}
	}
	return m, nil
}

func (m Model) wizSaveModelField() (tea.Model, tea.Cmd) {
	switch m.editFieldKey {
	case "displayName":
		val := strings.TrimSpace(m.textInput.Value())
		if val != "" {
			m.editingModel.DisplayName = val
		}
		m.textInput.Blur()
	case "model":
		val := strings.TrimSpace(m.textInput.Value())
		if val != "" {
			m.editingModel.Model = val
		}
		m.textInput.Blur()
	case "baseUrl":
		val := strings.TrimSpace(m.textInput.Value())
		if val != "" {
			m.editingModel.BaseURL = val
		}
		m.textInput.Blur()
	case "apiKey":
		val := strings.TrimSpace(m.textInput.Value())
		if val != "" {
			m.editingModel.APIKey = val
		}
		m.textInput.Blur()
	case "maxOutputTokens":
		val := strings.TrimSpace(m.textInput.Value())
		n, err := strconv.Atoi(val)
		if err != nil || n < 1 {
			m.err = "enter a valid number"
			return m, nil
		}
		m.editingModel.MaxOutputTokens = n
		m.textInput.Blur()
	case "provider":
		m.editingModel.Provider = m.detailList.items[m.detailList.cursor].value
	case "supportsImages":
		m.editingModel.SupportsImages = m.detailList.items[m.detailList.cursor].value == "yes"
	case "delete":
		if m.detailList.items[m.detailList.cursor].value == "yes" {
			return m, wizDeleteModel(m.editingModel.ID)
		}
		m.byokStep = WizModelEdit
		m.detailList = buildModelEditList(m.editingModel)
		return m, nil
	}

	m.byokStep = WizModelEdit
	m.detailList = buildModelEditList(m.editingModel)
	return m, wizPersistModel(m.editingModel)
}

func (m Model) wizSubmitURL() (tea.Model, tea.Cmd) {
	raw := strings.TrimSpace(m.textInput.Value())
	normalized, err := api.NormalizeURL(raw)
	if err != nil {
		m.err = err.Error()
		return m, nil
	}
	m.baseURL = normalized
	m.extractedName = api.ExtractProviderName(normalized)
	m.displayTitle = m.extractedName
	m.byokStep = WizTitle
	m.focusInput("Display name", m.extractedName)
	return m, nil
}

func (m Model) wizSubmitTitle() (tea.Model, tea.Cmd) {
	val := strings.TrimSpace(m.textInput.Value())
	if val == "" || strings.Contains(val, "://") {
		val = m.extractedName
	}
	m.displayTitle = val
	m.byokStep = WizKey
	m.focusInput("sk-... or ${ENV_VAR}", "")
	return m, nil
}

func (m Model) wizSubmitKey() (tea.Model, tea.Cmd) {
	key := strings.TrimSpace(m.textInput.Value())
	if key == "" {
		m.err = "API key cannot be empty"
		return m, nil
	}
	m.apiKey = key
	m.byokStep = WizFetching
	return m, cmdFetchModels(m)
}

func (m *Model) focusInput(placeholder, defaultVal string) {
	m.textInput.Reset()
	m.textInput.Placeholder = placeholder
	m.textInput.SetValue(defaultVal)
	m.textInput.Focus()
}

// ─────────────────────────────────────────────────────────────────────────────
// List builders
// ─────────────────────────────────────────────────────────────────────────────

func buildProviderList(groups []config.ProviderGroup) customList {
	var items []listItem

	// Existing provider groups first
	for _, g := range groups {
		name := groupDisplayName(g)
		sub := fmt.Sprintf("%d model(s)", len(g.Models))
		items = append(items, listItem{label: name, value: "group:" + g.Prefix, sub: sub})
	}

	// Then known provider templates (skip ones that already have a group)
	existing := map[string]bool{}
	for _, g := range groups {
		existing[g.Prefix] = true
	}
	for _, p := range providers.All {
		if existing[p.Key] {
			continue
		}
		items = append(items, listItem{label: p.Provider.Name, value: p.Key})
	}
	return newList(items, false, 12)
}

// groupDisplayName returns the group prefix as its display name.
func groupDisplayName(g config.ProviderGroup) string {
	return g.Prefix
}

func buildGroupDetailList(g config.ProviderGroup) customList {
	var items []listItem
	for _, model := range g.Models {
		dn := model.DisplayName
		if dn == "" {
			dn = model.Model
		}
		items = append(items, listItem{label: dn, value: "model:" + model.ID, sub: model.Model})
	}
	items = append(items,
		listItem{label: "+ Add more models", value: "add-models"},
		listItem{label: "← Back", value: "back"},
	)
	return newList(items, false, 12)
}

func buildModelEditList(model config.ModelConfig) customList {
	apiTypeLabel := model.Provider
	for _, pt := range providers.ProviderTypes {
		if pt.Value == model.Provider {
			apiTypeLabel = pt.Label
			break
		}
	}
	images := "No"
	if model.SupportsImages {
		images = "Yes"
	}
	items := []listItem{
		{label: "Display Name", value: "displayName", sub: model.DisplayName},
		{label: "Model ID", value: "model", sub: model.Model},
		{label: "Base URL", value: "baseUrl", sub: model.BaseURL},
		{label: "API Key", value: "apiKey", sub: maskKey(model.APIKey)},
		{label: "API Type", value: "provider", sub: apiTypeLabel},
		{label: "Max Tokens", value: "maxOutputTokens", sub: strconv.Itoa(model.MaxOutputTokens)},
		{label: "Image Support", value: "supportsImages", sub: images},
		{label: "Delete Model", value: "delete"},
		{label: "← Back", value: "back"},
	}
	return newList(items, false, 12)
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func buildModelList(models []api.ModelInfo) customList {
	items := make([]listItem, len(models))
	for i, m := range models {
		items[i] = listItem{label: m.Name, value: m.ID}
	}
	return newList(items, true, 12)
}

func buildImagesList() customList {
	return newList([]listItem{
		{label: "No", value: "no"},
		{label: "Yes", value: "yes"},
	}, false, 3)
}

func buildConfirmList() customList {
	return newList([]listItem{
		{label: "Save", value: "yes"},
		{label: "Cancel", value: "no"},
	}, false, 3)
}

func buildDoneList() customList {
	return newList([]listItem{
		{label: "Add more models to this provider", value: "same"},
		{label: "Add models for another provider", value: "other"},
		{label: "Back to main menu", value: "exit"},
	}, false, 5)
}

func listHeight(h int) int {
	// Reserve lines for header (~4), footer (2), flash/error (2), padding.
	usable := h - 8
	if usable < 4 {
		return 4
	}
	return usable
}

// ─────────────────────────────────────────────────────────────────────────────
// Async commands
// ─────────────────────────────────────────────────────────────────────────────

func loadAllSettings() tea.Cmd {
	return func() tea.Msg {
		s, raw, _ := config.ReadSettings()
		return settingsLoadedMsg{settings: s, raw: raw}
	}
}

func loadProviderGroups() tea.Cmd {
	return func() tea.Msg {
		groups, _ := config.ReadProviderGroups()
		return groupsLoadedMsg{groups: groups}
	}
}

func saveSettings(s config.Settings, raw map[string]any) tea.Cmd {
	return func() tea.Msg {
		if err := config.WriteSettings(s, raw); err != nil {
			return errMsg{err: err}
		}
		return settingsSavedMsg{}
	}
}

func cmdFetchModels(m Model) tea.Cmd {
	baseURL := m.baseURL
	apiKey := m.apiKey
	modelsEndpoint := m.modelsEndpoint
	providerType := m.providerType
	noAuth := m.noAuth

	return func() tea.Msg {
		models, _ := api.FetchModels(baseURL, apiKey, modelsEndpoint, providerType, noAuth)
		displayNames := make(map[string]string, len(models))
		for i, model := range models {
			dn := api.GetDisplayName(model.ID)
			displayNames[model.ID] = dn
			models[i].Name = fmt.Sprintf("%s  \x1b[38;5;241m%s\x1b[0m", dn, model.ID)
		}
		return modelsLoadedMsg{models: models, displayNames: displayNames}
	}
}

func wizSaveAll(m Model) tea.Cmd {
	providerKey := m.providerKey
	return func() tea.Msg {
		nextIdx, _ := config.GetNextModelIndex(providerKey)
		for i, modelID := range m.selectedModels {
			dn := m.modelDisplayNames[modelID]
			if dn == "" {
				dn = modelID
			}
			idx := nextIdx + i
			cfg := config.ModelConfig{
				ID:              config.GenerateModelID(providerKey, idx),
				Index:           idx,
				Model:           modelID,
				DisplayName:     fmt.Sprintf("%s [%s]", dn, m.displayTitle),
				BaseURL:         m.baseURL,
				APIKey:          m.apiKey,
				Provider:        m.providerType,
				MaxOutputTokens: m.maxOutputTokens,
				SupportsImages:  m.supportsImages,
			}
			if err := config.AddModelToSettings(cfg); err != nil {
				return errMsg{err: err}
			}
		}
		return byokSavedMsg{path: config.SettingsPath()}
	}
}

type modelDeletedMsg struct{}

func wizPersistModel(model config.ModelConfig) tea.Cmd {
	return func() tea.Msg {
		if err := config.AddModelToSettings(model); err != nil {
			return errMsg{err: err}
		}
		return settingsSavedMsg{}
	}
}

func wizDeleteModel(id string) tea.Cmd {
	return func() tea.Msg {
		if err := config.DeleteModelFromSettings(id); err != nil {
			return errMsg{err: err}
		}
		return modelDeletedMsg{}
	}
}

func clearFlashAfter() tea.Cmd {
	return tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
		return clearFlashMsg{}
	})
}
