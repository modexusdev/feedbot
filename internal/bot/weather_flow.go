// modexusBot/internal/bot/weather_flow.go
package bot

import (
	"log"

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
			b.sendMessage(
				chatID,
				reply.Format("❌", "Could not load weather data."),
			)
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
			b.sendMessage(
				chatID,
				reply.Format("❌", "Could not load weather data."),
			)
			return true
		}

		b.sendMessage(chatID, msg)
		return true

	case "warnings":
		b.sendMessage(chatID, reply.Format("⚠️", "Weather warnings will be shown here."))
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
