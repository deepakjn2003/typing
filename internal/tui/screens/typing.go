package screens

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/config"
	"github.com/deepakjn2003/typing/internal/tui/components"
	"github.com/deepakjn2003/typing/internal/tui/layout"
	"github.com/deepakjn2003/typing/internal/tui/styles"
	"github.com/deepakjn2003/typing/internal/typing"
)

type TypingTickMsg time.Time

type TypingModel struct {
	engine        *typing.Engine
	title         string
	cfg           *config.Config
	textSize      styles.TextSizePreset
	visibleLines  int
	width, height int
}

func NewTypingModel(passage, title string, cfg *config.Config) TypingModel {
	engine := typing.NewEngine(passage)
	if cfg != nil && cfg.IsTimedMode() {
		engine.SetTimeLimit(time.Duration(cfg.MaxPracticeSeconds) * time.Second)
	}

	visibleLines := 3
	textSize := styles.TextSizeMedium
	if cfg != nil {
		visibleLines = cfg.VisibleLines
		textSize = styles.ParseTextSize(cfg.TextSize)
	}

	return TypingModel{
		engine:       engine,
		title:        title,
		cfg:          cfg,
		textSize:     textSize,
		visibleLines: visibleLines,
	}
}

func (m *TypingModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *TypingModel) Engine() *typing.Engine {
	return m.engine
}

func (m TypingModel) IsComplete() bool {
	return m.engine.IsComplete()
}

func (m TypingModel) IsTimeUp() bool {
	return m.engine.IsTimeUp()
}

func (m TypingModel) Update(msg tea.Msg) (TypingModel, tea.Cmd) {
	if m.engine.IsComplete() {
		return m, nil
	}

	switch msg := msg.(type) {
	case TypingTickMsg:
		if m.engine.State() == typing.StateRunning {
			if m.engine.CheckTimeLimit() {
				return m, nil
			}
			return m, tickTyping()
		}
		return m, nil

	case tea.KeyPressMsg:
		key := msg.Key()

		if key.Code == tea.KeyBackspace {
			m.engine.Backspace()
			return m, nil
		}

		if key.Code == tea.KeyEscape {
			if m.engine.State() == typing.StateRunning {
				m.engine.Pause()
			} else if m.engine.State() == typing.StatePaused {
				m.engine.Resume()
				return m, tickTyping()
			}
			return m, nil
		}

		if key.Text == "" {
			return m, nil
		}

		for _, r := range key.Text {
			m.engine.ProcessRune(r)
		}

		if m.engine.IsComplete() {
			return m, nil
		}

		if m.engine.State() == typing.StateRunning && m.engine.Position() == 1 {
			return m, tickTyping()
		}

		return m, nil
	}

	return m, nil
}

func tickTyping() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TypingTickMsg(t)
	})
}

func StartTypingTick() tea.Cmd {
	return tickTyping()
}

func (m TypingModel) Title() string {
	return m.title
}

func (m TypingModel) TimerLine() string {
	snap := m.engine.Snapshot()
	timer := formatDuration(snap.Elapsed)
	if m.cfg != nil && m.cfg.IsTimedMode() {
		remaining := time.Duration(m.cfg.MaxPracticeSeconds)*time.Second - snap.Elapsed
		if remaining < 0 {
			remaining = 0
		}
		timer = formatDuration(remaining)
	}
	return timer
}

func (m TypingModel) View() string {
	snap := m.engine.Snapshot()
	stats := typing.SessionStats(snap)

	lineWidth := m.textSize.CharsPerLine(layout.ContentWidth(m.width))

	styleChar := func(pos int, r rune) string {
		charState := m.engine.CharStateAt(pos)
		ch := string(r)
		var base lipgloss.Style
		switch charState {
		case typing.CharCorrect:
			base = styles.CorrectCharStyle
		case typing.CharIncorrect:
			base = styles.IncorrectCharStyle
		case typing.CharCursor:
			base = styles.CursorCharStyle
		default:
			base = styles.UntypedCharStyle
		}
		styled := m.textSize.StyleChar(base)
		padding := m.textSize.CharPadding()
		if padding > 0 {
			styled = styled.Padding(0, padding)
		}
		return styled.Render(ch)
	}

	passageView := BuildTypingWindow(
		m.engine.TargetRunes(),
		m.engine.Position(),
		lineWidth,
		m.visibleLines,
		m.textSize.LineSpacing(),
		styleChar,
	)

	passageBox := lipgloss.NewStyle().
		Padding(1, 3).
		Render(passageView)

	var sections []string

	statGroup := components.StatGroup{
		Stats: []components.Stat{
			{Label: "WPM", Value: fmt.Sprintf("%.0f", stats.WPM), Width: 4},
			{Label: "Accuracy", Value: fmt.Sprintf("%.1f%%", stats.Accuracy), Width: 6},
			{Label: "Errors", Value: fmt.Sprintf("%d", stats.TotalErrors), Width: 3},
			{Label: "Time", Value: m.TimerLine(), Width: 5},
		},
	}
	sections = append(sections, statGroup.View())
	sections = append(sections, "")
	sections = append(sections, passageBox)

	if m.engine.State() == typing.StatePaused {
		sections = append(sections, "")
		pausedStyle := styles.WarningStyle.Bold(true)
		if m.textSize == styles.TextSizeLarge {
			pausedStyle = pausedStyle.Padding(1, 0)
		}
		sections = append(sections, pausedStyle.Render("PAUSED — press Esc to resume"))
	}

	return strings.Join(sections, "\n")
}

func formatDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
