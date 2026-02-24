package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/kaan-escober/wrench/internal/theme"
)

// ─── Category list ─────────────────────────────────────────────────────────────

func (m Model) viewCategory() string {
	defs := categorySettings[m.currentCat]
	badge := m.catBadge()

	var sb strings.Builder
	sb.WriteString(viewHeader(badge, m.catSubtitle()))

	nameStyle := lipgloss.NewStyle().Width(24)
	for i, def := range defs {
		isCursor := i == m.catCursor

		current := m.settingValueDisplay(def)
		desc := m.settingDesc(def)

		var nameStr string
		if isCursor {
			nameStr = nameStyle.Inherit(theme.Accent).Bold(true).Render(def.Label)
		} else {
			nameStr = nameStyle.Inherit(theme.Primary).Render(def.Label)
		}

		cursor := "  "
		if isCursor {
			cursor = theme.Accent.Render("> ")
		}

		valStr := theme.Teal.Render(current)
		if desc != "" && isCursor {
			valStr = theme.Teal.Render(current) + "  " + theme.Muted.Render(desc)
		}

		sb.WriteString(cursor + nameStr + "  " + valStr + "\n")
	}
	return sb.String()
}

func (m Model) catBadge() string {
	for _, e := range menuEntries {
		if e.cat == m.currentCat {
			return e.badge
		}
	}
	return "CFG"
}

func (m Model) catSubtitle() string {
	switch m.currentCat {
	case CatModel:
		return "Choose the default AI model and reasoning behaviour"
	case CatAutonomy:
		return "How proactively Droid executes commands"
	case CatDisplay:
		return "Control how Droid presents information in the TUI"
	case CatSound:
		return "Audio feedback for Droid events"
	case CatSecurity:
		return "Shield, commit attribution, and process controls"
	case CatBehavior:
		return "Session sync, IDE, spec saving, and experimental features"
	}
	return ""
}

func (m Model) settingValueDisplay(def SettingDef) string {
	switch def.Kind {
	case KindEnum:
		v := m.settings.GetField(def.Key)
		if v == "" {
			return def.Default + " (default)"
		}
		return v
	case KindBool:
		b := m.settings.GetBool(def.Key)
		if b == nil {
			return def.Default + " (default)"
		}
		if *b {
			return "enabled"
		}
		return "disabled"
	case KindText:
		v := m.settings.GetField(def.Key)
		if v == "" {
			return def.Default + " (default)"
		}
		return v
	}
	return ""
}

func (m Model) settingDesc(def SettingDef) string {
	if len(def.Options) > 0 && def.Options[0].Desc != "" {
		// Only show desc for bool kinds (the single-option desc field)
		if def.Kind == KindBool {
			return def.Options[0].Desc
		}
	}
	return ""
}

// ─── Option picker ─────────────────────────────────────────────────────────────

func (m Model) viewOptionPick() string {
	def := m.currentSettingDef()
	current := m.settings.GetField(def.Key)
	if current == "" {
		current = def.Default
	}

	header := viewHeader(def.Label, "Currently: "+theme.Teal.Render(current))
	list := m.optionList.render(true)

	// For custom path sound: show a note
	note := ""
	if def.Key == "completionSound" || def.Key == "awaitingInputSound" {
		note = "\n" + theme.Muted.Render("  Select \"Custom file path...\" to enter your own audio file")
	}

	return header + list + note
}

// ─── Bool picker ───────────────────────────────────────────────────────────────

func (m Model) viewBoolPick() string {
	def := m.currentSettingDef()
	b := m.settings.GetBool(def.Key)

	current := def.Default + " (default)"
	if b != nil {
		if *b {
			current = "enabled"
		} else {
			current = "disabled"
		}
	}

	desc := ""
	if len(def.Options) > 0 {
		desc = def.Options[0].Desc
	}

	header := viewHeader(def.Label, "Currently: "+theme.Teal.Render(current))

	var note string
	if desc != "" {
		note = theme.Muted.Render("  "+desc) + "\n\n"
	}

	return header + note + m.optionList.render(true)
}

// ─── Text input ────────────────────────────────────────────────────────────────

func (m Model) viewTextInput() string {
	def := m.currentSettingDef()

	var subtitle string
	if m.customInput {
		// Sound custom path
		subtitle = "Enter a file path to your audio file (.wav, .mp3, .ogg)"
	} else {
		current := m.settings.GetField(def.Key)
		if current == "" {
			current = def.Default
		}
		subtitle = "Currently: " + theme.Teal.Render(current)
	}

	inputLine := theme.PromptStr() + m.textInput.View()

	return viewHeader(def.Label, subtitle) + inputLine
}

