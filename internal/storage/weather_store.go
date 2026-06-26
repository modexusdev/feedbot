// modexusBot/internal/storage/weather_store.go

package storage

import (
	"database/sql"
	"time"
)

func createWeatherTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS weather_location (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		country TEXT,
		latitude REAL NOT NULL,
		longitude REAL NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}

func SaveWeatherLocation(location WeatherLocation) (WeatherLocation, error) {
	db, err := openDB()
	if err != nil {
		return WeatherLocation{}, err
	}
	defer db.Close()

	now := time.Now().Format(time.RFC3339)

	location.ID = "default"

	if location.CreatedAt == "" {
		location.CreatedAt = now
	}

	location.UpdatedAt = now

	query := `
	INSERT INTO weather_location (
		id,
		name,
		country,
		latitude,
		longitude,
		created_at,
		updated_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		country = excluded.country,
		latitude = excluded.latitude,
		longitude = excluded.longitude,
		updated_at = excluded.updated_at;
	`

	_, err = db.Exec(
		query,
		location.ID,
		location.Name,
		location.Country,
		location.Latitude,
		location.Longitude,
		location.CreatedAt,
		location.UpdatedAt,
	)

	if err != nil {
		return WeatherLocation{}, err
	}

	return location, nil
}

func GetWeatherLocation() (WeatherLocation, error) {
	db, err := openDB()
	if err != nil {
		return WeatherLocation{}, err
	}
	defer db.Close()

	var location WeatherLocation

	query := `
	SELECT
		id,
		name,
		country,
		latitude,
		longitude,
		created_at,
		updated_at
	FROM weather_location
	WHERE id = 'default'
	LIMIT 1;
	`

	err = db.QueryRow(query).Scan(
		&location.ID,
		&location.Name,
		&location.Country,
		&location.Latitude,
		&location.Longitude,
		&location.CreatedAt,
		&location.UpdatedAt,
	)

	return location, err
}
