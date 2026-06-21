package tracker

import (
	"time"

	"github.com/modexusdev/feedbot/internal/scheduler"
	"github.com/modexusdev/feedbot/internal/youtube"
)

// Watch starts all enabled background content trackers.
func Watch() {
	go scheduler.Check(30*time.Minute, youtube.CheckAllChannels)
}
