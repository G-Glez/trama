package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS greetings (
			id INTEGER PRIMARY KEY AUTOINCREMENT
		)`,
		`CREATE TABLE IF NOT EXISTS game_systems (
			id         TEXT PRIMARY KEY,
			name       TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS editions (
			id             TEXT PRIMARY KEY,
			game_system_id TEXT NOT NULL REFERENCES game_systems(id),
			name           TEXT NOT NULL,
			version        TEXT NOT NULL,
			created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS factions (
			id         TEXT PRIMARY KEY,
			edition_id TEXT NOT NULL REFERENCES editions(id),
			name       TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
