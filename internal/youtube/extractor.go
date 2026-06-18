// modexusBot/internal/youtube/extractor.go
package youtube

import (
	"path/filepath"

	"github.com/modexusdev/feedbot/internal/helper"
	"github.com/modexusdev/feedbot/internal/storage"
)

func ExtractYoutubeLink(link string) error {
	link = NormalizeYoutubeLink(link)

	html, err := FetchHTML(link)
	if err != nil {
		return err
	}

	avatarURL := ExtractAvatarURL(html)
	avatarURL = NormalizeAvatarSize(avatarURL)
	rssURL := ExtractRSSOrExternalID(html)

	name, err := FetchChannelNameFromRSS(rssURL)
	if err != nil {
		return err
	}
	channel := storage.YoutubeChannel{
		Name:      name,
		Handle:    ExtractHandle(link),
		RSSURL:    ExtractRSSOrExternalID(html),
		AvatarURL: avatarURL,
	}

	channel, err = storage.SaveYoutubeChannel(channel)
	if err != nil {
		return err
	}

	if channel.AvatarURL != "" {
		channel.AvatarPath = filepath.Join("data", "images", channel.ID+".jpg")

		if err := helper.DownloadFile(channel.AvatarURL, channel.AvatarPath); err != nil {
			return err
		}

		if _, err := storage.SaveYoutubeChannel(channel); err != nil {
			return err
		}
	}

	return nil
}
