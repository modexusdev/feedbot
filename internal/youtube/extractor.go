// modexusBot/internal/youtube/extractor.go
package youtube

import (
	"github.com/modexusdev/feedbot/internal/storage"
)

// ExtractYoutubeChannel extracts channel metadata from a YouTube link or handle.
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

// ExtractYoutubeLink extracts and saves a YouTube channel from a link or handle.
func ExtractYoutubeLink(link string) error {
	channel, err := ExtractYoutubeChannel(link)
	if err != nil {
		return err
	}

	_, err = storage.SaveYoutubeChannel(channel)
	return err
}
