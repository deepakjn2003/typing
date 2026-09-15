package layout

const MaxContentWidth = 90

func ContentWidth(termWidth int) int {
	w := termWidth - 4
	if w > MaxContentWidth {
		return MaxContentWidth
	}
	if w < 1 {
		return 1
	}
	return w
}

func ContentPadding(termWidth int) int {
	pad := (termWidth - ContentWidth(termWidth)) / 2
	if pad < 0 {
		return 0
	}
	return pad
}
