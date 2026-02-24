package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/kaan-escober/wrench/internal/theme"
)

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}
	var body string
	switch m.mode {
	case ModeMenu:
		body = m.viewMenu()
	case ModeCategory:
		body = m.viewCategory()
	case ModeOptionPick:
		body = m.viewOptionPick()
	case ModeBoolPick:
		body = m.viewBoolPick()
	case ModeTextInput:
		body = m.viewTextInput()
	case ModeCommandEdit, ModeCommandAdd:
		body = m.viewCommandEdit()
	case ModeBYOK:
		body = m.viewBYOK()
	}

	parts := []string{body}
	if e := m.viewErr(); e != "" {
		parts = append(parts, e)
	}
	if f := m.viewFlash(); f != "" {
		parts = append(parts, f)
	}
	footer := m.viewFooter()

	// Pad body so footer is pinned to the bottom of the screen.
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	contentH := lipgloss.Height(content)
	footerH := lipgloss.Height(footer)
	gap := m.height - contentH - footerH
	if gap > 0 {
		content += strings.Repeat("\n", gap)
	}
	return content + footer
}

// ─── Shared helpers ───────────────────────────────────────────────────────────

func (m Model) viewErr() string {
	if m.err == "" {
		return ""
	}
	return theme.Error.Render("  △  " + m.err)
}

func (m Model) viewFlash() string {
	if m.flash == "" {
		return ""
	}
	return theme.Success.Render(m.flash)
}

func (m Model) viewFooter() string {
	var hints string
	switch m.mode {
	case ModeMenu:
		hints = "↑↓ navigate  enter · open  ctrl+c quit"
	case ModeCategory:
		hints = "↑↓ navigate  enter · edit  esc · back"
	case ModeOptionPick, ModeBoolPick:
		hints = "↑↓ navigate  enter · select  esc · back"
	case ModeTextInput:
		hints = "enter · confirm  esc · back"
	case ModeCommandEdit:
		hints = "↑↓ navigate  tab · switch  a · add  d · delete  esc · save & back"
	case ModeCommandAdd:
		hints = "enter · add command  esc · cancel"
	case ModeBYOK:
		hints = m.byokFooterHints()
	}

	right := theme.Muted.Render(m.viewModeLabel())
	rightW := lipgloss.Width(right)
	// Truncate hints so the footer never exceeds terminal width.
	maxHints := m.width - rightW - 3 // 3 = minimum gap + space
	if maxHints < 0 {
		maxHints = 0
	}
	if len(hints) > maxHints {
		hints = hints[:maxHints]
	}
	left := theme.Muted.Render(hints)
	gap := m.width - lipgloss.Width(left) - rightW
	if gap < 1 {
		gap = 1
	}
	line := left + strings.Repeat(" ", gap) + right
	sep := theme.Muted.Render(strings.Repeat("─", m.width))
	return sep + "\n" + line
}

func (m Model) viewModeLabel() string {
	switch m.mode {
	case ModeMenu:
		return "DROID CONFIG"
	case ModeBYOK:
		return "BYOK"
	case ModeCommandEdit, ModeCommandAdd:
		return "CMD"
	default:
		defs := categorySettings[m.currentCat]
		if m.catCursor >= 0 && m.catCursor < len(defs) {
			return defs[m.catCursor].Key
		}
		for _, e := range menuEntries {
			if e.cat == m.currentCat {
				return e.badge
			}
		}
	}
	return ""
}

func (m Model) byokFooterHints() string {
	switch m.byokStep {
	case WizProvider, WizGroupDetail:
		return "↑↓ navigate  enter · select  esc · back"
	case WizModelEdit:
		return "↑↓ navigate  enter · edit  esc · back"
	case WizModelField:
		if m.editFieldKey == "provider" || m.editFieldKey == "supportsImages" || m.editFieldKey == "delete" {
			return "↑↓ navigate  enter · select  esc · cancel"
		}
		return "enter · save  esc · cancel"
	case WizModels:
		return "space · toggle  ↑↓ navigate  enter · confirm  esc · back"
	case WizSettingsImages, WizConfirm, WizDone:
		return "↑↓ navigate  enter · select  esc · back"
	default:
		return "enter · confirm  esc · back"
	}
}

// viewHeader renders the orange badge + optional subtitle line for a screen.
func viewHeader(badge, subtitle string) string {
	h := theme.Badge.Render(badge)
	if subtitle != "" {
		h += "\n\n" + theme.Muted.Render(subtitle)
	}
	return h + "\n\n"
}
