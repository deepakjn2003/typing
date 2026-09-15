package components

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type KeyBinding struct {
	Key  string
	Desc string
}

type HelpBar struct {
	Bindings []KeyBinding
	Width    int
}

func NewHelpBar(bindings []KeyBinding) HelpBar {
	return HelpBar{Bindings: bindings}
}

func (h *HelpBar) SetWidth(w int) {
	h.Width = w
}

func (h HelpBar) View() string {
	if len(h.Bindings) == 0 {
		return ""
	}

	var parts []string
	for _, b := range h.Bindings {
		key := styles.KeyStyle.Render(b.Key)
		desc := styles.KeyHintStyle.Render(b.Desc)
		parts = append(parts, key+" "+desc)
	}

	joined := strings.Join(parts, styles.DimStyle.Render("  •  "))

	return lipgloss.NewStyle().
		Padding(0, 1).
		Faint(true).
		Render(joined)
}
