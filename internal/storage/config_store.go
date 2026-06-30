// modexusBot/internal/storage/config_store.go

package storage

import "database/sql"

func createConfigTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			id INTEGER PRIMARY KEY CHECK(id = 1),
			language TEXT NOT NULL DEFAULT 'en'
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT OR IGNORE INTO config (id, language)
		VALUES (1, 'en')
	`)

	return err
}
func GetLanguage() (string, error) {
	db, err := openDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var language string

	err = db.QueryRow(`
		SELECT language
		FROM config
		WHERE id = 1
	`).Scan(&language)

	return language, err
}

func SaveLanguage(language string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE config
		SET language = ?
		WHERE id = 1
	`, language)

	return err
}
