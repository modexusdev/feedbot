// modexusBot/internal/storage/db.go
package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const dbPath = "data/feedbot.db"

func openDB() (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	if err := createYoutubeTable(db); err != nil {
		return err
	}

	if err := createWeatherTable(db); err != nil {
		return err
	}

	return nil
}
