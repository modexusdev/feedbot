// modexusBot/internal/weather/weather.go
package weather

import (
	"fmt"

	"github.com/modexusdev/feedbot/internal/commands"
	"github.com/modexusdev/feedbot/internal/i18n"
	"github.com/modexusdev/feedbot/internal/scheduler"
	"github.com/modexusdev/feedbot/internal/storage"
)

func PushTodayReport() {
	pushWeatherReport(i18n.T("weather.today_report"), 0)
}

func PushTomorrowReport() {
	pushWeatherReport(i18n.T("weather.tomorrow_report"), 1)
}

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

	data, err := FetchWeather(location.Latitude, location.Longitude, dayOffset)
	if err != nil {
		fmt.Println("weather automation error:", err)
		return
	}

	overview := FormatWeatherOverview(location.Name, data, dayOffset)
	keyboard := commands.BuildWeatherPageKeyboard(dayOffset, "overview")

	scheduler.Push(scheduler.ScheduledMessage{
		SourceEmoji: "🌦️",
		SourceName:  "Weather",
		Message:     "<b>" + title + " • " + i18n.T("button.overview") + "</b>\n\n" + overview,
		ReplyMarkup: keyboard,
	})
}
