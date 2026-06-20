// modexusBot/internal/scheduler/scheduler.go
package scheduler

import (
	"time"

	"github.com/modexusdev/feedbot/internal/reply"
)

type ScheduledMessage struct {
	SourceEmoji string
	SourceName  string
	Message     string
}

type AutomationMessage struct {
	Text string
}

var Queue = make(chan AutomationMessage, 100)

func Push(msg ScheduledMessage) {
	text := reply.AutomationFormat(
		msg.SourceEmoji,
		msg.SourceName,
		msg.Message,
	)

	Queue <- AutomationMessage{
		Text: text,
	}
}
func Check(interval time.Duration, checkFunc func()) {
	time.Sleep(4 * time.Second)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		checkFunc()
		<-ticker.C
	}
}
