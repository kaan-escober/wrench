package ui

import tea "github.com/charmbracelet/bubbletea"

// Run starts the TUI program and blocks until it exits.
func Run() error {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)
	_, err := p.Run()
	return err
}
