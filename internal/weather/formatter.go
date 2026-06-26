// modexusBot/internal/weather/formatter.go
package weather

import (
	"fmt"
	"strings"
	"time"
)

// FormatWeatherMessage returns the formatted weather message for the given location and weather data
func FormatWeatherMessage(location string, data WeatherResponse, dayOffset int) string {
	var b strings.Builder

	fmt.Fprintf(&b, "📍 %s\n", location)

	fmt.Fprintf(
		&b,
		"📅 %s\n\n",
		formatDate(data.Daily.Date),
	)

	if dayOffset == 0 {
		fmt.Fprintf(
			&b,
			"🌡️ Aktuell: %.1f°C \n",
			data.Current.CurrentTemperature,
		)

		fmt.Fprintf(
			&b,
			"%s Gefühlt: %.1f°C\n",
			feelsLikeIcon(
				data.Current.CurrentTemperature,
				data.Current.FeelsLike,
			),
			data.Current.FeelsLike,
		)

		fmt.Fprintf(
			&b,
			"☁️ Bewölkung: %d%%\n",
			data.Current.CloudCover,
		)

		fmt.Fprintf(
			&b,
			"🌤️ Status: %s\n\n",
			weatherStatus(data.Current.WeatherCode),
		)
	}

	fmt.Fprintf(
		&b,
		"🔻 Min: %.1f°C\n",
		data.Daily.MinTemperature,
	)

	fmt.Fprintf(
		&b,
		"🔺 Max: %.1f°C\n\n",
		data.Daily.MaxTemperature,
	)

	fmt.Fprintf(
		&b,
		"🌄 Sunrise: %s\n",
		formatTime(data.Daily.Sunrise),
	)

	fmt.Fprintf(
		&b,
		"🌆 Sunset: %s\n\n",
		formatTime(data.Daily.Sunset),
	)

	b.WriteString("🕘 Forecast:\n\n")

	for _, h := range data.Hourly {
		b.WriteString(formatHourlyLine(h))
		b.WriteString("\n\n")
	}

	return b.String()
}

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

// weatherIcon returns the appropriate icon for the given weather code
func weatherIcon(code int) string {
	switch code {
	case 0:
		return "☀️"
	case 1, 2:
		return "🌤️"
	case 3:
		return "☁️"
	case 45, 48:
		return "🌫️"
	case 51, 53, 55, 61, 63, 65, 80, 81, 82:
		return "🌧️"
	case 71, 73, 75, 77, 85, 86:
		return "❄️"
	case 95, 96, 99:
		return "⛈️"
	default:
		return "🌦️"
	}
}

// weatherStatus returns the status of the weather based on the given code
func weatherStatus(code int) string {
	switch code {
	case 0:
		return "Klar"
	case 1:
		return "Überwiegend klar"
	case 2:
		return "Teilweise bewölkt"
	case 3:
		return "Bewölkt"
	case 45, 48:
		return "Nebel"
	case 51, 53, 55:
		return "Nieselregen"
	case 61, 63, 65:
		return "Regen"
	case 80, 81, 82:
		return "Regenschauer"
	case 71, 73, 75:
		return "Schnee"
	case 95, 96, 99:
		return "Gewitter"
	default:
		return "Unbekannt"
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

// timeIcon returns the appropriate icon for the given time
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
