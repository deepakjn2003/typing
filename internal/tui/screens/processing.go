package screens

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/deepakjn2003/typing/internal/tui/components"
	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type ProcessingTickMsg time.Time

type ProcessingCompleteMsg struct {
	Title      string
	Passage    string
	Source     string
	SourceType string
}

type ProcessingErrorMsg struct {
	Err error
}

type StageState int

const (
	StagePending StageState = iota
	StageActive
	StageDone
)

type ProcessingStage struct {
	Label string
	State StageState
}

type ProcessingModel struct {
	spinner       components.Spinner
	stages        []ProcessingStage
	activeStage   int
	tickCount     int
	err           string
	width, height int
}

var defaultStages = []string{
	"Extracting content",
	"Normalizing text",
	"Generating typing passage",
	"Preparing session",
}

func NewProcessingModel() ProcessingModel {
	return ProcessingModel{
		spinner: components.NewSpinner(""),
		stages:  newStageList(0),
	}
}

func newStageList(activeIdx int) []ProcessingStage {
	stages := make([]ProcessingStage, len(defaultStages))
	for i, label := range defaultStages {
		state := StagePending
		if i < activeIdx {
			state = StageDone
		} else if i == activeIdx {
			state = StageActive
		}
		stages[i] = ProcessingStage{Label: label, State: state}
	}
	return stages
}

func (m *ProcessingModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *ProcessingModel) SetError(err string) {
	m.err = err
}

func (m *ProcessingModel) Reset() {
	m.stages = newStageList(0)
	m.activeStage = 0
	m.tickCount = 0
	m.spinner = components.NewSpinner("")
	m.err = ""
}

func (m ProcessingModel) Update(msg tea.Msg) (ProcessingModel, tea.Cmd) {
	switch msg.(type) {
	case ProcessingTickMsg:
		m.spinner.Tick()
		m.tickCount++

		// Cosmetic stage advancement (~1.5s per stage, cap before last stage)
		if m.tickCount%19 == 0 && m.activeStage < len(m.stages)-2 {
			m.activeStage++
			m.stages = newStageList(m.activeStage)
		}

		return m, tickProcessing()
	}
	return m, nil
}

func tickProcessing() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return ProcessingTickMsg(t)
	})
}

func StartProcessingTick() tea.Cmd {
	return tickProcessing()
}

func (m ProcessingModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Preparing your typing passage"))
	b.WriteString("\n\n")

	for _, stage := range m.stages {
		var prefix string
		switch stage.State {
		case StageDone:
			prefix = styles.SuccessStyle.Render("✓")
		case StageActive:
			frame := spinnerFrames[m.spinner.Frame%len(spinnerFrames)]
			prefix = styles.InfoStyle.Render(frame)
		default:
			prefix = styles.MutedStyle.Render("○")
		}
		b.WriteString("  ")
		b.WriteString(prefix)
		b.WriteString(" ")
		b.WriteString(styles.BodyStyle.Render(stage.Label))
		b.WriteString("\n")
	}

	if m.err != "" {
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err))
		b.WriteString("\n")
	}

	return b.String()
}

func (m ProcessingModel) HasError() bool {
	return m.err != ""
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
