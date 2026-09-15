package styles

import (
	"charm.land/lipgloss/v2"
	"github.com/deepakjn2003/typing/internal/config"
)

type TextSizePreset int

const (
	TextSizeSmall  TextSizePreset = 0
	TextSizeMedium TextSizePreset = 1
	TextSizeLarge  TextSizePreset = 2
)

func ParseTextSize(s string) TextSizePreset {
	switch s {
	case config.TextSizeSmall:
		return TextSizeSmall
	case config.TextSizeLarge:
		return TextSizeLarge
	default:
		return TextSizeMedium
	}
}

func (p TextSizePreset) String() string {
	switch p {
	case TextSizeSmall:
		return config.TextSizeSmall
	case TextSizeLarge:
		return config.TextSizeLarge
	default:
		return config.TextSizeMedium
	}
}

func (p TextSizePreset) Label() string {
	switch p {
	case TextSizeSmall:
		return "Small"
	case TextSizeLarge:
		return "Large"
	default:
		return "Medium"
	}
}

func (p TextSizePreset) CharsPerLine(termWidth int) int {
	base := termWidth - 20
	if base < 30 {
		base = 30
	}
	switch p {
	case TextSizeSmall:
		return base
	case TextSizeLarge:
		// Fewer chars per line so each character reads larger on screen.
		n := int(float64(base) / 2.2)
		if n < 24 {
			n = 24
		}
		return n
	default:
		return int(float64(base) / 1.35)
	}
}

func (p TextSizePreset) LineSpacing() int {
	switch p {
	case TextSizeSmall:
		return 0
	case TextSizeLarge:
		return 2
	default:
		return 1
	}
}

// Horizontal padding applied around each rendered character (terminal "size" simulation).
func (p TextSizePreset) CharPadding() int {
	switch p {
	case TextSizeSmall:
		return 0
	case TextSizeLarge:
		return 1
	default:
		return 0
	}
}

func (p TextSizePreset) StyleChar(base lipgloss.Style) lipgloss.Style {
	switch p {
	case TextSizeLarge:
		return base.Bold(true)
	case TextSizeSmall:
		return base
	default:
		return base
	}
}

func (p TextSizePreset) StatsStyle(base lipgloss.Style) lipgloss.Style {
	if p == TextSizeLarge {
		return base.Bold(true)
	}
	return base
}
