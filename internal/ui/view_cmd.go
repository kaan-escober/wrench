package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/kaan-escober/wrench/internal/theme"
)

func (m Model) viewCommandEdit() string {
	header := viewHeader("CMD", "Commands that always allow or always deny, overriding autonomy")

	// Column widths
	colW := (m.width - 6) / 2
	if colW < 20 {
		colW = 20
	}

	allowHeader := m.colHeader("ALLOWLIST", m.cmdFocusCol == 0)
	denyHeader := m.colHeader("DENYLIST", m.cmdFocusCol == 1)

	allowBody := m.renderCmdList(m.allowCmds, m.cmdFocusCol == 0, colW)
	denyBody := m.renderCmdList(m.denyCmds, m.cmdFocusCol == 1, colW)

	// Pad columns to same height
	allowLines := strings.Split(allowBody, "\n")
	denyLines := strings.Split(denyBody, "\n")
	maxLines := len(allowLines)
	if len(denyLines) > maxLines {
		maxLines = len(denyLines)
	}
	for len(allowLines) < maxLines {
		allowLines = append(allowLines, "")
	}
	for len(denyLines) < maxLines {
		denyLines = append(denyLines, "")
	}

	// Render side-by-side
	var sb strings.Builder
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		allowHeader,
		strings.Repeat(" ", colW+2-lipgloss.Width(allowHeader)+2),
		denyHeader,
	) + "\n")

	for i := range allowLines {
		al := allowLines[i]
		dl := ""
		if i < len(denyLines) {
			dl = denyLines[i]
		}
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(colW+2).Render(al),
			"  ",
			lipgloss.NewStyle().Width(colW+2).Render(dl),
		) + "\n")
	}

	// Adding mode
	if m.mode == ModeCommandAdd {
		col := "ALLOWLIST"
		if m.cmdFocusCol == 1 {
			col = "DENYLIST"
		}
		sb.WriteString("\n" + theme.Muted.Render(fmt.Sprintf("  Adding to %s:", col)) + "\n")
		sb.WriteString("  " + theme.PromptStr() + m.cmdInput.View())
	}

	return header + sb.String()
}

func (m Model) colHeader(label string, active bool) string {
	if active {
		return theme.Badge.Render(" "+label+" ") + "\n"
	}
	return lipgloss.NewStyle().
		Foreground(theme.ColorMuted).
		Padding(0, 1).
		Render(label) + "\n"
}

func (m Model) renderCmdList(cmds []string, focused bool, colW int) string {
	if len(cmds) == 0 {
		empty := theme.Muted.Render("  (none)")
		return empty
	}
	var sb strings.Builder
	for i, cmd := range cmds {
		isCursor := focused && i == m.cmdCursor

		display := cmd
		if len(display) > colW-4 {
			display = display[:colW-7] + "..."
		}

		var line string
		if isCursor {
			line = theme.Accent.Render("> ") + theme.Primary.Render(display)
		} else {
			line = "  " + theme.Primary.Render(display)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}
