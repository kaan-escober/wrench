package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/kaan-escober/wrench/internal/config"
	"github.com/kaan-escober/wrench/internal/theme"
)

func (m Model) viewMenu() string {
	return m.renderLogo() + "\n\n" + m.renderMenuRows()
}

func (m Model) renderLogo() string {
	title := theme.Accent.Bold(true).Render("D R O I D")
	cfg := theme.Badge.Render("CONFIG")
	line1 := title + "  " + cfg
	line2 := theme.Muted.Render("  " + config.SettingsPath())
	return line1 + "\n" + line2
}

// labelCol is the fixed width for the label column — wide enough for the longest entry.
const labelCol = 20

// badgeTextWidth is the widest badge text across all menu entries (e.g. "BYOK" = 4).
// We pad shorter badges to this width so all badge blocks are the same rendered size.
var badgeTextWidth = func() int {
	w := 0
	for _, e := range menuEntries {
		if len(e.badge) > w {
			w = len(e.badge)
		}
	}
	return w
}()

var badgeMuted = lipgloss.NewStyle().
	Foreground(theme.ColorMuted).
	Background(lipgloss.Color("#2A2B3A")).
	Bold(true).
	Padding(0, 1)

// padBadge centres/pads badge text to badgeTextWidth so all blocks render identically.
func padBadge(s string) string {
	for len(s) < badgeTextWidth {
		s += " "
	}
	return s
}

func (m Model) renderMenuRows() string {
	var sb strings.Builder
	labelStyle := lipgloss.NewStyle().Width(labelCol)

	for i, entry := range menuEntries {
		isCursor := i == m.menuCursor

		text := padBadge(entry.badge)
		var badge string
		if isCursor {
			badge = theme.Badge.Render(text)
		} else {
			badge = badgeMuted.Render(text)
		}

		var label string
		if isCursor {
			label = labelStyle.Inherit(theme.Accent).Bold(true).Render(entry.label)
		} else {
			label = labelStyle.Inherit(theme.Primary).Render(entry.label)
		}

		summary := m.menuSummary(entry.cat)

		cursor := "  "
		if isCursor {
			cursor = theme.Accent.Render("> ")
		}

		row := cursor + badge + "  " + label + "  " + theme.Muted.Render(summary)
		sb.WriteString(row + "\n")
	}
	return sb.String()
}

func (m Model) menuSummary(cat Category) string {
	s := m.settings
	raw := m.rawCfg

	switch cat {
	case CatBYOK:
		n := customModelCount(raw)
		if n == 0 {
			return "no models configured"
		}
		return fmt.Sprintf("%d model(s) configured", n)

	case CatModel:
		model := orDef(s.Model, "opus")
		effort := orDef(s.ReasoningEffort, "auto")
		return model + "  ·  " + effort

	case CatAutonomy:
		return orDef(s.AutonomyLevel, "normal")

	case CatDisplay:
		diff := orDef(s.DiffMode, "github")
		todo := orDef(s.TodoDisplayMode, "pinned")
		return diff + "  ·  " + todo

	case CatSound:
		c := orDef(s.CompletionSound, "fx-ok01")
		f := orDef(s.SoundFocusMode, "always")
		return c + "  ·  " + f

	case CatSecurity:
		shield := dotBool(s.EnableDroidShield, true, true) + " shield"
		coauth := dotBool(s.IncludeCoAuthoredByDroid, true, true) + " co-author"
		bgProc := dotBool(s.AllowBackgroundProcesses, false, true) + " bg-proc"
		return shield + "  " + coauth + "  " + bgProc

	case CatBehavior:
		cloud := dotBool(s.CloudSessionSync, true, true) + " cloud"
		hooks := dotBool(s.HooksDisabled, false, false) + " hooks"
		droids := dotBool(s.EnableCustomDroids, true, true) + " droids"
		return cloud + "  " + hooks + "  " + droids

	case CatCommands:
		a := len(s.CommandAllowlist)
		d := len(s.CommandDenylist)
		if a == 0 && d == 0 {
			return "factory defaults"
		}
		return fmt.Sprintf("%d allowed  ·  %d denied", a, d)
	}
	return ""
}

// orDef returns val if non-empty, else def.
func orDef(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

// dotBool returns an orange ● or muted ○.
// defaultOn: what the factory default is (true=on).
// activeMeansTrue: orange when value==true.
func dotBool(b *bool, defaultOn bool, activeMeansTrue bool) string {
	effective := defaultOn
	if b != nil {
		effective = *b
	}
	isActive := effective == activeMeansTrue
	if isActive {
		return theme.Accent.Render("●")
	}
	return theme.Muted.Render("○")
}

func customModelCount(raw map[string]any) int {
	if v, ok := raw["customModels"]; ok {
		if arr, ok := v.([]any); ok {
			return len(arr)
		}
	}
	return 0
}
