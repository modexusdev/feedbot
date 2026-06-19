// modexusBot/internal/scheduler/watcher.go
package scheduler

import "time"

func Watch() {
	go startDevYoutubeWatcher()
	go startDevGithubWatcher()
	go startDevNewHackerWatcher()
}

func startDevYoutubeWatcher() {
	time.Sleep(10 * time.Second)

	Push(ScheduledMessage{
		SourceEmoji: "🎥",
		SourceName:  "YouTube",
		Message:     "New video uploaded by Pix",
	})
}
func startDevGithubWatcher() {
	time.Sleep(14 * time.Second)

	Push(ScheduledMessage{
		SourceEmoji: "🎥",
		SourceName:  "GitHub",
		Message:     "New repository created",
	})
}
func startDevNewHackerWatcher() {
	time.Sleep(10 * time.Second)

	Push(ScheduledMessage{
		SourceEmoji: "🎥",
		SourceName:  "New Hacker",
		Message:     "new hacker joined",
	})
}
