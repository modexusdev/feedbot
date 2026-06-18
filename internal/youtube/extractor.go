// modexusBot/internal/youtube/extractor.go
package youtube

import (
	"path/filepath"

	"github.com/modexusdev/feedbot/internal/helper"
	"github.com/modexusdev/feedbot/internal/storage"
)

func ExtractYoutubeChannel(link string) (storage.YoutubeChannel, error) {
	link = NormalizeYoutubeLink(link)

	html, err := FetchHTML(link)
	if err != nil {
		return storage.YoutubeChannel{}, err
	}

	avatarURL := ExtractAvatarURL(html)
	avatarURL = NormalizeAvatarSize(avatarURL)

	rssURL := ExtractRSSOrExternalID(html)

	name, err := FetchChannelNameFromRSS(rssURL)
	if err != nil {
		return storage.YoutubeChannel{}, err
	}

	channel := storage.YoutubeChannel{
		Name:      name,
		Handle:    ExtractHandle(link),
		RSSURL:    rssURL,
		AvatarURL: avatarURL,
	}

	return channel, nil
}

func SaveYoutubeChannelWithAvatar(channel storage.YoutubeChannel) error {
	channel, err := storage.SaveYoutubeChannel(channel)
	if err != nil {
		return err
	}

	if channel.AvatarURL == "" {
		return nil
	}

	channel.AvatarPath = filepath.Join("data", "images", channel.ID+".jpg")

	if err := helper.DownloadFile(channel.AvatarURL, channel.AvatarPath); err != nil {
		return err
	}

	_, err = storage.SaveYoutubeChannel(channel)
	return err
}

func ExtractYoutubeLink(link string) error {
	channel, err := ExtractYoutubeChannel(link)
	if err != nil {
		return err
	}

	return SaveYoutubeChannelWithAvatar(channel)
}
