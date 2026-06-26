// modexusBot/internal/weather/getweather.go
package weather

import "github.com/modexusdev/feedbot/internal/reply"

func GetWeatherRaw(lat, lon float64, location string, dayOffset int) (string, error) {
	data, err := FetchWeather(lat, lon, dayOffset)
	if err != nil {
		return "", err
	}

	return FormatWeatherMessage(location, data, dayOffset), nil
}

func GetWeather(lat, lon float64, location string, dayOffset int) (string, error) {
	data, err := FetchWeather(lat, lon, dayOffset)
	if err != nil {
		return "", err
	}

	msg := FormatWeatherMessage(location, data, dayOffset)

	return reply.WeatherFormat(msg), nil
}
