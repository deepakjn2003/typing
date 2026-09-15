package styles

import (
	"os"

	"charm.land/lipgloss/v2"
)

// TerminalIsDark reports whether the terminal uses a dark background.
func TerminalIsDark() bool {
	return lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
}
