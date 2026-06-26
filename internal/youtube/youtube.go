// modexusBot/internal/youtube/youtube.go
package youtube

import (
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/modexusdev/feedbot/internal/scheduler"
	"github.com/modexusdev/feedbot/internal/storage"
)

// Feed represents the YouTube RSS feed structure.
type Feed struct {
	Entries []Entry `xml:"entry"`
}

// Entry represents a single video entry in a YouTube RSS feed.
type Entry struct {
	Title     string `xml:"title"`
	Link      Link   `xml:"link"`
	VideoID   string `xml:"videoId"`
	Published string `xml:"published"`
}

// Link represents the video link inside a YouTube RSS entry.
type Link struct {
	Href string `xml:"href,attr"`
}

// LatestVideo contains normalized video data used by FeedBot.
type LatestVideo struct {
	Title     string
	Link      string
	VideoID   string
	Published string
}

// GetFeedVideos fetches and parses videos from a YouTube RSS feed.
func GetFeedVideos(rssURL string) ([]LatestVideo, error) {
	resp, err := http.Get(rssURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube rss request failed: %s", resp.Status)
	}

	var feed Feed

	err = xml.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		return nil, err
	}

	var videos []LatestVideo

	for _, entry := range feed.Entries {
		videos = append(videos, LatestVideo{
			Title:     entry.Title,
			Link:      entry.Link.Href,
			VideoID:   entry.VideoID,
			Published: entry.Published,
		})
	}

	return videos, nil
}

// CheckAllChannels checks all saved YouTube channels for new videos.
func CheckAllChannels() {
	channels, err := storage.GetYoutubeChannels()
	if err != nil {
		fmt.Println("load youtube channels error:", err)
		return
	}

	for i := range channels {
		channel := &channels[i]

		videos, err := GetFeedVideos(channel.RSSURL)
		if err != nil {
			fmt.Println("youtube check error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(videos) == 0 {
			time.Sleep(4 * time.Second)
			continue
		}

		latestVideo := videos[0]

		if channel.LastVideoID == "" {
			if !IsShort(latestVideo) {
				pushYoutubeVideo(channel, latestVideo)
			}

			saveLatestVideoID(channel, latestVideo.VideoID)
			time.Sleep(5 * time.Second)
			continue
		}

		if channel.LastVideoID == latestVideo.VideoID {
			time.Sleep(5 * time.Second)
			continue
		}

		var newVideos []LatestVideo

		for _, video := range videos {
			if video.VideoID == channel.LastVideoID {
				break
			}

			if IsShort(video) {
				continue
			}

			newVideos = append(newVideos, video)
		}

		for i := len(newVideos) - 1; i >= 0; i-- {
			pushYoutubeVideo(channel, newVideos[i])
		}

		saveLatestVideoID(channel, latestVideo.VideoID)

		time.Sleep(5 * time.Second)
	}
}

func pushYoutubeVideo(channel *storage.YoutubeChannel, video LatestVideo) {

	date := video.Published

	t, err := time.Parse(time.RFC3339, video.Published)
	if err == nil {
		date = t.Format("02.01.2006 15:04")
	}

	scheduler.Push(scheduler.ScheduledMessage{
		SourceEmoji: "🎥",
		SourceName:  "YouTube",
		Message: fmt.Sprintf(
			"👤 <b>%s</b>\n\n"+
				"🎬 <b>Neues Video</b>\n\n"+
				"━━━━━━━━━━━━\n"+
				"%s\n"+
				"━━━━━━━━━━━━\n\n"+
				"🗓 %s\n\n"+
				"🔗 %s",
			html.EscapeString(channel.Name),
			html.EscapeString(video.Title),
			html.EscapeString(date),
			html.EscapeString(video.Link),
		),
	})
}

func saveLatestVideoID(channel *storage.YoutubeChannel, videoID string) {
	channel.LastVideoID = videoID
	channel.UpdatedAt = time.Now().Format(time.RFC3339)

	_, err := storage.SaveYoutubeChannel(*channel)
	if err != nil {
		fmt.Println("save youtube channel error:", err)
	}
}
