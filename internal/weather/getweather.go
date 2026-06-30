// modexusBot/internal/weather/getweather.go
package weather

import "github.com/modexusdev/feedbot/internal/reply"

func GetWeatherPage(lat, lon float64, location string, dayOffset int, page string) (string, error) {
	data, err := FetchWeather(lat, lon, dayOffset)
	if err != nil {
		return "", err
	}

	switch page {
	case "forecast":
		return reply.WeatherFormat(FormatWeatherForecast(location, data, dayOffset)), nil
	default:
		return reply.WeatherFormat(FormatWeatherOverview(location, data, dayOffset)), nil
	}
}
