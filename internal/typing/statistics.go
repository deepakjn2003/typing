package typing

import (
	"math"
	"sort"
	"time"
)

func CalculateWPM(correctChars int, elapsed time.Duration) float64 {
	if elapsed <= 0 || correctChars <= 0 {
		return 0
	}
	minutes := elapsed.Minutes()
	if minutes == 0 {
		return 0
	}
	words := float64(correctChars) / 5.0
	return math.Round((words/minutes)*10) / 10
}

func CalculateRawWPM(totalChars int, elapsed time.Duration) float64 {
	if elapsed <= 0 || totalChars <= 0 {
		return 0
	}
	minutes := elapsed.Minutes()
	if minutes == 0 {
		return 0
	}
	words := float64(totalChars) / 5.0
	return math.Round((words/minutes)*10) / 10
}

func CalculateAccuracy(correct, total int) float64 {
	if total <= 0 {
		return 100.0
	}
	acc := (float64(correct) / float64(total)) * 100.0
	return math.Round(acc*10) / 10
}

type ErrorAnalysis struct {
	TotalErrors  int
	UniqueErrors int
	TopErrors    []ErrorFrequency
}

func AnalyzeErrors(errors []TypingError) ErrorAnalysis {
	if len(errors) == 0 {
		return ErrorAnalysis{}
	}

	freq := make(map[ErrorPair]int)
	for _, e := range errors {
		pair := ErrorPair{Expected: e.Expected, Actual: e.Actual}
		freq[pair]++
	}

	frequencies := make([]ErrorFrequency, 0, len(freq))
	for pair, count := range freq {
		frequencies = append(frequencies, ErrorFrequency{Pair: pair, Count: count})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].Count > frequencies[j].Count
	})

	return ErrorAnalysis{
		TotalErrors:  len(errors),
		UniqueErrors: len(frequencies),
		TopErrors:    frequencies,
	}
}

type Stats struct {
	WPM             float64
	RawWPM          float64
	Accuracy        float64
	Elapsed         time.Duration
	Characters      int
	CorrectChars    int
	IncorrectChars  int
	TotalErrors     int
	Corrections     int
	Backspaces      int
	TotalKeystrokes int
	ErrorAnalysis   ErrorAnalysis
	TimeUp          bool
}

func SessionStats(snap EngineSnapshot) Stats {
	return Stats{
		WPM:             CalculateWPM(snap.CorrectChars, snap.Elapsed),
		RawWPM:          CalculateRawWPM(snap.InputLen, snap.Elapsed),
		Accuracy:        CalculateAccuracy(snap.CorrectChars, snap.InputLen),
		Elapsed:         snap.Elapsed,
		Characters:      snap.InputLen,
		CorrectChars:    snap.CorrectChars,
		IncorrectChars:  snap.IncorrectChars,
		TotalErrors:     snap.TotalErrors,
		Corrections:     snap.Corrections,
		Backspaces:      snap.Backspaces,
		TotalKeystrokes: snap.TotalKeystrokes,
		ErrorAnalysis:   AnalyzeErrors(snap.Errors),
		TimeUp:          snap.TimeUp,
	}
}
