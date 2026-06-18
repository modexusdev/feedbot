// modexusBot/internal/youtube/helper.go
package youtube

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NormalizeYoutubeLink(link string) string {
	link = strings.TrimSpace(link)

	if link == "" {
		return ""
	}

	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	parsed, err := url.Parse(link)
	if err != nil {
		return ""
	}

	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")

	for _, part := range parts {
		if strings.HasPrefix(part, "@") {
			return "https://www.youtube.com/" + part
		}
	}

	return ""
}

func ExtractRSSOrExternalID(html string) string {
	rssURL := extractBetween(html, `"rssUrl":"`, `"`)

	if rssURL != "" {
		return rssURL
	}

	externalID := extractBetween(html, `"externalId":"`, `"`)

	if externalID != "" {
		return "https://www.youtube.com/feeds/videos.xml?channel_id=" + externalID
	}

	return ""
}

func extractBetween(text string, startText string, endText string) string {
	start := strings.Index(text, startText)
	if start == -1 {
		return ""
	}

	start += len(startText)

	end := strings.Index(text[start:], endText)
	if end == -1 {
		return ""
	}

	return text[start : start+end]
}
func FetchHTML(link string) (string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func ExtractHandle(link string) string {
	link = NormalizeYoutubeLink(link)

	if link == "" {
		return ""
	}

	parsed, err := url.Parse(link)
	if err != nil {
		return ""
	}

	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")

	for _, part := range parts {
		if strings.HasPrefix(part, "@") {
			return strings.ToLower(part)
		}
	}

	return ""
}
func ExtractAvatarURL(html string) string {
	metadata := extractBetween(html, `"channelMetadataRenderer":{`, `,"metadataRowContainer"`)

	avatarURL := ""

	if metadata != "" {
		avatarURL = extractBetween(metadata, `"avatar":{"thumbnails":[{"url":"`, `"`)
	}

	if avatarURL == "" {
		avatarURL = extractBetween(html, `"channelMetadataRenderer":{`, `}`)
		avatarURL = extractBetween(avatarURL, `"avatar":{"thumbnails":[{"url":"`, `"`)
	}

	if avatarURL == "" {
		return ""
	}

	avatarURL = strings.ReplaceAll(avatarURL, `\/`, `/`)
	avatarURL = strings.ReplaceAll(avatarURL, `\u0026`, `&`)

	return avatarURL
}
func NormalizeAvatarSize(avatarURL string) string {
	if avatarURL == "" {
		return ""
	}

	parts := strings.Split(avatarURL, "=")
	if len(parts) < 2 {
		return avatarURL
	}

	return parts[0] + "=s128-c-k-c0x00ffffff-no-rj"
}

type youtubeRSSFeed struct {
	Title string `xml:"title"`
}

func FetchChannelNameFromRSS(rssURL string) (string, error) {
	resp, err := http.Get(rssURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var feed youtubeRSSFeed

	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return "", err
	}

	return feed.Title, nil
}
