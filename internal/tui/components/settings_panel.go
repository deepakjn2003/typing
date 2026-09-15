package components

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/config"
	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type SettingsClosedMsg struct{}

type settingItem struct {
	label   string
	group   string
	cycle   func(*config.Config) string
	visible func(*config.Config) bool
}

type SettingsPanel struct {
	cfg      *config.Config
	items    []settingItem
	selected int
	width    int
	closing  bool
}

var (
	textSizeOptions  = []string{config.TextSizeSmall, config.TextSizeMedium, config.TextSizeLarge}
	visibleLineOpts  = []int{2, 3, 4, 5}
	maxWordOpts      = []int{50, 100, 150, 250, 300}
	practiceModeOpts = []string{config.PracticeModeTimed, config.PracticeModePassage}
	maxTimeOpts      = []int{15, 30, 60, 120, 180}
)

func NewSettingsPanel(cfg *config.Config) SettingsPanel {
	sp := SettingsPanel{cfg: cfg}
	sp.items = []settingItem{
		{
			label: "Text size",
			group: "Display",
			cycle: func(c *config.Config) string {
				idx := indexOf(textSizeOptions, c.TextSize)
				c.TextSize = textSizeOptions[(idx+1)%len(textSizeOptions)]
				return styles.ParseTextSize(c.TextSize).Label()
			},
		},
		{
			label: "Visible lines",
			group: "Display",
			cycle: func(c *config.Config) string {
				idx := indexOfInt(visibleLineOpts, c.VisibleLines)
				c.VisibleLines = visibleLineOpts[(idx+1)%len(visibleLineOpts)]
				return fmt.Sprintf("%d", c.VisibleLines)
			},
		},
		{
			label: "Max words",
			group: "Practice",
			cycle: func(c *config.Config) string {
				idx := indexOfInt(maxWordOpts, c.MaxWords)
				c.MaxWords = maxWordOpts[(idx+1)%len(maxWordOpts)]
				if c.MinWords > c.MaxWords {
					c.MinWords = c.MaxWords / 2
				}
				return fmt.Sprintf("%d", c.MaxWords)
			},
		},
		{
			label: "Practice mode",
			group: "Practice",
			cycle: func(c *config.Config) string {
				idx := indexOf(practiceModeOpts, c.PracticeMode)
				c.PracticeMode = practiceModeOpts[(idx+1)%len(practiceModeOpts)]
				return practiceModeLabel(c.PracticeMode)
			},
		},
		{
			label: "Max time",
			group: "Practice",
			cycle: func(c *config.Config) string {
				idx := indexOfInt(maxTimeOpts, c.MaxPracticeSeconds)
				c.MaxPracticeSeconds = maxTimeOpts[(idx+1)%len(maxTimeOpts)]
				return fmt.Sprintf("%ds", c.MaxPracticeSeconds)
			},
			visible: func(c *config.Config) bool { return c.IsTimedMode() },
		},
	}
	return sp
}

func practiceModeLabel(mode string) string {
	if mode == config.PracticeModeTimed {
		return "Timed"
	}
	return "Full passage"
}

func indexOf(slice []string, val string) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return 0
}

func indexOfInt(slice []int, val int) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return 0
}

func (s *SettingsPanel) visibleItems() []int {
	var indices []int
	for i, item := range s.items {
		if item.visible == nil || item.visible(s.cfg) {
			indices = append(indices, i)
		}
	}
	return indices
}

func (s *SettingsPanel) SetSize(w int) {
	s.width = w
}

func (s SettingsPanel) currentValue(item settingItem) string {
	switch item.label {
	case "Text size":
		return styles.ParseTextSize(s.cfg.TextSize).Label()
	case "Visible lines":
		return fmt.Sprintf("%d", s.cfg.VisibleLines)
	case "Max words":
		return fmt.Sprintf("%d", s.cfg.MaxWords)
	case "Practice mode":
		return practiceModeLabel(s.cfg.PracticeMode)
	case "Max time":
		return fmt.Sprintf("%ds", s.cfg.MaxPracticeSeconds)
	default:
		return ""
	}
}

func (s SettingsPanel) cycleSelected() {
	visible := s.visibleItems()
	if len(visible) == 0 {
		return
	}
	item := s.items[visible[s.selected]]
	item.cycle(s.cfg)
	s.cfg.Validate()
}

func (s SettingsPanel) Update(msg tea.Msg) (SettingsPanel, tea.Cmd) {
	if s.closing {
		return s, nil
	}

	visible := s.visibleItems()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.Key()
		switch {
		case key.Code == tea.KeyEscape:
			s.closing = true
			_ = s.cfg.Save()
			return s, func() tea.Msg { return SettingsClosedMsg{} }
		case key.Code == tea.KeyUp || msg.String() == "k":
			if s.selected > 0 {
				s.selected--
			}
			return s, nil
		case key.Code == tea.KeyDown || msg.String() == "j":
			if s.selected < len(visible)-1 {
				s.selected++
			}
			return s, nil
		case key.Code == tea.KeyEnter, key.Code == tea.KeySpace:
			s.cycleSelected()
			return s, nil
		}
	}

	return s, nil
}

func (s SettingsPanel) IsClosing() bool {
	return s.closing
}

func (s SettingsPanel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Settings"))
	b.WriteString("\n\n")

	visible := s.visibleItems()
	lastGroup := ""
	labelColWidth := 18

	for vi, idx := range visible {
		item := s.items[idx]
		if item.group != lastGroup {
			if lastGroup != "" {
				b.WriteString("\n")
			}
			b.WriteString(styles.MutedStyle.Render(item.group))
			b.WriteString("\n")
			b.WriteString(styles.Separator(labelColWidth + 20))
			b.WriteString("\n")
			lastGroup = item.group
		}

		value := s.currentValue(item)
		isSelected := vi == s.selected

		var prefix, indicator string
		if isSelected {
			prefix = "▸ "
			indicator = styles.MenuItemSelectedStyle.Render("▸")
		} else {
			prefix = "  "
			indicator = " "
		}

		label := item.label
		if isSelected {
			label = styles.MenuItemSelectedStyle.Render(prefix + label)
		} else {
			label = styles.MenuItemStyle.Render(prefix + label)
		}

		valueStyled := styles.InfoStyle.Render(value)
		if isSelected {
			valueStyled = styles.MenuItemSelectedStyle.Render(value)
		}

		labelPlain := prefix + item.label
		gap := labelColWidth - len(labelPlain)
		if gap < 2 {
			gap = 2
		}

		line := label + strings.Repeat(" ", gap) + valueStyled + "    " + indicator
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styles.DimStyle.Render("enter/space cycle  •  esc save & close"))

	panelWidth := 52
	if s.width > 0 && s.width < panelWidth+8 {
		panelWidth = s.width - 8
		if panelWidth < 36 {
			panelWidth = 36
		}
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Padding(1, 2).
		Width(panelWidth).
		Render(strings.TrimRight(b.String(), "\n"))
}
