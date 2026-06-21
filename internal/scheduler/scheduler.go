// modexusBot/internal/scheduler/scheduler.go
package scheduler

import (
	"time"

	"github.com/modexusdev/feedbot/internal/reply"
)

// ScheduledMessage contains the raw data for an automation notification.
type ScheduledMessage struct {
	SourceEmoji string
	SourceName  string
	Message     string
}

// AutomationMessage contains the final formatted message text.
type AutomationMessage struct {
	Text string
}

// Queue stores automation messages before they are sent by the bot.
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

// Check runs a function immediately after a short startup delay
// and then repeatedly at the given interval.
func Check(interval time.Duration, checkFunc func()) {
	time.Sleep(4 * time.Second)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		checkFunc()
		<-ticker.C
	}
}
