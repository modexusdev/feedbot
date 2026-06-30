// modexusBot/internal/bot/weather_flow.go
package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/commands"
	"github.com/modexusdev/feedbot/internal/i18n"
	"github.com/modexusdev/feedbot/internal/reply"
	"github.com/modexusdev/feedbot/internal/storage"
	"github.com/modexusdev/feedbot/internal/weather"
)

func (b *Bot) handleWeatherCommand(chatID int64, cmd commands.Command) bool {
	if cmd.Action == "" {
		b.sendWeatherMenu(chatID)
		return true
	}

	switch cmd.Action {
	case "today":
		location := getWeatherLocation()
		msg, err := weather.GetWeather(location.Latitude, location.Longitude, location.Name, 0)
		if err != nil {
			b.sendMessage(chatID, reply.Format("❌", "Weather", i18n.T("weather.load_data_error")))
			return true
		}

		b.sendMessage(chatID, msg)
		return true

	case "tomorrow":
		location := getWeatherLocation()
		msg, err := weather.GetWeather(location.Latitude, location.Longitude, location.Name, 1)
		if err != nil {
			b.sendMessage(chatID, reply.Format("❌", "Weather", i18n.T("weather.load_data_error")))
			return true
		}

		b.sendMessage(chatID, msg)
		return true

	case "location":
		b.waitingForWeatherLocation[chatID] = true

		b.sendMessage(chatID, reply.Format("📍", "Weather", i18n.T("weather.send_city_name")))
		return true
	}

	return false
}

func (b *Bot) sendWeatherMenu(chatID int64) {
	msg := tgbotapi.NewMessage(
		chatID,
		reply.Format("🌤", "Weather", i18n.T("weather.choose_action")),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = commands.BuildWeatherKeyboard()

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("failed to send weather menu: %v", err)
	}
}

func (b *Bot) handleWeatherLocation(chatID int64, city string) {
	delete(b.waitingForWeatherLocation, chatID)

	locations, err := weather.GeocodeCity(city)
	if err != nil {
		b.sendMessage(chatID, reply.Format("❌", "Weather", i18n.T("weather.city_not_found")))
		return
	}

	if len(locations) == 0 {
		b.sendMessage(chatID, reply.Format("❌", "Weather", i18n.T("weather.no_city_found")))
		return
	}

	b.pendingWeatherLocations[chatID] = locations
	b.waitingForWeatherLocationNumber[chatID] = true

	b.sendMessage(chatID, reply.Format("📍", "Weather", weather.FormatLocationList(locations)))
}

func (b *Bot) handleWeatherLocationNumber(chatID int64, text string) {
	locations := b.pendingWeatherLocations[chatID]

	number, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		b.sendMessage(chatID, reply.Format("❌", "Weather", i18n.T("weather.send_only_number")))
		return
	}

	if number < 1 || number > len(locations) {
		b.sendMessage(
			chatID,
			reply.Format(
				"❌",
				"Weather",
				fmt.Sprintf(i18n.T("weather.choose_number_between"), len(locations)),
			),
		)
		return
	}

	loc := locations[number-1]

	_, err = storage.SaveWeatherLocation(storage.WeatherLocation{
		Name:      loc.Name,
		Country:   loc.Country,
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
	})

	if err != nil {
		b.sendMessage(chatID, reply.Format("❌", "Weather", i18n.T("weather.save_location_error")))
		return
	}

	delete(b.pendingWeatherLocations, chatID)
	delete(b.waitingForWeatherLocationNumber, chatID)

	b.sendMessage(chatID, reply.Format("✅", "Weather", weather.FormatSelectedLocation(loc)))
}

func getWeatherLocation() storage.WeatherLocation {
	location, err := storage.GetWeatherLocation()
	if err == nil {
		return location
	}

	return storage.WeatherLocation{
		ID:        "default",
		Name:      "Berlin",
		Country:   "Deutschland",
		Latitude:  52.5173885,
		Longitude: 13.3951309,
	}
}
