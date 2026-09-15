package typing

import (
	"time"
)

type EngineState int

const (
	StateNotStarted EngineState = iota
	StateRunning
	StatePaused
	StateCompleted
)

type Engine struct {
	target []rune
	input  []rune
	state  EngineState

	startTime time.Time
	endTime   time.Time
	pauseTime time.Time
	pausedDur time.Duration

	errors []TypingError

	totalKeystrokes int
	backspaces      int
	corrections     int

	timeLimit time.Duration
	timeUp    bool
}

func NewEngine(target string) *Engine {
	return &Engine{
		target: []rune(target),
		input:  make([]rune, 0, len(target)),
		state:  StateNotStarted,
		errors: make([]TypingError, 0),
	}
}

func (e *Engine) SetTimeLimit(d time.Duration) {
	e.timeLimit = d
}

func (e *Engine) IsTimeUp() bool {
	return e.timeUp
}

func (e *Engine) State() EngineState {
	return e.state
}

func (e *Engine) Target() string {
	return string(e.target)
}

func (e *Engine) TargetRunes() []rune {
	return e.target
}

func (e *Engine) Input() string {
	return string(e.input)
}

func (e *Engine) InputRunes() []rune {
	return e.input
}

func (e *Engine) Position() int {
	return len(e.input)
}

func (e *Engine) TargetLength() int {
	return len(e.target)
}

func (e *Engine) IsComplete() bool {
	return e.state == StateCompleted
}

func (e *Engine) Errors() []TypingError {
	result := make([]TypingError, len(e.errors))
	copy(result, e.errors)
	return result
}

func (e *Engine) ProcessRune(r rune) {
	if e.state == StateCompleted || e.state == StatePaused {
		return
	}

	if e.state == StateNotStarted {
		e.state = StateRunning
		e.startTime = time.Now()
	}

	e.totalKeystrokes++

	pos := len(e.input)
	if pos >= len(e.target) {
		return
	}

	e.input = append(e.input, r)

	expected := e.target[pos]
	if r != expected {
		e.errors = append(e.errors, TypingError{
			Position: pos,
			Expected: expected,
			Actual:   r,
			At:       e.Elapsed(),
		})
	}

	if len(e.input) >= len(e.target) {
		e.state = StateCompleted
		e.endTime = time.Now()
	}
}

func (e *Engine) Backspace() {
	if e.state != StateRunning || len(e.input) == 0 {
		return
	}

	e.totalKeystrokes++
	e.backspaces++

	lastPos := len(e.input) - 1
	if e.input[lastPos] != e.target[lastPos] {
		e.corrections++
	}

	e.input = e.input[:lastPos]
}

func (e *Engine) CheckTimeLimit() bool {
	if e.timeLimit <= 0 || e.state == StateCompleted || e.state == StateNotStarted {
		return false
	}
	if e.Elapsed() >= e.timeLimit {
		e.timeUp = true
		e.state = StateCompleted
		e.endTime = time.Now()
		return true
	}
	return false
}

func (e *Engine) Pause() {
	if e.state != StateRunning {
		return
	}
	e.state = StatePaused
	e.pauseTime = time.Now()
}

func (e *Engine) Resume() {
	if e.state != StatePaused {
		return
	}
	e.pausedDur += time.Since(e.pauseTime)
	e.state = StateRunning
}

func (e *Engine) Elapsed() time.Duration {
	switch e.state {
	case StateNotStarted:
		return 0
	case StateCompleted:
		return e.endTime.Sub(e.startTime) - e.pausedDur
	case StatePaused:
		return e.pauseTime.Sub(e.startTime) - e.pausedDur
	default:
		return time.Since(e.startTime) - e.pausedDur
	}
}

func (e *Engine) CorrectCount() int {
	correct := 0
	for i, r := range e.input {
		if i < len(e.target) && r == e.target[i] {
			correct++
		}
	}
	return correct
}

func (e *Engine) IncorrectCount() int {
	incorrect := 0
	for i, r := range e.input {
		if i < len(e.target) && r != e.target[i] {
			incorrect++
		}
	}
	return incorrect
}

func (e *Engine) TotalErrors() int {
	return len(e.errors)
}

func (e *Engine) TotalKeystrokes() int {
	return e.totalKeystrokes
}

func (e *Engine) Backspaces() int {
	return e.backspaces
}

func (e *Engine) Corrections() int {
	return e.corrections
}

type CharState int

const (
	CharUntyped CharState = iota
	CharCorrect
	CharIncorrect
	CharCursor
)

func (e *Engine) CharStateAt(pos int) CharState {
	if pos < 0 || pos >= len(e.target) {
		return CharUntyped
	}

	inputLen := len(e.input)
	if pos == inputLen {
		return CharCursor
	}
	if pos > inputLen {
		return CharUntyped
	}
	if e.input[pos] == e.target[pos] {
		return CharCorrect
	}
	return CharIncorrect
}

func (e *Engine) Snapshot() EngineSnapshot {
	return EngineSnapshot{
		TargetLen:       len(e.target),
		InputLen:        len(e.input),
		CorrectChars:    e.CorrectCount(),
		IncorrectChars:  e.IncorrectCount(),
		TotalErrors:     len(e.errors),
		Corrections:     e.corrections,
		Backspaces:      e.backspaces,
		TotalKeystrokes: e.totalKeystrokes,
		Elapsed:         e.Elapsed(),
		Errors:          e.Errors(),
		State:           e.state,
		TimeUp:          e.timeUp,
	}
}

type EngineSnapshot struct {
	TargetLen       int
	InputLen        int
	CorrectChars    int
	IncorrectChars  int
	TotalErrors     int
	Corrections     int
	Backspaces      int
	TotalKeystrokes int
	Elapsed         time.Duration
	Errors          []TypingError
	State           EngineState
	TimeUp          bool
}
