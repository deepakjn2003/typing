package components

import (
	"strings"

	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type HelpSection struct {
	Title    string
	Bindings []KeyBinding
}

type HelpOverlay struct {
	Sections []HelpSection
	Width    int
}

func NewHelpOverlay(sections []HelpSection, width int) HelpOverlay {
	return HelpOverlay{Sections: sections, Width: width}
}

func DefaultHelpOverlay() HelpOverlay {
	return NewHelpOverlay([]HelpSection{
		{
			Title: "Navigation",
			Bindings: []KeyBinding{
				{Key: "↑/↓", Desc: "move selection"},
				{Key: "enter", Desc: "confirm"},
				{Key: "esc", Desc: "go back"},
			},
		},
		{
			Title: "Commands",
			Bindings: []KeyBinding{
				{Key: ":", Desc: "open command palette"},
				{Key: "ctrl+,", Desc: "open command palette"},
			},
		},
		{
			Title: "Global",
			Bindings: []KeyBinding{
				{Key: "?", Desc: "toggle this help"},
				{Key: "ctrl+c", Desc: "force quit"},
			},
		},
	}, 0)
}

func (h HelpOverlay) View() string {
	var b strings.Builder

	title := styles.TitleStyle.Render("Keyboard Shortcuts")
	b.WriteString(title)
	b.WriteString("\n\n")

	for _, section := range h.Sections {
		b.WriteString(styles.MutedStyle.Render(section.Title))
		b.WriteString("\n")

		for _, binding := range section.Bindings {
			key := styles.KeyStyle.Render(binding.Key)
			desc := styles.KeyHintStyle.Render(binding.Desc)
			b.WriteString("  " + key + "  " + desc + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(styles.DimStyle.Render("Press ? or esc to close"))

	panelWidth := 44
	if h.Width > 0 && h.Width < panelWidth+8 {
		panelWidth = h.Width - 8
		if panelWidth < 30 {
			panelWidth = 30
		}
	}

	return styles.BoxStyle.
		Width(panelWidth).
		BorderForeground(styles.ColorPrimary).
		Render(strings.TrimRight(b.String(), "\n"))
}
