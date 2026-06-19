// modexusBot/internal/storage/youtube_store.go
package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/modexusdev/feedbot/internal/helper"
)

const youtubeDB = "data/youtube/channels.json"

func SaveYoutubeChannel(channel YoutubeChannel) (YoutubeChannel, error) {
	if err := os.MkdirAll(filepath.Dir(youtubeDB), 0755); err != nil {
		return channel, err
	}

	channels, err := loadYoutubeChannels()
	if err != nil {
		return channel, err
	}

	now := time.Now().Format(time.RFC3339)

	for i, existing := range channels {
		if existing.RSSURL == channel.RSSURL || existing.Handle == channel.Handle {
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

			channels[i] = channel
			return channel, saveYoutubeChannels(channels)
		}
	}

	if channel.ID == "" {
		channel.ID = generateUniqueYoutubeID(channels)
	}

	channel.CreatedAt = now
	channel.UpdatedAt = now

	channels = append(channels, channel)

	return channel, saveYoutubeChannels(channels)
}
func loadYoutubeChannels() ([]YoutubeChannel, error) {
	var channels []YoutubeChannel

	data, err := os.ReadFile(youtubeDB)
	if err != nil {
		if os.IsNotExist(err) {
			return channels, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return channels, nil
	}

	if err := json.Unmarshal(data, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

func saveYoutubeChannels(channels []YoutubeChannel) error {
	jsonData, err := json.MarshalIndent(channels, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(youtubeDB, jsonData, 0644)
}

func generateUniqueYoutubeID(channels []YoutubeChannel) string {
	for {
		id := helper.GenerateID("yo")

		exists := false
		for _, channel := range channels {
			if channel.ID == id {
				exists = true
				break
			}
		}

		if !exists {
			return id
		}
	}
}
func YoutubeChannelExists(handle, rssURL string) bool {
	channels, err := loadYoutubeChannels()
	if err != nil {
		return false
	}

	for _, channel := range channels {
		if rssURL != "" && channel.RSSURL == rssURL {
			return true
		}

		if handle != "" && channel.Handle == handle {
			return true
		}
	}

	return false
}
func GetYoutubeChannels() ([]YoutubeChannel, error) {
	return loadYoutubeChannels()
}
