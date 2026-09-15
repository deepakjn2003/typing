package screens

import (
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/tui/components"
	"github.com/deepakjn2003/typing/internal/tui/layout"
	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type InputMode int

const (
	InputModeSelect InputMode = iota
	InputModePaste
	InputModeURL
)

type InputModel struct {
	mode          InputMode
	modeMenu      components.Menu
	textarea      textarea.Model
	textinput     textinput.Model
	width, height int
	err           string
}

func nativeTextInputStyles() textinput.Styles {
	s := textinput.DefaultStyles(styles.TerminalIsDark())
	s.Focused.Placeholder = lipgloss.NewStyle().Faint(true)
	s.Blurred.Placeholder = lipgloss.NewStyle().Faint(true)
	return s
}

func nativeTextareaStyles() textarea.Styles {
	s := textarea.DefaultStyles(styles.TerminalIsDark())
	for _, state := range []*textarea.StyleState{&s.Focused, &s.Blurred} {
		state.Placeholder = lipgloss.NewStyle().Faint(true)
	}
	return s
}

func NewInputModel() InputModel {
	menu := components.NewMenu([]components.MenuItem{
		{Label: "Paste text", Description: "Provide text directly", ID: "paste"},
		{Label: "Article URL", Description: "Fetch and summarize an article", ID: "url"},
	})

	ta := textarea.New()
	ta.Placeholder = "Paste your article, notes, or any text here..."
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.CharLimit = 50000
	ta.SetStyles(nativeTextareaStyles())

	ti := textinput.New()
	ti.Placeholder = "https://example.com/article"
	ti.CharLimit = 2048
	ti.Prompt = "  "
	ti.SetStyles(nativeTextInputStyles())

	return InputModel{
		mode:      InputModeSelect,
		modeMenu:  menu,
		textarea:  ta,
		textinput: ti,
	}
}

func (m *InputModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	contentW := layout.ContentWidth(w)
	m.modeMenu.SetWidth(contentW)

	taWidth := contentW - 6
	if taWidth < 20 {
		taWidth = 20
	}
	taHeight := h - 12
	if taHeight < 5 {
		taHeight = 5
	}
	m.textarea.SetWidth(taWidth)
	m.textarea.SetHeight(taHeight)
	m.textinput.CharLimit = 2048
}

func (m *InputModel) Reset() {
	m.mode = InputModeSelect
	m.err = ""
	m.modeMenu.Selected = 0
	m.textarea.Reset()
	m.textinput.SetValue("")
}

func (m InputModel) Mode() InputMode {
	return m.mode
}

func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	switch m.mode {
	case InputModeSelect:
		return m.updateSelect(msg)
	case InputModePaste:
		return m.updatePaste(msg)
	case InputModeURL:
		return m.updateURL(msg)
	}
	return m, nil
}

func (m InputModel) updateSelect(msg tea.Msg) (InputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case msg.Code == tea.KeyUp || msg.String() == "k":
			m.modeMenu.Up()
		case msg.Code == tea.KeyDown || msg.String() == "j":
			m.modeMenu.Down()
		case msg.Code == tea.KeyEnter:
			switch m.modeMenu.SelectedID() {
			case "paste":
				m.mode = InputModePaste
				return m, m.textarea.Focus()
			case "url":
				m.mode = InputModeURL
				return m, m.textinput.Focus()
			}
		}
	}
	return m, nil
}

func (m InputModel) updatePaste(msg tea.Msg) (InputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEscape {
			m.mode = InputModeSelect
			m.textarea.Blur()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m InputModel) updateURL(msg tea.Msg) (InputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEscape {
			m.mode = InputModeSelect
			m.textinput.Blur()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m InputModel) Value() string {
	switch m.mode {
	case InputModePaste:
		return m.textarea.Value()
	case InputModeURL:
		return m.textinput.Value()
	}
	return ""
}

func (m InputModel) SourceType() string {
	switch m.mode {
	case InputModePaste:
		return "paste"
	case InputModeURL:
		return "url"
	}
	return ""
}

func (m InputModel) HasContent() bool {
	return strings.TrimSpace(m.Value()) != ""
}

func (m InputModel) View() string {
	var b strings.Builder
	contentW := layout.ContentWidth(m.width)

	switch m.mode {
	case InputModeSelect:
		b.WriteString(styles.TitleStyle.Render("How do you want to practice?"))
		b.WriteString("\n\n")
		b.WriteString(m.modeMenu.View())

	case InputModePaste:
		b.WriteString(styles.TitleStyle.Render("Paste Content"))
		b.WriteString("\n\n")

		boxContent := m.textarea.View()
		box := styles.BoxStyle.Width(contentW).Render(boxContent)
		b.WriteString(box)

		if m.err != "" {
			b.WriteString("\n\n")
			b.WriteString(styles.ErrorStyle.Render(m.err))
		}

	case InputModeURL:
		b.WriteString(styles.TitleStyle.Render("Article URL"))
		b.WriteString("\n\n")

		boxContent := m.textinput.View()
		box := styles.BoxStyle.Width(contentW).Render(boxContent)
		b.WriteString(box)

		if m.err != "" {
			b.WriteString("\n\n")
			b.WriteString(styles.ErrorStyle.Render(m.err))
		}
	}

	return b.String()
}
