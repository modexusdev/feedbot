// modexusBot/internal/weather/weather.go
package weather

import (
	"fmt"

	"github.com/modexusdev/feedbot/internal/scheduler"
	"github.com/modexusdev/feedbot/internal/storage"
)

// PushTodayReport pushes the weather report for today to the scheduler.
func PushTodayReport() {
	pushWeatherReport("🌅 Morgenbericht für heute", 0)
}

// PushTomorrowReport pushes the weather report for tomorrow to the scheduler.
func PushTomorrowReport() {
	pushWeatherReport("🌆 Abendbericht für morgen", 1)
}

// pushWeatherReport pushes the weather report for a given day offset to the scheduler.
func pushWeatherReport(title string, dayOffset int) {
	location, err := storage.GetWeatherLocation()
	if err != nil {
		location = storage.WeatherLocation{
			ID:        "default",
			Name:      "Berlin",
			Country:   "Deutschland",
			Latitude:  52.5173885,
			Longitude: 13.3951309,
		}
	}

	msg, err := GetWeatherRaw(
		location.Latitude,
		location.Longitude,
		location.Name,
		dayOffset,
	)

	if err != nil {
		fmt.Println("weather automation error:", err)
		return
	}

	scheduler.Push(scheduler.ScheduledMessage{
		SourceEmoji: "🌦️",
		SourceName:  "Weather",
		Message:     "<b>" + title + "</b>\n\n" + msg,
	})
}
