package typing

import "time"

type TypingError struct {
	Position int
	Expected rune
	Actual   rune
	At       time.Duration
}

type ErrorPair struct {
	Expected rune
	Actual   rune
}

type ErrorFrequency struct {
	Pair  ErrorPair
	Count int
}
