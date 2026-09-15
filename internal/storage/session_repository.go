package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/deepakjn2003/typing/internal/typing"
)

type WeakCharacter struct {
	Char  rune
	Count int
}

type SessionRepository interface {
	Save(ctx context.Context, session *typing.Session) error
	List(ctx context.Context, limit, offset int) ([]typing.Session, error)
	GetByID(ctx context.Context, id string) (*typing.Session, error)
	Recent(ctx context.Context, limit int) ([]typing.Session, error)
	WeakCharacters(ctx context.Context, limit int) ([]WeakCharacter, error)
}

type SQLiteSessionRepository struct {
	db *sql.DB
}

func NewSQLiteSessionRepository(db *sql.DB) *SQLiteSessionRepository {
	return &SQLiteSessionRepository{db: db}
}

func (r *SQLiteSessionRepository) Save(ctx context.Context, session *typing.Session) error {
	var errorDetailsJSON []byte
	var err error
	if len(session.ErrorDetails) > 0 {
		errorDetailsJSON, err = json.Marshal(session.ErrorDetails)
		if err != nil {
			return fmt.Errorf("encoding error details: %w", err)
		}
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO sessions (
			id, created_at, source_type, source, title, text,
			duration_ms, wpm, raw_wpm, accuracy,
			characters, correct_characters, incorrect_characters,
			errors, corrected_errors, backspaces, error_details
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		session.ID,
		session.CreatedAt.UTC().Format(time.RFC3339),
		session.SourceType,
		session.Source,
		session.Title,
		session.Text,
		session.Duration.Milliseconds(),
		session.WPM,
		session.RawWPM,
		session.Accuracy,
		session.Characters,
		session.CorrectCharacters,
		session.IncorrectCharacters,
		session.Errors,
		session.CorrectedErrors,
		session.Backspaces,
		string(errorDetailsJSON),
	)
	if err != nil {
		return fmt.Errorf("saving session: %w", err)
	}
	return nil
}

func (r *SQLiteSessionRepository) List(ctx context.Context, limit, offset int) ([]typing.Session, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, created_at, source_type, source, title, text,
			   duration_ms, wpm, raw_wpm, accuracy,
			   characters, correct_characters, incorrect_characters,
			   errors, corrected_errors, backspaces, error_details
		FROM sessions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("querying sessions: %w", err)
	}
	defer rows.Close()

	return scanSessions(rows)
}

func (r *SQLiteSessionRepository) GetByID(ctx context.Context, id string) (*typing.Session, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, created_at, source_type, source, title, text,
			   duration_ms, wpm, raw_wpm, accuracy,
			   characters, correct_characters, incorrect_characters,
			   errors, corrected_errors, backspaces, error_details
		FROM sessions
		WHERE id = ?
	`, id)

	session, err := scanSession(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting session %s: %w", id, err)
	}
	return session, nil
}

func (r *SQLiteSessionRepository) Recent(ctx context.Context, limit int) ([]typing.Session, error) {
	return r.List(ctx, limit, 0)
}

func scanSessions(rows *sql.Rows) ([]typing.Session, error) {
	var sessions []typing.Session
	for rows.Next() {
		var s typing.Session
		var createdAtStr string
		var durationMs int64
		var errorDetailsJSON sql.NullString

		err := rows.Scan(
			&s.ID, &createdAtStr, &s.SourceType, &s.Source, &s.Title, &s.Text,
			&durationMs, &s.WPM, &s.RawWPM, &s.Accuracy,
			&s.Characters, &s.CorrectCharacters, &s.IncorrectCharacters,
			&s.Errors, &s.CorrectedErrors, &s.Backspaces, &errorDetailsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning session row: %w", err)
		}

		s.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		s.Duration = time.Duration(durationMs) * time.Millisecond

		if errorDetailsJSON.Valid && errorDetailsJSON.String != "" {
			if err := json.Unmarshal([]byte(errorDetailsJSON.String), &s.ErrorDetails); err != nil {
				s.ErrorDetails = nil
			}
		}

		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

func scanSession(row *sql.Row) (*typing.Session, error) {
	var s typing.Session
	var createdAtStr string
	var durationMs int64
	var errorDetailsJSON sql.NullString

	err := row.Scan(
		&s.ID, &createdAtStr, &s.SourceType, &s.Source, &s.Title, &s.Text,
		&durationMs, &s.WPM, &s.RawWPM, &s.Accuracy,
		&s.Characters, &s.CorrectCharacters, &s.IncorrectCharacters,
		&s.Errors, &s.CorrectedErrors, &s.Backspaces, &errorDetailsJSON,
	)
	if err != nil {
		return nil, err
	}

	s.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	s.Duration = time.Duration(durationMs) * time.Millisecond

	if errorDetailsJSON.Valid && errorDetailsJSON.String != "" {
		if err := json.Unmarshal([]byte(errorDetailsJSON.String), &s.ErrorDetails); err != nil {
			s.ErrorDetails = nil
		}
	}

	return &s, nil
}

func (r *SQLiteSessionRepository) WeakCharacters(ctx context.Context, limit int) ([]WeakCharacter, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT error_details FROM sessions
		WHERE error_details IS NOT NULL AND error_details != ''
		ORDER BY created_at DESC
		LIMIT 20
	`)
	if err != nil {
		return nil, fmt.Errorf("querying sessions for weak characters: %w", err)
	}
	defer rows.Close()

	freq := make(map[rune]int)
	for rows.Next() {
		var detailsJSON string
		if err := rows.Scan(&detailsJSON); err != nil {
			continue
		}
		var errors []typing.TypingError
		if err := json.Unmarshal([]byte(detailsJSON), &errors); err != nil {
			continue
		}
		for _, e := range errors {
			freq[e.Expected]++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating session rows: %w", err)
	}

	if len(freq) == 0 {
		return nil, nil
	}

	result := make([]WeakCharacter, 0, len(freq))
	for ch, count := range freq {
		result = append(result, WeakCharacter{Char: ch, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	if len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}
