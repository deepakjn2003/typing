package components

import (
	"strings"

	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type MenuItem struct {
	Label       string
	Description string
	ID          string
}

type Menu struct {
	Items    []MenuItem
	Selected int
	Width    int
}

func NewMenu(items []MenuItem) Menu {
	return Menu{
		Items: items,
	}
}

func (m *Menu) SetWidth(w int) {
	m.Width = w
}

func (m *Menu) Up() {
	if m.Selected > 0 {
		m.Selected--
	}
}

func (m *Menu) Down() {
	if m.Selected < len(m.Items)-1 {
		m.Selected++
	}
}

func (m *Menu) SelectedItem() MenuItem {
	if m.Selected >= 0 && m.Selected < len(m.Items) {
		return m.Items[m.Selected]
	}
	return MenuItem{}
}

func (m *Menu) SelectedID() string {
	return m.SelectedItem().ID
}

func (m Menu) View() string {
	var lines []string
	for i, item := range m.Items {
		var labelLine string
		if i == m.Selected {
			labelLine = styles.MenuItemSelectedStyle.Render("▸ " + item.Label)
		} else {
			labelLine = styles.MenuItemStyle.Render("  " + item.Label)
		}
		lines = append(lines, labelLine)

		if item.Description != "" {
			descStyle := styles.MutedStyle
			if i == m.Selected {
				descStyle = styles.DimStyle
			}
			lines = append(lines, descStyle.Render("  "+item.Description))
		}

		if i < len(m.Items)-1 {
			lines = append(lines, "")
		}
	}
	return strings.Join(lines, "\n")
}
