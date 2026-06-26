// modexusBot/internal/bot/weather_flow.go
package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/commands"
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
		msg, err := weather.GetWeather(
			location.Latitude,
			location.Longitude,
			location.Name,
			0,
		)

		if err != nil {
			b.sendMessage(chatID, reply.Format("❌", "Could not load weather data."))
			return true
		}

		b.sendMessage(chatID, msg)
		return true

	case "tomorrow":
		location := getWeatherLocation()
		msg, err := weather.GetWeather(
			location.Latitude,
			location.Longitude,
			location.Name,
			1,
		)

		if err != nil {
			b.sendMessage(chatID, reply.Format("❌", "Could not load weather data."))
			return true
		}

		b.sendMessage(chatID, msg)
		return true

	case "location":
		b.waitingForWeatherLocation[chatID] = true

		b.sendMessage(
			chatID,
			reply.Format("📍", "Bitte schick mir den Namen deiner Stadt.\n\nBeispiel: Halle"),
		)

		return true
	}

	return false
}

func (b *Bot) sendWeatherMenu(chatID int64) {
	msg := tgbotapi.NewMessage(
		chatID,
		reply.Format("🌤", "Choose a weather action."),
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
		b.sendMessage(chatID, reply.Format("❌", "Stadt konnte nicht gefunden werden."))
		return
	}

	if len(locations) == 0 {
		b.sendMessage(chatID, reply.Format("❌", "Keine Stadt gefunden."))
		return
	}

	b.pendingWeatherLocations[chatID] = locations
	b.waitingForWeatherLocationNumber[chatID] = true

	b.sendMessage(
		chatID,
		reply.Format("📍", weather.FormatLocationList(locations)),
	)
}
func (b *Bot) handleWeatherLocationNumber(chatID int64, text string) {
	locations := b.pendingWeatherLocations[chatID]

	number, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		b.sendMessage(chatID, reply.Format("❌", "Bitte sende nur eine Nummer aus der Liste."))
		return
	}

	if number < 1 || number > len(locations) {
		b.sendMessage(
			chatID,
			reply.Format("❌", fmt.Sprintf("Bitte wähle eine Nummer zwischen 1 und %d.", len(locations))),
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
		b.sendMessage(chatID, reply.Format("❌", "Standort konnte nicht gespeichert werden."))
		return
	}

	delete(b.pendingWeatherLocations, chatID)
	delete(b.waitingForWeatherLocationNumber, chatID)

	b.sendMessage(
		chatID,
		reply.Format("✅", weather.FormatSelectedLocation(loc)),
	)
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
