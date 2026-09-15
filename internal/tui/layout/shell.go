package layout

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/tui/components"
	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type Shell struct {
	Width, Height int
	Bindings      []components.KeyBinding
	Content       string
	ContentAlign  lipgloss.Position
}

func (s Shell) View() string {
	footer := components.NewHelpBar(s.Bindings)
	footer.SetWidth(s.Width)
	footerView := styles.FooterBarStyle.Width(s.Width).Render(footer.View())

	footerHeight := lipgloss.Height(footerView)
	contentHeight := s.Height - footerHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	maxWidth := ContentWidth(s.Width)
	boundedContent := lipgloss.NewStyle().Width(maxWidth).Render(s.Content)
	hCentered := lipgloss.PlaceHorizontal(s.Width, lipgloss.Center, boundedContent)

	vAlign := s.ContentAlign
	if vAlign != lipgloss.Top && vAlign != lipgloss.Center {
		vAlign = lipgloss.Top
	}

	placed := lipgloss.Place(s.Width, contentHeight, lipgloss.Left, vAlign, hCentered)

	return lipgloss.JoinVertical(lipgloss.Top, placed, footerView)
}

// ContentBlockWidth returns the widest visible line in content, capped at maxWidth.
func ContentBlockWidth(content string, maxWidth int) int {
	if content == "" {
		return maxWidth
	}

	blockWidth := 0
	for _, line := range strings.Split(content, "\n") {
		if w := lipgloss.Width(line); w > blockWidth {
			blockWidth = w
		}
	}

	if blockWidth < 1 {
		return maxWidth
	}
	if blockWidth > maxWidth {
		return maxWidth
	}
	return blockWidth
}
