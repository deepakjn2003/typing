package styles

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// ANSI semantic colors — each terminal maps these to its own palette.
var (
	ColorPrimary   = lipgloss.Color("4")
	ColorAccent    = lipgloss.Color("6")
	ColorSecondary = lipgloss.Color("5")
	ColorBorder    = lipgloss.Color("8")
	ColorBorderDim = lipgloss.Color("8")

	ColorSuccess = lipgloss.Color("2")
	ColorError   = lipgloss.Color("1")
	ColorWarning = lipgloss.Color("3")
	ColorInfo    = lipgloss.Color("6")

	ColorCorrect   = ColorSuccess
	ColorIncorrect = ColorError
	ColorUntyped   = lipgloss.Color("8")
)

// No background fills — the terminal provides the canvas.
var (
	AppStyle       = lipgloss.NewStyle()
	FooterBarStyle = lipgloss.NewStyle()
)

var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary)

	SubtitleStyle = lipgloss.NewStyle().
		Faint(true).
		Italic(true)

	BodyStyle = lipgloss.NewStyle()

	MutedStyle = lipgloss.NewStyle().Faint(true)

	DimStyle = lipgloss.NewStyle().Faint(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(ColorSuccess)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(ColorError)

	WarningStyle = lipgloss.NewStyle().
		Foreground(ColorWarning)

	InfoStyle = lipgloss.NewStyle().
		Foreground(ColorInfo)
)

var (
	CorrectCharStyle = lipgloss.NewStyle().
		Foreground(ColorCorrect)

	IncorrectCharStyle = lipgloss.NewStyle().
		Foreground(ColorIncorrect).
		Underline(true)

	CursorCharStyle = lipgloss.NewStyle().
		Reverse(true).
		Bold(true)

	UntypedCharStyle = lipgloss.NewStyle().
		Faint(true)
)

var (
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2)

	FocusedBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 1)

	StatusBarStyle = lipgloss.NewStyle().
		Faint(true).
		Padding(0, 1)
)

var (
	MenuItemStyle = lipgloss.NewStyle().
		Padding(0, 2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 2)

	KeyHintStyle = lipgloss.NewStyle().Faint(true)

	KeyStyle = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	SeparatorStyle = lipgloss.NewStyle().
		Faint(true)
)

var (
	StatLabelStyle = lipgloss.NewStyle().Faint(true)

	StatValueStyle = lipgloss.NewStyle().Bold(true)

	BigStatStyle = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)
)

func Separator(width int) string {
	if width <= 0 {
		return ""
	}
	return SeparatorStyle.Render(strings.Repeat("─", width))
}
