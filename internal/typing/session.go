package typing

import "time"

type Session struct {
	ID                  string
	CreatedAt           time.Time
	SourceType          string
	Source              string
	Title               string
	Text                string
	Duration            time.Duration
	WPM                 float64
	RawWPM              float64
	Accuracy            float64
	Characters          int
	CorrectCharacters   int
	IncorrectCharacters int
	Errors              int
	CorrectedErrors     int
	Backspaces          int
	ErrorDetails        []TypingError
}

func NewSessionFromEngine(id string, engine *Engine, sourceType, source, title string) Session {
	snap := engine.Snapshot()
	stats := SessionStats(snap)

	displaySource := source
	if sourceType == "paste" && len(source) > 100 {
		displaySource = source[:100]
	}

	return Session{
		ID:                  id,
		CreatedAt:           time.Now(),
		SourceType:          sourceType,
		Source:              displaySource,
		Title:               title,
		Text:                engine.Target(),
		Duration:            stats.Elapsed,
		WPM:                 stats.WPM,
		RawWPM:              stats.RawWPM,
		Accuracy:            stats.Accuracy,
		Characters:          stats.Characters,
		CorrectCharacters:   stats.CorrectChars,
		IncorrectCharacters: stats.IncorrectChars,
		Errors:              stats.TotalErrors,
		CorrectedErrors:     stats.Corrections,
		Backspaces:          stats.Backspaces,
		ErrorDetails:        snap.Errors,
	}
}
