package tracker

import (
	"time"

	"github.com/modexusdev/feedbot/internal/scheduler"
	"github.com/modexusdev/feedbot/internal/youtube"
)

func Watch() {
	go scheduler.Check(30*time.Minute, youtube.CheckAllChannels)
}

func startDevYoutubeWatcher() {
	time.Sleep(10 * time.Second)

	scheduler.Push(scheduler.ScheduledMessage{
		SourceEmoji: "🎥",
		SourceName:  "YouTube",
		Message:     "New video uploaded by Pix",
	})
}
func startDevGithubWatcher() {
	time.Sleep(14 * time.Second)

	scheduler.Push(scheduler.ScheduledMessage{
		SourceEmoji: "🎥",
		SourceName:  "GitHub",
		Message:     "New repository created",
	})
}
func startDevNewHackerWatcher() {
	time.Sleep(10 * time.Second)

	scheduler.Push(scheduler.ScheduledMessage{
		SourceEmoji: "🎥",
		SourceName:  "New Hacker",
		Message:     "new hacker joined",
	})
}
