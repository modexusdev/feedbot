package tracker

import (
	"time"

	"github.com/modexusdev/feedbot/internal/scheduler"
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
}
