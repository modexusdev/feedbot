// modexusBot/internal/storage/youtube_store.go
package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/modexusdev/feedbot/internal/helper"

	_ "modernc.org/sqlite"
)

const youtubeDB = "data/feedbot.db"

func openYoutubeDB() (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(youtubeDB), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", youtubeDB)
	if err != nil {
		return nil, err
	}

	if err := createYoutubeTable(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createYoutubeTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS youtube_channels (
			id TEXT PRIMARY KEY,
			name TEXT,
			handle TEXT UNIQUE,
			rss_url TEXT UNIQUE,
			avatar_url TEXT,
			last_video_id TEXT,
			created_at TEXT,
			updated_at TEXT
		);
	`)
	return err
}

func SaveYoutubeChannel(channel YoutubeChannel) (YoutubeChannel, error) {
	db, err := openYoutubeDB()
	if err != nil {
		return channel, err
	}
	defer db.Close()

	now := time.Now().Format(time.RFC3339)

	existing, found, err := findYoutubeChannel(db, channel.Handle, channel.RSSURL)
	if err != nil {
		return channel, err
	}

	if found {
		channel.ID = existing.ID
		channel.CreatedAt = existing.CreatedAt
		channel.UpdatedAt = now

		if channel.Name == "" {
			channel.Name = existing.Name
		}
		if channel.Handle == "" {
			channel.Handle = existing.Handle
		}
		if channel.RSSURL == "" {
			channel.RSSURL = existing.RSSURL
		}
		if channel.AvatarURL == "" {
			channel.AvatarURL = existing.AvatarURL
		}
		if channel.LastVideoID == "" {
			channel.LastVideoID = existing.LastVideoID
		}

		_, err := db.Exec(`
			UPDATE youtube_channels
			SET name = ?,
				handle = ?,
				rss_url = ?,
				avatar_url = ?,
				last_video_id = ?,
				updated_at = ?
			WHERE id = ?
		`,
			channel.Name,
			channel.Handle,
			channel.RSSURL,
			channel.AvatarURL,
			channel.LastVideoID,
			channel.UpdatedAt,
			channel.ID,
		)

		return channel, err
	}

	if channel.ID == "" {
		channel.ID = helper.GenerateID("yo")
	}

	channel.CreatedAt = now
	channel.UpdatedAt = now

	_, err = db.Exec(`
		INSERT INTO youtube_channels (
			id,
			name,
			handle,
			rss_url,
			avatar_url,
			last_video_id,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		channel.ID,
		channel.Name,
		channel.Handle,
		channel.RSSURL,
		channel.AvatarURL,
		channel.LastVideoID,
		channel.CreatedAt,
		channel.UpdatedAt,
	)

	return channel, err
}

func findYoutubeChannel(db *sql.DB, handle, rssURL string) (YoutubeChannel, bool, error) {
	var channel YoutubeChannel

	row := db.QueryRow(`
		SELECT id,
			   name,
			   handle,
			   rss_url,
			   avatar_url,
			   last_video_id,
			   created_at,
			   updated_at
		FROM youtube_channels
		WHERE handle = ? OR rss_url = ?
		LIMIT 1
	`, handle, rssURL)

	err := row.Scan(
		&channel.ID,
		&channel.Name,
		&channel.Handle,
		&channel.RSSURL,
		&channel.AvatarURL,
		&channel.LastVideoID,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return channel, false, nil
	}

	if err != nil {
		return channel, false, err
	}

	return channel, true, nil
}

func YoutubeChannelExists(handle, rssURL string) bool {
	db, err := openYoutubeDB()
	if err != nil {
		return false
	}
	defer db.Close()

	_, found, err := findYoutubeChannel(db, handle, rssURL)
	if err != nil {
		return false
	}

	return found
}

func GetYoutubeChannels() ([]YoutubeChannel, error) {
	db, err := openYoutubeDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT id,
			   name,
			   handle,
			   rss_url,
			   avatar_url,
			   last_video_id,
			   created_at,
			   updated_at
		FROM youtube_channels
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []YoutubeChannel

	for rows.Next() {
		var channel YoutubeChannel

		err := rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.Handle,
			&channel.RSSURL,
			&channel.AvatarURL,
			&channel.LastVideoID,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		channels = append(channels, channel)
	}

	return channels, rows.Err()
}
