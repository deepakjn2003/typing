package components

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type PaletteCommand int

const (
	PaletteCmdSettings PaletteCommand = iota
	PaletteCmdHelp
	PaletteCmdQuit
	PaletteCmdRestart
)

type PaletteResultMsg struct {
	Command PaletteCommand
}

type PaletteContext struct {
	ShowRestart bool
}

type CommandPalette struct {
	input    textinput.Model
	commands []paletteEntry
	filtered []int
	selected int
	width    int
	height   int
	quitting bool
}

type paletteEntry struct {
	label    string
	shortcut string
	group    string
	command  PaletteCommand
}

func paletteTextInputStyles() textinput.Styles {
	s := textinput.DefaultStyles(styles.TerminalIsDark())
	s.Focused.Placeholder = lipgloss.NewStyle().Faint(true)
	s.Blurred.Placeholder = lipgloss.NewStyle().Faint(true)
	return s
}

func NewCommandPalette() CommandPalette {
	return NewCommandPaletteWithContext(PaletteContext{})
}

func NewCommandPaletteWithContext(ctx PaletteContext) CommandPalette {
	ti := textinput.New()
	ti.Placeholder = "Search commands..."
	ti.Prompt = ""
	ti.CharLimit = 64
	ti.SetStyles(paletteTextInputStyles())

	var commands []paletteEntry
	if ctx.ShowRestart {
		commands = append(commands, paletteEntry{
			label:    "Restart Session",
			shortcut: "r",
			group:    "Session",
			command:  PaletteCmdRestart,
		})
	}
	commands = append(commands,
		paletteEntry{label: "Settings", shortcut: "ctrl+,", group: "General", command: PaletteCmdSettings},
		paletteEntry{label: "Help", shortcut: "?", group: "General", command: PaletteCmdHelp},
		paletteEntry{label: "Quit", shortcut: "q", group: "General", command: PaletteCmdQuit},
	)

	cp := CommandPalette{
		input:    ti,
		commands: commands,
	}
	cp.refilter()
	return cp
}

func (c *CommandPalette) Init() tea.Cmd {
	return c.input.Focus()
}

func (c *CommandPalette) SetSize(w, h int) {
	c.width = w
	c.height = h
	inputWidth := paletteWidth(w) - 8
	if inputWidth < 20 {
		inputWidth = 20
	}
	c.input.SetWidth(inputWidth)
}

func paletteWidth(termWidth int) int {
	w := int(float64(termWidth) * 0.6)
	if w > 60 {
		w = 60
	}
	if w < 30 {
		w = 30
	}
	return w
}

func (c *CommandPalette) refilter() {
	query := strings.ToLower(strings.TrimSpace(c.input.Value()))
	c.filtered = nil
	for i, cmd := range c.commands {
		if query == "" || strings.Contains(strings.ToLower(cmd.label), query) {
			c.filtered = append(c.filtered, i)
		}
	}
	if c.selected >= len(c.filtered) {
		c.selected = len(c.filtered) - 1
	}
	if c.selected < 0 {
		c.selected = 0
	}
}

func (c CommandPalette) Update(msg tea.Msg) (CommandPalette, tea.Cmd) {
	if c.quitting {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.Key()
		switch {
		case key.Code == tea.KeyEscape:
			c.quitting = true
			return c, nil
		case key.Code == tea.KeyUp || msg.String() == "k":
			if c.selected > 0 {
				c.selected--
			}
			return c, nil
		case key.Code == tea.KeyDown || msg.String() == "j":
			if c.selected < len(c.filtered)-1 {
				c.selected++
			}
			return c, nil
		case key.Code == tea.KeyEnter:
			if len(c.filtered) > 0 {
				cmd := c.commands[c.filtered[c.selected]].command
				c.quitting = true
				return c, func() tea.Msg { return PaletteResultMsg{Command: cmd} }
			}
			return c, nil
		}
	}

	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)
	c.refilter()
	return c, cmd
}

func (c CommandPalette) InputValue() string {
	return c.input.Value()
}

func (c CommandPalette) IsQuitting() bool {
	return c.quitting
}

func (c CommandPalette) View() string {
	var b strings.Builder

	title := styles.TitleStyle.Render("Commands")
	escHint := styles.MutedStyle.Render("esc")
	b.WriteString(title)
	b.WriteString(strings.Repeat(" ", 4))
	b.WriteString(escHint)
	b.WriteString("\n\n")

	b.WriteString(c.input.View())
	b.WriteString("\n\n")

	pw := paletteWidth(c.width)
	lastGroup := ""
	for fi, idx := range c.filtered {
		cmd := c.commands[idx]
		if cmd.group != lastGroup {
			if lastGroup != "" {
				b.WriteString("\n")
			}
			b.WriteString(styles.MutedStyle.Render(cmd.group))
			b.WriteString("\n")
			lastGroup = cmd.group
		}

		label := cmd.label
		if fi == c.selected {
			label = styles.MenuItemSelectedStyle.Render("▸ " + label)
		} else {
			label = styles.MenuItemStyle.Render("  " + label)
		}

		shortcut := ""
		if cmd.shortcut != "" {
			shortcut = styles.KeyStyle.Render(cmd.shortcut)
		}

		labelWidth := lipgloss.Width(label)
		gap := pw - labelWidth - lipgloss.Width(shortcut) - 4
		if gap < 1 {
			gap = 1
		}

		line := label + strings.Repeat(" ", gap) + shortcut
		b.WriteString(line)
		b.WriteString("\n")
	}

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Padding(1, 2).
		Width(pw).
		Render(strings.TrimRight(b.String(), "\n"))

	return panel
}

