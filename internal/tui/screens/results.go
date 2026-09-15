package screens

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/deepakjn2003/typing/internal/tui/components"
	"github.com/deepakjn2003/typing/internal/tui/styles"
	"github.com/deepakjn2003/typing/internal/typing"
)

type ResultsModel struct {
	stats         typing.Stats
	title         string
	menu          components.Menu
	width, height int
	timeUp        bool
}

func NewResultsModel(stats typing.Stats, title string, timeUp bool) ResultsModel {
	menu := components.NewMenu([]components.MenuItem{
		{Label: "Practice Again", ID: "again"},
		{Label: "Home", ID: "home"},
	})

	return ResultsModel{
		stats:  stats,
		title:  title,
		menu:   menu,
		timeUp: timeUp,
	}
}

func (m *ResultsModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.menu.SetWidth(w)
}

func (m ResultsModel) Update(msg tea.Msg) (ResultsModel, tea.Cmd) {
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

func (m ResultsModel) SelectedAction() string {
	return m.menu.SelectedID()
}

func (m ResultsModel) View() string {
	stats := m.stats
	var sections []string

	header := "SESSION COMPLETE"
	if m.timeUp || stats.TimeUp {
		header = "TIME'S UP"
	}

	sections = append(sections, styles.SuccessStyle.Bold(true).Render(header))
	sections = append(sections, "")
	sections = append(sections, styles.BigStatStyle.Render(fmt.Sprintf("%.0f WPM", stats.WPM)))
	sections = append(sections, "")

	statGroup := components.StatGroup{
		Stats: []components.Stat{
			{Label: "Accuracy", Value: fmt.Sprintf("%.1f%%", stats.Accuracy), Width: 6},
			{Label: "Errors", Value: fmt.Sprintf("%d", stats.TotalErrors), Width: 3},
			{Label: "Time", Value: formatDuration(stats.Elapsed), Width: 5},
		},
	}
	sections = append(sections, statGroup.View())
	sections = append(sections, "")

	sections = append(sections, styles.MutedStyle.Render("Performance"))
	sections = append(sections, styles.Separator(40))

	var details strings.Builder
	addStat := func(label, value string) {
		l := styles.StatLabelStyle.Render(fmt.Sprintf("  %-20s", label))
		v := styles.StatValueStyle.Render(value)
		details.WriteString(l + v + "\n")
	}

	addStat("Characters", fmt.Sprintf("%d", stats.Characters))
	addStat("Correct", fmt.Sprintf("%d", stats.CorrectChars))
	addStat("Incorrect", fmt.Sprintf("%d", stats.IncorrectChars))
	addStat("Corrections", fmt.Sprintf("%d", stats.Corrections))

	sections = append(sections, strings.TrimRight(details.String(), "\n"))

	if len(stats.ErrorAnalysis.TopErrors) > 0 {
		sections = append(sections, "")
		sections = append(sections, styles.MutedStyle.Render("Common mistakes"))
		limit := 5
		if len(stats.ErrorAnalysis.TopErrors) < limit {
			limit = len(stats.ErrorAnalysis.TopErrors)
		}
		for _, ef := range stats.ErrorAnalysis.TopErrors[:limit] {
			sections = append(sections, styles.ErrorStyle.Render(
				fmt.Sprintf("  %c → %c    %d times", ef.Pair.Expected, ef.Pair.Actual, ef.Count),
			))
		}
	}

	sections = append(sections, "")
	sections = append(sections, m.menu.View())

	return strings.Join(sections, "\n")
}
