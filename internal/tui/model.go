package tui

import (
	"context"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/deepakjn2003/typing/internal/config"
	"github.com/deepakjn2003/typing/internal/content"
	"github.com/deepakjn2003/typing/internal/storage"
	"github.com/deepakjn2003/typing/internal/tui/components"
	"github.com/deepakjn2003/typing/internal/tui/layout"
	"github.com/deepakjn2003/typing/internal/tui/screens"
	"github.com/deepakjn2003/typing/internal/typing"
	"github.com/google/uuid"
)

type Screen int

const (
	ScreenHome Screen = iota
	ScreenInput
	ScreenProcessing
	ScreenTyping
	ScreenResults
	ScreenHistory
)

type Model struct {
	screen Screen
	width  int
	height int

	showHelp     bool
	showQuit     bool
	showPalette  bool
	showSettings bool
	quitDlg      components.ConfirmDialog
	help         components.HelpOverlay
	palette      components.CommandPalette
	settings     components.SettingsPanel

	home       screens.HomeModel
	input      screens.InputModel
	processing screens.ProcessingModel
	typing     screens.TypingModel
	results    screens.ResultsModel
	history    screens.HistoryModel

	pipeline *content.Pipeline
	repo     storage.SessionRepository
	cfg      *config.Config

	currentSource     string
	currentSourceType string
	currentTitle      string
}

func NewModel(pipeline *content.Pipeline, repo storage.SessionRepository, cfg *config.Config) Model {
	return Model{
		screen:     ScreenHome,
		home:       screens.NewHomeModel(),
		input:      screens.NewInputModel(),
		processing: screens.NewProcessingModel(),
		palette:    components.NewCommandPalette(),
		cfg:        cfg,
		pipeline:   pipeline,
		repo:       repo,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.palette.SetSize(msg.Width, msg.Height)
		m.settings.SetSize(msg.Width)
		m.home.SetSize(msg.Width, msg.Height)
		m.input.SetSize(msg.Width, msg.Height)
		m.processing.SetSize(msg.Width, msg.Height)
		m.typing.SetSize(msg.Width, msg.Height)
		m.results.SetSize(msg.Width, msg.Height)
		m.history.SetSize(msg.Width, msg.Height)
		return m, nil

	case components.ConfirmResultMsg:
		m.showQuit = false
		if msg.Confirmed {
			return m, tea.Quit
		}
		return m, nil

	case components.PaletteResultMsg:
		m.showPalette = false
		switch msg.Command {
		case components.PaletteCmdSettings:
			m.showSettings = true
			m.settings = components.NewSettingsPanel(m.cfg)
			m.settings.SetSize(m.width)
			return m, nil
		case components.PaletteCmdHelp:
			m.showHelp = true
			m.help = components.NewHelpOverlay(HelpSections(m.screen, m.screenContext()), m.width)
			return m, nil
		case components.PaletteCmdQuit:
			m.showQuit = true
			m.quitDlg = components.NewQuitConfirmDialog()
			return m, m.quitDlg.Init()
		case components.PaletteCmdRestart:
			if m.screen == ScreenTyping || m.screen == ScreenResults {
				m.screen = ScreenInput
				m.input.Reset()
				return m, nil
			}
		}
		return m, nil

	case components.SettingsClosedMsg:
		m.showSettings = false
		return m, nil
	}

	if m.showQuit {
		return m.updateQuitOverlay(msg)
	}

	if m.showSettings {
		return m.updateSettingsOverlay(msg)
	}

	if m.showPalette {
		return m.updatePaletteOverlay(msg)
	}

	if m.showHelp {
		return m.updateHelpOverlay(msg)
	}

	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		if keyMsg.String() == "?" {
			m.showHelp = true
			m.help = components.NewHelpOverlay(HelpSections(m.screen, m.screenContext()), m.width)
			return m, nil
		}
		if isPaletteKey(keyMsg) {
			m.showPalette = true
			m.palette = components.NewCommandPaletteWithContext(m.paletteContext())
			m.palette.SetSize(m.width, m.height)
			return m, m.palette.Init()
		}
	}

	switch m.screen {
	case ScreenHome:
		return m.updateHome(msg)
	case ScreenInput:
		return m.updateInput(msg)
	case ScreenProcessing:
		return m.updateProcessing(msg)
	case ScreenTyping:
		return m.updateTyping(msg)
	case ScreenResults:
		return m.updateResults(msg)
	case ScreenHistory:
		return m.updateHistory(msg)
	}

	return m, nil
}

