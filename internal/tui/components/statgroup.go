package components

import (
	"fmt"
	"strings"

	"github.com/deepakjn2003/typing/internal/tui/styles"
)

type Stat struct {
	Label string
	Value string
	Width int
}

type StatGroup struct {
	Stats []Stat
}

func (g StatGroup) View() string {
	var parts []string
	for _, stat := range g.Stats {
		width := stat.Width
		if width < 1 {
			width = 8
		}
		label := styles.StatLabelStyle.Render(stat.Label)
		value := styles.StatValueStyle.Render(fmt.Sprintf("%-*s", width, stat.Value))
		parts = append(parts, label+"  "+value)
	}
	return strings.Join(parts, "      ")
}
