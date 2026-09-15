package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	created_at DATETIME NOT NULL,
	source_type TEXT NOT NULL,
	source TEXT NOT NULL,
	title TEXT NOT NULL,
	text TEXT NOT NULL,
	duration_ms INTEGER NOT NULL,
	wpm REAL NOT NULL,
	raw_wpm REAL NOT NULL,
	accuracy REAL NOT NULL,
	characters INTEGER NOT NULL,
	correct_characters INTEGER NOT NULL,
	incorrect_characters INTEGER NOT NULL,
	errors INTEGER NOT NULL,
	corrected_errors INTEGER NOT NULL,
	backspaces INTEGER NOT NULL,
	error_details TEXT
);

CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at DESC);
`

func OpenDatabase(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating database directory %s: %w", dir, err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("setting WAL mode: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("applying schema: %w", err)
	}

	return db, nil
}
