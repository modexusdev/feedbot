// modexusBot/internal/weather/getweather.go
package weather

import "github.com/modexusdev/feedbot/internal/reply"

func GetWeather(lat, lon float64, location string) (string, error) {
	data, err := FetchWeather(lat, lon)
	if err != nil {
		return "", err
	}

	msg := FormatWeatherMessage(location, data)

	return reply.WeatherFormat(msg), nil
}
