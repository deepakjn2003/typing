package screens

import (
	"strings"
)

type textLine struct {
	start int
	end   int
	text  string
}

func wrapPassage(target []rune, lineWidth int) []textLine {
	if lineWidth < 10 {
		lineWidth = 10
	}

	var lines []textLine
	start := 0
	currentLen := 0

	for i, r := range target {
		if r == '\n' {
			lines = append(lines, textLine{start: start, end: i + 1, text: string(target[start : i+1])})
			start = i + 1
			currentLen = 0
			continue
		}

		currentLen++
		if currentLen >= lineWidth && r == ' ' {
			lines = append(lines, textLine{start: start, end: i + 1, text: string(target[start : i+1])})
			start = i + 1
			currentLen = 0
		}
	}

	if start < len(target) {
		lines = append(lines, textLine{start: start, end: len(target), text: string(target[start:])})
	}

	if len(lines) == 0 {
		lines = append(lines, textLine{start: 0, end: 0, text: ""})
	}

	return lines
}

func cursorLineIndex(lines []textLine, cursorPos int) int {
	for i, line := range lines {
		if cursorPos >= line.start && cursorPos < line.end {
			return i
		}
	}
	if cursorPos >= lines[len(lines)-1].end {
		return len(lines) - 1
	}
	return 0
}

func BuildTypingWindow(
	target []rune,
	cursorPos int,
	lineWidth int,
	visibleLines int,
	lineSpacing int,
	styleChar func(pos int, r rune) string,
) string {
	lines := wrapPassage(target, lineWidth)
	cursorLine := cursorLineIndex(lines, cursorPos)

	if visibleLines < 2 {
		visibleLines = 2
	}
	if visibleLines > 5 {
		visibleLines = 5
	}

	half := visibleLines / 2
	startLine := cursorLine - half
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + visibleLines
	if endLine > len(lines) {
		endLine = len(lines)
		startLine = endLine - visibleLines
		if startLine < 0 {
			startLine = 0
		}
	}

	var rendered []string
	for i := startLine; i < endLine; i++ {
		line := lines[i]
		var b strings.Builder
		for pos := line.start; pos < line.end; pos++ {
			b.WriteString(styleChar(pos, target[pos]))
		}
		rendered = append(rendered, b.String())
		if lineSpacing > 0 && i < endLine-1 {
			for s := 0; s < lineSpacing; s++ {
				rendered = append(rendered, "")
			}
		}
	}

	return strings.Join(rendered, "\n")
}
