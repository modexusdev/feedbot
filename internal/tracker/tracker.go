package tracker

import (
	"time"

	"github.com/modexusdev/feedbot/internal/scheduler"
	"github.com/modexusdev/feedbot/internal/weather"
	"github.com/modexusdev/feedbot/internal/youtube"
)

// Watch starts all enabled background content trackers.
func Watch() {
	go scheduler.Check(
		scheduler.Schedule{
			Interval:  30 * time.Minute,
			QuietFrom: 23,
			QuietTo:   5,
		},
		youtube.CheckAllChannels,
	)
	// Push today's weather report at 6 AM and tomorrow's at 6 PM
	go scheduler.DailyAt(1, 04, weather.PushTodayReport)
	// Push tomorrow's weather report at 6 PM
	go scheduler.DailyAt(1, 05, weather.PushTomorrowReport)
}
