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
	// create data directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}
	// open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	// run database migrations
	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, err
	}
	// return database connection
	return db, nil
}

func runMigrations(db *sql.DB) error {
	// create youtube table
	if err := createYoutubeTable(db); err != nil {
		return err
	}

	// create weather table
	if err := createWeatherTable(db); err != nil {
		return err
	}
	if err := createConfigTable(db); err != nil {
		return err
	}

	return nil
}
