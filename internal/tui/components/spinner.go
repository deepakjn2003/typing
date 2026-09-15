package components

import (
	"github.com/deepakjn2003/typing/internal/tui/styles"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type Spinner struct {
	Message string
	Frame   int
}

func NewSpinner(message string) Spinner {
	return Spinner{Message: message}
}

func (s *Spinner) Tick() {
	s.Frame = (s.Frame + 1) % len(spinnerFrames)
}

func (s *Spinner) SetMessage(msg string) {
	s.Message = msg
}

func (s Spinner) View() string {
	frame := styles.InfoStyle.Render(spinnerFrames[s.Frame])
	message := styles.MutedStyle.Render(s.Message)
	return frame + " " + message
}
