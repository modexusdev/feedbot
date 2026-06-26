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

// Schedule contains the schedule for an automation notification.
type Schedule struct {
	Interval  time.Duration
	QuietFrom int
	QuietTo   int
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
func Check(schedule Schedule, checkFunc func()) {
	time.Sleep(4 * time.Second)

	ticker := time.NewTicker(schedule.Interval)
	defer ticker.Stop()

	for {
		if !schedule.IsQuietTime() {
			checkFunc()
		}

		<-ticker.C
	}
}

// IsQuietTime reports whether the current local time
// falls within the configured quiet period.
func (s Schedule) IsQuietTime() bool {
	if s.QuietFrom < 0 || s.QuietTo < 0 {
		return false
	}

	hour := time.Now().Hour()

	if s.QuietFrom > s.QuietTo {
		return hour >= s.QuietFrom || hour < s.QuietTo
	}

	return hour >= s.QuietFrom && hour < s.QuietTo
}

// DailyAt runs a function at a specific time every day.
func DailyAt(hour, minute int, checkFunc func()) {
	for {
		now := time.Now()

		next := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			hour,
			minute,
			0,
			0,
			now.Location(),
		)

		if !next.After(now) {
			next = next.Add(24 * time.Hour)
		}

		time.Sleep(time.Until(next))
		checkFunc()
	}
}
