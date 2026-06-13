package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Run(db *sql.DB) error {
	current, err := currentVersion(db)
	if err != nil {
		return fmt.Errorf("read current version: %w", err)
	}

	entries, err := fs.Glob(migrationsFS, "migrations/v*.sql")
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(entries)

	for _, entry := range entries {
		v, err := parseVersion(entry)
		if err != nil {
			return fmt.Errorf("parse %s: %w", entry, err)
		}

		if v <= current {
			continue
		}

		sqlContent, err := migrationsFS.ReadFile(entry)
		if err != nil {
			return fmt.Errorf("read %s: %w", entry, err)
		}

		if _, err := db.Exec(string(sqlContent)); err != nil {
			return fmt.Errorf("apply %s: %w", entry, err)
		}

		if _, err := db.Exec("UPDATE schema_version SET version = ?", v); err != nil {
			return fmt.Errorf("update schema version to %d: %w", v, err)
		}

		fmt.Printf("  Applied migration v%04d\n", v)
	}

	return nil
}

func currentVersion(db *sql.DB) (int, error) {
	var exists int
	err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='schema_version'`).Scan(&exists)
	if err != nil {
		return -1, nil
	}
	if exists == 0 {
		return -1, nil
	}

	var v int
	if err := db.QueryRow(`SELECT version FROM schema_version`).Scan(&v); err != nil {
		return -1, nil
	}
	return v, nil
}

func parseVersion(path string) (int, error) {
	base := filepath.Base(path)
	base = strings.TrimPrefix(base, "v")
	base = strings.TrimSuffix(base, ".sql")
	return strconv.Atoi(base)
}
