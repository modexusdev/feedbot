// modexusBot/internal/weather/formatter.go
package weather

import (
	"fmt"
	"strings"
	"time"

	"github.com/modexusdev/feedbot/internal/i18n"
)

// formatHourlyLine returns the formatted hourly line for the given weather item
func formatHourlyLine(h WeatherHourlyItem) string {
	parts := []string{
		fmt.Sprintf(
			"%s %s %.1f°C",
			timeIcon(h.Time),
			formatTime(h.Time),
			h.Temperature,
		),
	}

	if h.CloudCover > 0 {
		parts = append(parts, fmt.Sprintf("☁️ %d%%", h.CloudCover))
	}

	if h.Rain > 0 {
		parts = append(parts, fmt.Sprintf("🌧️ %.1f mm", h.Rain))
	}

	if h.Snowfall > 0 {
		parts = append(parts, fmt.Sprintf("❄️ %.1f", h.Snowfall))
	}

	return strings.Join(parts, " • ")
}

// formatDate returns the date part of the given timestamp string
func formatDate(value string) string {
	layout := "2006-01-02"

	if strings.Contains(value, "T") {
		layout = "2006-01-02T15:04"
	}

	t, err := time.Parse(layout, value)
	if err != nil {
		return value
	}

	return t.Format("02.01.2006")
}

// formatTime returns the time part of the given timestamp string
func formatTime(value string) string {
	t, err := time.Parse("2006-01-02T15:04", value)
	if err != nil {
		return value
	}

	return t.Format("15:04")
}

// weatherStatus returns the status of the weather based on the given code
func weatherStatus(code int) string {
	switch code {
	case 0:
		return i18n.T("weather.status.clear")
	case 1:
		return i18n.T("weather.status.mainly_clear")
	case 2:
		return i18n.T("weather.status.partly_cloudy")
	case 3:
		return i18n.T("weather.status.cloudy")
	case 45, 48:
		return i18n.T("weather.status.fog")
	case 51, 53, 55:
		return i18n.T("weather.status.drizzle")
	case 61, 63, 65:
		return i18n.T("weather.status.rain")
	case 80, 81, 82:
		return i18n.T("weather.status.rain_showers")
	case 71, 73, 75:
		return i18n.T("weather.status.snow")
	case 95, 96, 99:
		return i18n.T("weather.status.thunderstorm")
	default:
		return i18n.T("weather.status.unknown")
	}
}

// timeIcon returns the appropriate icon for the given time
func timeIcon(value string) string {
	t, err := time.Parse("2006-01-02T15:04", value)
	if err != nil {
		return "🕘"
	}

	hour := t.Hour()

	switch {
	case hour >= 5 && hour < 9:
		return "🌄"
	case hour >= 9 && hour < 12:
		return "🌞"
	case hour >= 12 && hour < 17:
		return "🏙️"
	case hour >= 17 && hour < 21:
		return "🌆"
	default:
		return "🌙"
	}
}

// feelsLikeIcon returns the appropriate icon for the perceived temperature difference.
func feelsLikeIcon(actual, feels float64) string {
	diff := feels - actual

	switch {
	case diff >= 4:
		return "🥵"
	case diff >= 2:
		return "🔥"
	case diff <= -4:
		return "🥶"
	case diff <= -2:
		return "🧊"
	default:
		return "😊"
	}
}
func FormatWeatherOverview(location string, data WeatherResponse, dayOffset int) string {
	var b strings.Builder

	fmt.Fprintf(&b, "📍 %s\n", location)
	fmt.Fprintf(&b, "📅 %s\n\n", formatDate(data.Daily.Date))

	if dayOffset == 0 {
		fmt.Fprintf(&b, "🌡️ %s: %.1f°C\n", i18n.T("weather.current"), data.Current.CurrentTemperature)

		fmt.Fprintf(
			&b,
			"%s %s: %.1f°C\n",
			feelsLikeIcon(data.Current.CurrentTemperature, data.Current.FeelsLike),
			i18n.T("weather.feels_like"),
			data.Current.FeelsLike,
		)

		fmt.Fprintf(&b, "☁️ %s: %d%%\n", i18n.T("weather.cloud_cover"), data.Current.CloudCover)
		fmt.Fprintf(&b, "🌤️ %s: %s\n\n", i18n.T("weather.status"), weatherStatus(data.Current.WeatherCode))
	}

	fmt.Fprintf(&b, "🔻 %s: %.1f°C\n", i18n.T("weather.min"), data.Daily.MinTemperature)
	fmt.Fprintf(&b, "🔺 %s: %.1f°C\n\n", i18n.T("weather.max"), data.Daily.MaxTemperature)

	fmt.Fprintf(&b, "🌄 %s: %s\n", i18n.T("weather.sunrise"), formatTime(data.Daily.Sunrise))
	fmt.Fprintf(&b, "🌆 %s: %s\n", i18n.T("weather.sunset"), formatTime(data.Daily.Sunset))

	return b.String()
}

func FormatWeatherForecast(location string, data WeatherResponse, dayOffset int) string {
	var b strings.Builder

	fmt.Fprintf(&b, "📍 %s\n", location)
	fmt.Fprintf(&b, "📅 %s\n\n", formatDate(data.Daily.Date))

	fmt.Fprintf(&b, "🕘 %s:\n\n", i18n.T("weather.forecast"))

	for _, h := range data.Hourly {
		b.WriteString(formatHourlyLine(h))
		b.WriteString("\n\n")
	}

	return strings.TrimSpace(b.String())
}
