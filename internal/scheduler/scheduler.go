// modexusBot/internal/scheduler/scheduler.go
package scheduler

import (
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
