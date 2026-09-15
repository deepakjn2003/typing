package screens

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/storage"
	"github.com/deepakjn2003/typing/internal/tui/layout"
	"github.com/deepakjn2003/typing/internal/tui/styles"
	"github.com/deepakjn2003/typing/internal/typing"
)

type HistoryLoadedMsg struct {
	Sessions []typing.Session
	Err      error
}

type HistoryModel struct {
	sessions      []typing.Session
	selected      int
	width, height int
	repo          storage.SessionRepository
	loaded        bool
	err           string
}

const (
	historyPrefixW = 2
	historyWPMW    = 9
	historyAccW    = 8
	historyDurW    = 8
)

func NewHistoryModel(repo storage.SessionRepository) HistoryModel {
	return HistoryModel{
		repo: repo,
	}
}

func (m HistoryModel) Init() tea.Cmd {
	return m.loadHistory()
}

func (m HistoryModel) loadHistory() tea.Cmd {
	repo := m.repo
	return func() tea.Msg {
		if repo == nil {
			return HistoryLoadedMsg{}
		}
		sessions, err := repo.List(context.Background(), 50, 0)
		return HistoryLoadedMsg{Sessions: sessions, Err: err}
	}
}

func (m *HistoryModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m HistoryModel) Update(msg tea.Msg) (HistoryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case HistoryLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err.Error()
		} else {
			m.sessions = msg.Sessions
		}
		m.loaded = true

	case tea.KeyPressMsg:
		switch {
		case msg.Code == tea.KeyUp || msg.String() == "k":
			if m.selected > 0 {
				m.selected--
			}
		case msg.Code == tea.KeyDown || msg.String() == "j":
			if m.selected < len(m.sessions)-1 {
				m.selected++
			}
		}
	}
	return m, nil
}

func (m HistoryModel) contentWidth() int {
	return layout.ContentWidth(m.width)
}

func (m HistoryModel) titleWidth() int {
	titleW := m.contentWidth() - historyPrefixW - historyWPMW - historyAccW - historyDurW
	if titleW < 16 {
		return 16
	}
	return titleW
}

func (m HistoryModel) renderHeaderRow() string {
	titleW := m.titleWidth()
	titleCol := styles.MutedStyle.Render(fmt.Sprintf("%s%-*s", strings.Repeat(" ", historyPrefixW), titleW, "Session"))
	wpmCol := styles.MutedStyle.Render(fmt.Sprintf("%*s", historyWPMW, "WPM"))
	accCol := styles.MutedStyle.Render(fmt.Sprintf("%*s", historyAccW, "Acc"))
	durCol := styles.MutedStyle.Render(fmt.Sprintf("%*s", historyDurW, "Time"))
	return titleCol + wpmCol + accCol + durCol
}

func historySelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(styles.ColorPrimary).Bold(true)
}

func (m HistoryModel) renderRow(prefix, title string, wpm, acc, dur string, selected bool) string {
	titleW := m.titleWidth()
	title = truncateWidth(title, titleW)

	titleCell := fmt.Sprintf("%s%-*s", prefix, titleW, title)
	wpmCell := fmt.Sprintf("%*s", historyWPMW, wpm)
	accCell := fmt.Sprintf("%*s", historyAccW, acc)
	durCell := fmt.Sprintf("%*s", historyDurW, dur)

	if selected {
		s := historySelectedStyle()
		return s.Render(titleCell) + s.Render(wpmCell) + s.Render(accCell) + s.Render(durCell)
	}

	return styles.BodyStyle.Render(titleCell) +
		styles.InfoStyle.Render(wpmCell) +
		styles.SuccessStyle.Render(accCell) +
		styles.MutedStyle.Render(durCell)
}

func (m HistoryModel) View() string {
	var b strings.Builder
	contentW := m.contentWidth()

	b.WriteString(lipgloss.NewStyle().Width(contentW).Align(lipgloss.Center).Render(styles.TitleStyle.Render("History")))
	b.WriteString("\n\n")

	if m.err != "" {
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err))
		b.WriteString("\n\n")
	}

	if !m.loaded {
		b.WriteString(styles.MutedStyle.Render("Loading..."))
		b.WriteString("\n")
	} else if len(m.sessions) == 0 {
		b.WriteString(styles.MutedStyle.Render("No sessions yet. Start practicing!"))
		b.WriteString("\n")
	} else {
		groups := groupByDate(m.sessions)
		sessionIdx := 0

		for _, group := range groups {
			b.WriteString(styles.MutedStyle.Render(group.Label))
			b.WriteString("\n")
			b.WriteString(styles.Separator(contentW))
			b.WriteString("\n")
			b.WriteString(m.renderHeaderRow())
			b.WriteString("\n")

			for _, s := range group.Sessions {
				prefix := "  "
				if sessionIdx == m.selected {
					prefix = "▸ "
				}

				line := m.renderRow(
					prefix,
					s.Title,
					fmt.Sprintf("%.0f WPM", s.WPM),
					fmt.Sprintf("%.1f%%", s.Accuracy),
					formatDuration(s.Duration),
					sessionIdx == m.selected,
				)

				b.WriteString(line)
				b.WriteString("\n")
				sessionIdx++
			}

			b.WriteString("\n")
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

func truncateWidth(s string, max int) string {
	if max < 4 {
		return s
	}
	if lipgloss.Width(s) <= max {
		return s
	}
	for len(s) > 0 && lipgloss.Width(s) > max-3 {
		_, size := utf8.DecodeLastRuneInString(s)
		if size == 0 {
			break
		}
		s = s[:len(s)-size]
	}
	return s + "..."
}

type dateGroup struct {
	Label    string
	Sessions []typing.Session
}

func groupByDate(sessions []typing.Session) []dateGroup {
	if len(sessions) == 0 {
		return nil
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	groups := make(map[string]*dateGroup)
	var order []string

	for _, s := range sessions {
		sessionDate := time.Date(s.CreatedAt.Year(), s.CreatedAt.Month(), s.CreatedAt.Day(),
			0, 0, 0, 0, s.CreatedAt.Location())

		var label string
		switch {
		case sessionDate.Equal(today):
			label = "Today"
		case sessionDate.Equal(yesterday):
			label = "Yesterday"
		default:
			label = s.CreatedAt.Format("January 2, 2006")
		}

		g, exists := groups[label]
		if !exists {
			g = &dateGroup{Label: label}
			groups[label] = g
			order = append(order, label)
		}
		g.Sessions = append(g.Sessions, s)
	}

	result := make([]dateGroup, 0, len(order))
	for _, label := range order {
		result = append(result, *groups[label])
	}
	return result
}
