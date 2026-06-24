// modexusBot/internal/weather/getweather.go
package weather

import "github.com/modexusdev/feedbot/internal/reply"

func GetWeather(lat, lon float64, location string, dayOffset int) (string, error) {
	data, err := FetchWeather(lat, lon, dayOffset)
	if err != nil {
		return "", err
	}

	msg := FormatWeatherMessage(location, data, dayOffset)

	return reply.WeatherFormat(msg), nil
}
