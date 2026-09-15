package layout

import (
	"charm.land/lipgloss/v2"
)

func RenderOverlay(base, overlay string, width, height int) string {
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		overlay,
		lipgloss.WithWhitespaceChars(base),
	)
}
