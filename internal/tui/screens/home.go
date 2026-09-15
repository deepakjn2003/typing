package screens

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/deepakjn2003/typing/internal/tui/components"
	"github.com/deepakjn2003/typing/internal/tui/layout"
	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type HomeModel struct {
	menu          components.Menu
	width, height int
}

func NewHomeModel() HomeModel {
	menu := components.NewMenu([]components.MenuItem{
		{Label: "Start Practice", ID: "practice"},
		{Label: "History", ID: "history"},
	})

	return HomeModel{menu: menu}
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.menu.SetWidth(layout.ContentWidth(w))
}

func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case msg.Code == tea.KeyUp || msg.String() == "k":
			m.menu.Up()
		case msg.Code == tea.KeyDown || msg.String() == "j":
			m.menu.Down()
		}
	}
	return m, nil
}

func (m HomeModel) SelectedAction() string {
	return m.menu.SelectedID()
}

func (m HomeModel) View() string {
	var sections []string

	sections = append(sections, styles.TitleStyle.Render("typing"))
	sections = append(sections, styles.SubtitleStyle.Render("Learn something. Type something."))
	sections = append(sections, "")
	sections = append(sections, m.menu.View())

	return strings.Join(sections, "\n")
}
