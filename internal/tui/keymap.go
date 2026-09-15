package tui

import (
	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/tui/components"
	"github.com/deepakjn2003/typing/internal/tui/screens"
)

type ScreenContext struct {
	InputMode     screens.InputMode
	TypingTitle   string
	TypingTimer   string
	ProcessingErr bool
	TimedMode     bool
}

func globalBindings() []components.KeyBinding {
	return []components.KeyBinding{
		{Key: ":", Desc: "commands"},
		{Key: "ctrl+,", Desc: "commands"},
		{Key: "?", Desc: "help"},
	}
}

func ScreenMeta(screen Screen, ctx ScreenContext) (bindings []components.KeyBinding, align lipgloss.Position) {
	align = lipgloss.Top

	switch screen {
	case ScreenHome:
		align = lipgloss.Center
		bindings = append([]components.KeyBinding{
			{Key: "↑/↓", Desc: "navigate"},
			{Key: "enter", Desc: "select"},
			{Key: "q", Desc: "quit"},
		}, globalBindings()...)

	case ScreenInput:
		switch ctx.InputMode {
		case screens.InputModeSelect:
			align = lipgloss.Center
			bindings = append([]components.KeyBinding{
				{Key: "↑/↓", Desc: "navigate"},
				{Key: "enter", Desc: "select"},
				{Key: "esc", Desc: "back"},
			}, globalBindings()...)
		case screens.InputModePaste:
			bindings = append([]components.KeyBinding{
				{Key: "ctrl+d", Desc: "submit"},
				{Key: "esc", Desc: "back"},
			}, globalBindings()...)
		case screens.InputModeURL:
			bindings = append([]components.KeyBinding{
				{Key: "enter", Desc: "submit"},
				{Key: "esc", Desc: "back"},
			}, globalBindings()...)
		}

	case ScreenProcessing:
		align = lipgloss.Center
		if ctx.ProcessingErr {
			bindings = append([]components.KeyBinding{
				{Key: "esc", Desc: "back"},
			}, globalBindings()...)
		} else {
			bindings = globalBindings()
		}

	case ScreenTyping:
		align = lipgloss.Center
		bindings = append([]components.KeyBinding{
			{Key: "esc", Desc: "pause"},
			{Key: "backspace", Desc: "correct"},
		}, globalBindings()...)
		if ctx.TimedMode {
			bindings = append(bindings, components.KeyBinding{Key: "timer", Desc: "countdown active"})
		}

	case ScreenResults:
		align = lipgloss.Center
		bindings = append([]components.KeyBinding{
			{Key: "↑/↓", Desc: "navigate"},
			{Key: "enter", Desc: "select"},
		}, globalBindings()...)

	case ScreenHistory:
		bindings = append([]components.KeyBinding{
			{Key: "↑/↓", Desc: "navigate"},
			{Key: "esc", Desc: "back"},
		}, globalBindings()...)
	}

	return bindings, align
}

func HelpSections(screen Screen, ctx ScreenContext) []components.HelpSection {
	sections := []components.HelpSection{
		{
			Title: "Navigation",
			Bindings: []components.KeyBinding{
				{Key: "↑/↓", Desc: "move selection"},
				{Key: "enter", Desc: "confirm"},
				{Key: "esc", Desc: "go back"},
			},
		},
	}

	switch screen {
	case ScreenTyping:
		sections = append(sections, components.HelpSection{
			Title: "Typing",
			Bindings: []components.KeyBinding{
				{Key: "esc", Desc: "pause / resume"},
				{Key: "backspace", Desc: "correct mistake"},
			},
		})
	case ScreenInput:
		switch ctx.InputMode {
		case screens.InputModePaste:
			sections = append(sections, components.HelpSection{
				Title: "Input",
				Bindings: []components.KeyBinding{
					{Key: "ctrl+d", Desc: "submit text"},
				},
			})
		case screens.InputModeURL:
			sections = append(sections, components.HelpSection{
				Title: "Input",
				Bindings: []components.KeyBinding{
					{Key: "enter", Desc: "fetch article"},
				},
			})
		}
	case ScreenHome:
		sections = append(sections, components.HelpSection{
			Title: "Home",
			Bindings: []components.KeyBinding{
				{Key: "q", Desc: "quit"},
			},
		})
	}

	sections = append(sections,
		components.HelpSection{
			Title: "Commands",
			Bindings: []components.KeyBinding{
				{Key: ":", Desc: "open command palette"},
				{Key: "ctrl+,", Desc: "open command palette"},
			},
		},
		components.HelpSection{
			Title: "Global",
			Bindings: []components.KeyBinding{
				{Key: "?", Desc: "toggle this help"},
				{Key: "ctrl+c", Desc: "force quit"},
			},
		},
	)

	return sections
}
