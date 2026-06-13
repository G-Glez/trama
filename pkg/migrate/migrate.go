package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	TableName     string
	ColumnName    string
	MigrationsDir string
}

func (c *Config) defaults() {
	if c.TableName == "" {
		c.TableName = "schema_version"
	}
	if c.ColumnName == "" {
		c.ColumnName = "version"
	}
	if c.MigrationsDir == "" {
		c.MigrationsDir = "migrations"
	}
}

var bootstrapSQL = `CREATE TABLE IF NOT EXISTS %s (%s INTEGER NOT NULL);
INSERT INTO %[1]s (%[2]s) SELECT 0 WHERE NOT EXISTS (SELECT 1 FROM %[1]s);`

func Run(db *sql.DB, cfg ...Config) error {
	var c Config
	if len(cfg) > 0 {
		c = cfg[0]
	}
	c.defaults()

	if _, err := db.Exec(fmt.Sprintf(bootstrapSQL, c.TableName, c.ColumnName)); err != nil {
		return fmt.Errorf("bootstrap %s: %w", c.TableName, err)
	}

	current, err := currentVersion(db, c.TableName, c.ColumnName)
	if err != nil {
		return fmt.Errorf("read current version: %w", err)
	}

	entries, err := filepath.Glob(filepath.Join(c.MigrationsDir, "v*.sql"))
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(entries)

	if len(entries) > 0 {
		maxV, err := parseVersion(entries[len(entries)-1])
		if err != nil {
			return fmt.Errorf("parse max version: %w", err)
		}

		if current > maxV {
			return fmt.Errorf("schema version %d exceeds maximum available migration v%04d", current, maxV)
		}
	}

	for _, entry := range entries {
		v, err := parseVersion(entry)
		if err != nil {
			return fmt.Errorf("parse %s: %w", entry, err)
		}

		if v <= current {
			continue
		}

		sqlContent, err := os.ReadFile(entry)
		if err != nil {
			return fmt.Errorf("read %s: %w", entry, err)
		}

		if _, err := db.Exec(string(sqlContent)); err != nil {
			return fmt.Errorf("apply %s: %w", entry, err)
		}

		updateSQL := fmt.Sprintf("UPDATE %s SET %s = ?", c.TableName, c.ColumnName)
		if _, err := db.Exec(updateSQL, v); err != nil {
			return fmt.Errorf("update %s to %d: %w", c.TableName, v, err)
		}
	}

	return nil
}

func currentVersion(db *sql.DB, tableName, columnName string) (int, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", columnName, tableName)
	var v int
	if err := db.QueryRow(query).Scan(&v); err != nil {
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
