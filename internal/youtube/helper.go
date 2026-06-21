// modexusBot/internal/youtube/helper.go
package youtube

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// NormalizeYoutubeLink normalizes a YouTube handle link to a canonical channel URL.
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

// ExtractRSSOrExternalID extracts the RSS feed URL or channel ID from YouTube page HTML.
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

// FetchHTML downloads and returns the HTML content of a page.
func FetchHTML(link string) (string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ExtractHandle extracts the YouTube handle from a channel URL.
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

// ExtractAvatarURL extracts the channel avatar URL from YouTube page HTML.
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

// NormalizeAvatarSize normalizes a YouTube avatar URL to a fixed thumbnail size.
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

// FetchChannelNameFromRSS fetches the channel name from a YouTube RSS feed.
func FetchChannelNameFromRSS(rssURL string) (string, error) {
	resp, err := http.Get(rssURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var feed youtubeRSSFeed

	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return "", err
	}

	return feed.Title, nil
}