func isPaletteKey(msg tea.KeyPressMsg) bool {
	if msg.String() == "ctrl+," {
		return true
	}
	return msg.Key().Text == ":"
}

func (m Model) updatePaletteOverlay(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		if keyMsg.Key().Code == tea.KeyEscape && m.palette.InputValue() == "" {
			m.showPalette = false
			return m, nil
		}
		if keyMsg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.palette, cmd = m.palette.Update(msg)
	if m.palette.IsQuitting() && cmd == nil {
		m.showPalette = false
	}
	return m, cmd
}

func (m Model) updateSettingsOverlay(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		if keyMsg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.settings, cmd = m.settings.Update(msg)
	if m.settings.IsClosing() {
		m.showSettings = false
	}
	return m, cmd
}

func (m Model) updateQuitOverlay(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		if keyMsg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.quitDlg, cmd = m.quitDlg.Update(msg)
	return m, cmd
}

func (m Model) updateHelpOverlay(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case keyMsg.String() == "?" || keyMsg.Key().Code == tea.KeyEscape:
			m.showHelp = false
			return m, nil
		case keyMsg.String() == "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) updateHome(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		keyStr := msg.String()
		switch {
		case keyStr == "q":
			m.showQuit = true
			m.quitDlg = components.NewQuitConfirmDialog()
			return m, m.quitDlg.Init()
		case keyStr == "ctrl+c":
			return m, tea.Quit
		case msg.Key().Code == tea.KeyEnter:
			switch m.home.SelectedAction() {
			case "practice":
				m.screen = ScreenInput
				m.input.Reset()
				return m, nil
			case "history":
				m.screen = ScreenHistory
				m.history = screens.NewHistoryModel(m.repo)
				m.history.SetSize(m.width, m.height)
				return m, m.history.Init()
			}
		}
	}

	var cmd tea.Cmd
	m.home, cmd = m.home.Update(msg)
	return m, cmd
}

func (m Model) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		keyStr := msg.String()

		if keyStr == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.Key().Code == tea.KeyEscape && m.input.Mode() == screens.InputModeSelect {
			m.screen = ScreenHome
			m.home = screens.NewHomeModel()
			m.home.SetSize(m.width, m.height)
			return m, nil
		}

		if keyStr == "ctrl+d" && m.input.Mode() == screens.InputModePaste {
			if m.input.HasContent() {
				return m.startProcessing(m.input.Value(), m.input.SourceType())
			}
		}

		if msg.Key().Code == tea.KeyEnter && m.input.Mode() == screens.InputModeURL {
			if m.input.HasContent() {
				return m.startProcessing(m.input.Value(), m.input.SourceType())
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) summaryOptions() content.SummaryOptions {
	if m.cfg == nil {
		return content.DefaultSummaryOptions()
	}
	return content.SummaryOptions{
		MinWords: m.cfg.MinWords,
		MaxWords: m.cfg.MaxWords,
	}
}

func truncateToMaxWords(text string, maxWords int) string {
	words := strings.Fields(text)
	if len(words) <= maxWords {
		return text
	}
	return strings.Join(words[:maxWords], " ")
}

func (m Model) startProcessing(source, sourceType string) (tea.Model, tea.Cmd) {
	m.currentSource = source
	m.currentSourceType = sourceType
	m.screen = ScreenProcessing
	m.processing.Reset()

	pipeline := m.pipeline
	opts := m.summaryOptions()
	cfg := m.cfg
	repo := m.repo
	return m, tea.Batch(
		screens.StartProcessingTick(),
		func() tea.Msg {
			if pipeline == nil {
				passage := source
				if cfg != nil {
					passage = truncateToMaxWords(passage, cfg.MaxWords)
				}
				return screens.ProcessingCompleteMsg{
					Title:      "Practice",
					Passage:    passage,
					Source:     source,
					SourceType: sourceType,
				}
			}

			// Query weak characters from recent session history
			if repo != nil {
				if weakChars, err := repo.WeakCharacters(context.Background(), 10); err == nil && len(weakChars) > 0 {
					chars := make([]string, len(weakChars))
					for i, wc := range weakChars {
						chars[i] = string(wc.Char)
					}
					opts.WeakCharacters = chars
				}
			}

			result, err := pipeline.Process(context.Background(), source, sourceType, opts)
			if err != nil {
				return screens.ProcessingErrorMsg{Err: err}
			}
			return screens.ProcessingCompleteMsg{
				Title:      result.Title,
				Passage:    result.Passage,
				Source:     result.Source,
				SourceType: result.SourceType,
			}
		},
	)
}

func (m Model) updateProcessing(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		keyStr := msg.String()
		if keyStr == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.Key().Code == tea.KeyEscape {
			m.screen = ScreenInput
			return m, nil
		}

	case screens.ProcessingCompleteMsg:
		m.currentTitle = msg.Title
		m.typing = screens.NewTypingModel(msg.Passage, msg.Title, m.cfg)
		m.typing.SetSize(m.width, m.height)
		m.screen = ScreenTyping
		return m, nil

	case screens.ProcessingErrorMsg:
		m.processing.SetError(msg.Err.Error())
		return m, nil
	}

	var cmd tea.Cmd
	m.processing, cmd = m.processing.Update(msg)
	return m, cmd
}

