package theme

import "github.com/charmbracelet/lipgloss"

// Colors — extracted from Factory/Droid screenshots
const (
	ColorAccent  = lipgloss.Color("#FF9E64") // orange — primary accent
	ColorPrimary = lipgloss.Color("#C0CAF5") // off-white lavender — body text
	ColorMuted   = lipgloss.Color("#565F89") // slate — secondary/hints
	ColorSuccess = lipgloss.Color("#9ECE6A") // green
	ColorError   = lipgloss.Color("#F7768E") // red
	ColorTeal    = lipgloss.Color("#7DCFFF") // teal — code/commands
	ColorBlack   = lipgloss.Color("#1A1B26") // badge text bg-color
)

// Base styles
var (
	Accent  = lipgloss.NewStyle().Foreground(ColorAccent)
	Primary = lipgloss.NewStyle().Foreground(ColorPrimary)
	Muted   = lipgloss.NewStyle().Foreground(ColorMuted)
	Success = lipgloss.NewStyle().Foreground(ColorSuccess)
	Error   = lipgloss.NewStyle().Foreground(ColorError)
	Teal    = lipgloss.NewStyle().Foreground(ColorTeal)
	Bold    = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
)

// Badge — solid orange rect, dark text, sharp corners. Matches ASK USER / PLAN etc.
var Badge = lipgloss.NewStyle().
	Background(ColorAccent).
	Foreground(ColorBlack).
	Bold(true).
	Padding(0, 1)

// BadgeSuccess — same shape, green
var BadgeSuccess = lipgloss.NewStyle().
	Background(ColorSuccess).
	Foreground(ColorBlack).
	Bold(true).
	Padding(0, 1)

// BadgeError — same shape, red
var BadgeError = lipgloss.NewStyle().
	Background(ColorError).
	Foreground(ColorBlack).
	Bold(true).
	Padding(0, 1)

// Input — the text input box, thin dim border
var Input = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(ColorMuted).
	Padding(0, 1)

// Prompt — orange `>` before input fields
const Prompt = "❯ "

// PromptStyled renders the orange prompt symbol
func PromptStr() string {
	return Accent.Render(Prompt)
}

// Hint renders a keyboard hint in muted color
func Hint(s string) string {
	return Muted.Render(s)
}
