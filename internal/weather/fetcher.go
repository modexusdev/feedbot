// modexusBot/internal/weather/fetcher.go
package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`

	Current WeatherCurrent      `json:"current"`
	Hourly  []WeatherHourlyItem `json:"hourly"`
	Daily   WeatherDaily        `json:"daily"`
}

type WeatherCurrent struct {
	Time               string  `json:"time"`
	CurrentTemperature float64 `json:"current_temperature"`
	FeelsLike          float64 `json:"feels_like"`
	Rain               float64 `json:"rain"`
	Snowfall           float64 `json:"snowfall"`
	WeatherCode        int     `json:"weather_code"`
	CloudCover         int     `json:"cloud_cover"`
	IsDay              int     `json:"is_day"`
}

type WeatherHourlyItem struct {
	Time        string  `json:"time"`
	Temperature float64 `json:"temperature"`
	Rain        float64 `json:"rain"`
	Snowfall    float64 `json:"snowfall"`
	WeatherCode int     `json:"weather_code"`
	CloudCover  int     `json:"cloud_cover"`
	IsDay       int     `json:"is_day"`
}

type WeatherDaily struct {
	Date           string  `json:"date"`
	MaxTemperature float64 `json:"max_temperature"`
	MinTemperature float64 `json:"min_temperature"`
	Sunrise        string  `json:"sunrise"`
	Sunset         string  `json:"sunset"`
}

type openMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`

	Current openMeteoCurrent `json:"current"`
	Hourly  openMeteoHourly  `json:"hourly"`
	Daily   openMeteoDaily   `json:"daily"`
}

type openMeteoCurrent struct {
	Time                string  `json:"time"`
	Temperature         float64 `json:"temperature_2m"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	Rain                float64 `json:"rain"`
	Snowfall            float64 `json:"snowfall"`
	WeatherCode         int     `json:"weather_code"`
	CloudCover          int     `json:"cloud_cover"`
	IsDay               int     `json:"is_day"`
}

type openMeteoHourly struct {
	Time        []string  `json:"time"`
	Temperature []float64 `json:"temperature_2m"`
	Rain        []float64 `json:"rain"`
	Snowfall    []float64 `json:"snowfall"`
	WeatherCode []int     `json:"weather_code"`
	CloudCover  []int     `json:"cloud_cover"`
	IsDay       []int     `json:"is_day"`
}

type openMeteoDaily struct {
	Time           []string  `json:"time"`
	TemperatureMax []float64 `json:"temperature_2m_max"`
	TemperatureMin []float64 `json:"temperature_2m_min"`
	Sunrise        []string  `json:"sunrise"`
	Sunset         []string  `json:"sunset"`
}

// FetchWeather fetches the weather forecast for the given latitude, longitude, and day offset
func FetchWeather(lat, lon float64, dayOffset int) (WeatherResponse, error) {
	forecastDays := dayOffset + 1

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&timezone=auto&forecast_days=%d&current=temperature_2m,apparent_temperature,rain,snowfall,weather_code,cloud_cover,is_day&hourly=temperature_2m,rain,snowfall,weather_code,cloud_cover,is_day&daily=temperature_2m_max,temperature_2m_min,sunrise,sunset",
		lat,
		lon,
		forecastDays,
	)

	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherResponse{}, fmt.Errorf("weather api returned status: %d", resp.StatusCode)
	}

	var raw openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return WeatherResponse{}, err
	}

	return filterWeather(raw, dayOffset), nil
}

// filterWeather filters the raw weather response into a structured WeatherResponse
func filterWeather(raw openMeteoResponse, dayOffset int) WeatherResponse {
	return WeatherResponse{
		Latitude:  raw.Latitude,
		Longitude: raw.Longitude,
		Timezone:  raw.Timezone,
		Current:   filterCurrent(raw.Current),
		Hourly:    filterHourly(raw.Hourly, raw.Daily.Time[dayOffset]),
		Daily:     filterDaily(raw.Daily, dayOffset),
	}
}

// filterCurrent filters the raw current weather data into a structured WeatherCurrent
func filterCurrent(current openMeteoCurrent) WeatherCurrent {
	return WeatherCurrent{
		Time:               current.Time,
		CurrentTemperature: current.Temperature,
		FeelsLike:          current.ApparentTemperature,
		Rain:               current.Rain,
		Snowfall:           current.Snowfall,
		WeatherCode:        current.WeatherCode,
		CloudCover:         current.CloudCover,
		IsDay:              current.IsDay,
	}
}

// filterHourly filters the raw hourly weather data into a structured []WeatherHourlyItem
func filterHourly(hourly openMeteoHourly, targetDate string) []WeatherHourlyItem {
	wantedHours := map[string]bool{
		"06:00": true,
		"09:00": true,
		"12:00": true,
		"15:00": true,
		"18:00": true,
		"21:00": true,
		"23:00": true,
	}

	var result []WeatherHourlyItem

	for i, t := range hourly.Time {
		date := getDate(t)
		hour := getHour(t)

		if date != targetDate {
			continue
		}

		if !wantedHours[hour] {
			continue
		}

		result = append(result, WeatherHourlyItem{
			Time:        t,
			Temperature: hourly.Temperature[i],
			Rain:        hourly.Rain[i],
			Snowfall:    hourly.Snowfall[i],
			WeatherCode: hourly.WeatherCode[i],
			CloudCover:  hourly.CloudCover[i],
			IsDay:       hourly.IsDay[i],
		})
	}

	return result
}
func filterDaily(daily openMeteoDaily, dayOffset int) WeatherDaily {
	if len(daily.Time) <= dayOffset {
		return WeatherDaily{}
	}

	return WeatherDaily{
		Date:           daily.Time[dayOffset],
		MaxTemperature: daily.TemperatureMax[dayOffset],
		MinTemperature: daily.TemperatureMin[dayOffset],
		Sunrise:        daily.Sunrise[dayOffset],
		Sunset:         daily.Sunset[dayOffset],
	}
}

// getDate extracts the date part from a timestamp string
func getDate(value string) string {
	parts := strings.Split(value, "T")
	if len(parts) != 2 {
		return ""
	}

	return parts[0]
}
func getHour(value string) string {
	parts := strings.Split(value, "T")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}
