package components

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/gum/v2/confirm"
	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type ConfirmResultMsg struct {
	Confirmed bool
}

// ConfirmDialog wraps gum confirm for embedded use with native terminal styling.
type ConfirmDialog struct {
	opts     confirm.Options
	selected bool
	done     bool
}

func NewQuitConfirmDialog() ConfirmDialog {
	opts := confirm.Options{
		Prompt:      "Quit typing?",
		Affirmative: "Yes",
		Negative:    "No",
		Default:     false,
		ShowHelp:    false,
	}
	opts.PromptStyle.Bold = true
	opts.SelectedStyle.Bold = true
	opts.SelectedStyle.Padding = "0 3"
	opts.UnselectedStyle.Padding = "0 3"

	return ConfirmDialog{
		opts:     opts,
		selected: false,
	}
}

func (c ConfirmDialog) Init() tea.Cmd {
	return nil
}

func (c ConfirmDialog) Update(msg tea.Msg) (ConfirmDialog, tea.Cmd) {
	if c.done {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.Key()
		switch {
		case key.Code == tea.KeyEscape:
			c.done = true
			return c, func() tea.Msg { return ConfirmResultMsg{Confirmed: false} }
		case key.Text == "y" || key.Text == "Y":
			c.done = true
			return c, func() tea.Msg { return ConfirmResultMsg{Confirmed: true} }
		case key.Text == "n" || key.Text == "N":
			c.done = true
			return c, func() tea.Msg { return ConfirmResultMsg{Confirmed: false} }
		case key.Code == tea.KeyLeft || key.Text == "h":
			c.selected = false
		case key.Code == tea.KeyRight || key.Text == "l":
			c.selected = true
		case key.Code == tea.KeyTab:
			c.selected = !c.selected
		case key.Code == tea.KeyEnter:
			c.done = true
			return c, func() tea.Msg { return ConfirmResultMsg{Confirmed: c.selected} }
		}
	}

	return c, nil
}

func (c ConfirmDialog) View() string {
	promptStyle := c.opts.PromptStyle.ToLipgloss()
	selectedStyle := c.opts.SelectedStyle.ToLipgloss().MarginRight(1).Reverse(true).Bold(true)
	unselectedStyle := c.opts.UnselectedStyle.ToLipgloss().MarginRight(1)

	var aff, neg string
	if c.selected {
		aff = selectedStyle.Render(c.opts.Affirmative)
		neg = unselectedStyle.Render(c.opts.Negative)
	} else {
		aff = unselectedStyle.Render(c.opts.Affirmative)
		neg = selectedStyle.Render(c.opts.Negative)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		promptStyle.Render(c.opts.Prompt),
		lipgloss.JoinHorizontal(lipgloss.Left, aff, neg),
	)

	return styles.BoxStyle.Render(content)
}

func (c ConfirmDialog) Done() bool {
	return c.done
}