func (m Model) updateTyping(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.typing, cmd = m.typing.Update(msg)

	if m.typing.IsComplete() {
		engine := m.typing.Engine()
		stats := typing.SessionStats(engine.Snapshot())
		timeUp := m.typing.IsTimeUp()

		session := typing.NewSessionFromEngine(
			uuid.New().String(),
			engine,
			m.currentSourceType,
			m.currentSource,
			m.currentTitle,
		)

		repo := m.repo
		saveCmd := func() tea.Msg {
			if repo != nil {
				_ = repo.Save(context.Background(), &session)
			}
			return nil
		}

		m.results = screens.NewResultsModel(stats, m.currentTitle, timeUp)
		m.results.SetSize(m.width, m.height)
		m.screen = ScreenResults
		return m, saveCmd
	}

	return m, cmd
}

func (m Model) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		keyStr := msg.String()
		if keyStr == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.Key().Code == tea.KeyEnter {
			switch m.results.SelectedAction() {
			case "again":
				m.screen = ScreenInput
				m.input.Reset()
				return m, nil
			case "home":
				m.screen = ScreenHome
				m.home = screens.NewHomeModel()
				m.home.SetSize(m.width, m.height)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.results, cmd = m.results.Update(msg)
	return m, cmd
}

func (m Model) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		keyStr := msg.String()
		if keyStr == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.Key().Code == tea.KeyEscape {
			m.screen = ScreenHome
			m.home = screens.NewHomeModel()
			m.home.SetSize(m.width, m.height)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.history, cmd = m.history.Update(msg)
	return m, cmd
}

func (m Model) screenContext() ScreenContext {
	ctx := ScreenContext{
		InputMode: m.input.Mode(),
	}

	if m.screen == ScreenTyping {
		ctx.TypingTitle = m.typing.Title()
		ctx.TypingTimer = m.typing.TimerLine()
		if m.cfg != nil && m.cfg.IsTimedMode() {
			ctx.TimedMode = true
		}
	}

	if m.screen == ScreenProcessing {
		ctx.ProcessingErr = m.processing.HasError()
	}

	return ctx
}

func (m Model) screenContent() string {
	switch m.screen {
	case ScreenHome:
		return m.home.View()
	case ScreenInput:
		return m.input.View()
	case ScreenProcessing:
		return m.processing.View()
	case ScreenTyping:
		return m.typing.View()
	case ScreenResults:
		return m.results.View()
	case ScreenHistory:
		return m.history.View()
	default:
		return "Unknown screen"
	}
}

func (m Model) paletteContext() components.PaletteContext {
	return components.PaletteContext{
		ShowRestart: m.screen == ScreenTyping || m.screen == ScreenResults,
	}
}

func (m Model) overlayContent() (string, bool) {
	if m.showSettings {
		return m.settings.View(), true
	}
	if m.showPalette {
		return m.palette.View(), true
	}
	if m.showHelp {
		return m.help.View(), true
	}
	if m.showQuit {
		return m.quitDlg.View(), true
	}
	return "", false
}

func (m Model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("")
	}

	bindings, align := ScreenMeta(m.screen, m.screenContext())

	shellView := layout.Shell{
		Width:        m.width,
		Height:       m.height,
		Bindings:     bindings,
		Content:      m.screenContent(),
		ContentAlign: align,
	}.View()

	if overlay, isOverlay := m.overlayContent(); isOverlay {
		shellView = layout.RenderOverlay(shellView, overlay, m.width, m.height)
	}

	v := tea.NewView(shellView)
	v.AltScreen = true
	return v
}
