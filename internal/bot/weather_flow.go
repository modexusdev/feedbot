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
	"github.com/modexusdev/feedbot/internal/weather"
)

type WeatherLocation struct {
	Name string
	Lat  float64
	Lon  float64
}

var weatherLocation = WeatherLocation{
	Name: "Halle (Saale)",
	Lat:  51.48158,
	Lon:  11.97947,
}

func (b *Bot) handleWeatherCommand(chatID int64, cmd commands.Command) bool {
	if cmd.Action == "" {
		b.sendWeatherMenu(chatID)
		return true
	}

	switch cmd.Action {
	case "today":
		msg, err := weather.GetWeather(
			weatherLocation.Lat,
			weatherLocation.Lon,
			weatherLocation.Name,
			0,
		)

		if err != nil {
			b.sendMessage(chatID, reply.Format("❌", "Could not load weather data."))
			return true
		}

		b.sendMessage(chatID, msg)
		return true

	case "tomorrow":
		msg, err := weather.GetWeather(
			weatherLocation.Lat,
			weatherLocation.Lon,
			weatherLocation.Name,
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

	weatherLocation = WeatherLocation{
		Name: loc.Name,
		Lat:  loc.Latitude,
		Lon:  loc.Longitude,
	}

	delete(b.pendingWeatherLocations, chatID)
	delete(b.waitingForWeatherLocationNumber, chatID)

	b.sendMessage(
		chatID,
		reply.Format("✅", weather.FormatSelectedLocation(loc)),
	)
}
