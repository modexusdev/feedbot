// modexusBot/internal/commands/helper.go
package commands

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/i18n"
)

type KeyboardButtonConfig struct {
	Text    string
	Command string
}

func buttonYoutube() string {
	return "🎥 " + i18n.T("button.youtube")
}

func buttonWeather() string {
	return "🌤 " + i18n.T("button.weather")
}

func buttonHelp() string {
	return "📚 " + i18n.T("button.help")
}

func buttonYoutubeAdd() string {
	return "➕ " + i18n.T("button.add")
}

func buttonYoutubeList() string {
	return "📋 " + i18n.T("button.list")
}

func buttonYoutubeRemove() string {
	return "➖ " + i18n.T("button.remove")
}

func buttonYoutubeCheck() string {
	return "🔄 " + i18n.T("button.check")
}

func buttonWeatherToday() string {
	return "🌤 " + i18n.T("button.today")
}

func buttonWeatherTomorrow() string {
	return "🌥 " + i18n.T("button.tomorrow")
}

func buttonWeatherSetLocation() string {
	return "📍 " + i18n.T("button.set_location")
}

func buttonBack() string {
	return "🔙 " + i18n.T("button.back")
}

// NormalizeKeyboardText converts keyboard button text into real commands.
func NormalizeKeyboardText(text string) string {
	text = strings.TrimSpace(text)

	switch text {
	case buttonYoutube():
		return "#youtube"
	case buttonWeather():
		return "#weather"
	case buttonHelp():
		return "#help"
	case buttonYoutubeAdd():
		return "#youtube add"
	case buttonYoutubeList():
		return "#youtube list"
	case buttonYoutubeRemove():
		return "#youtube remove"
	case buttonYoutubeCheck():
		return "#youtube check"
	case buttonWeatherToday():
		return "#weather today"
	case buttonWeatherTomorrow():
		return "#weather tomorrow"
	case buttonWeatherSetLocation():
		return "#weather location"
	case buttonBack():
		return "#help"
	default:
		return text
	}
}

// Parse converts a text message into a structured command.
func Parse(text string) Command {
	text = strings.TrimSpace(text)

	cmd := Command{}

	if strings.HasPrefix(text, "#") {
		cmd.IsCommand = true
	}

	text = strings.ToLower(text)
	text = strings.TrimPrefix(text, "#")

	parts := strings.Fields(text)

	if len(parts) > 0 {
		cmd.Name = parts[0]
	}

	if len(parts) > 1 {
		cmd.Action = parts[1]
	}

	if len(parts) > 2 {
		cmd.Args = parts[2:]
	}

	return cmd
}

// BuildHelpText generates the help text for all enabled commands.
func BuildHelpText(services EnabledServices) string {
	var b strings.Builder

	b.WriteString(i18n.T("help.choose_service"))
	b.WriteString("\n\n")

	if services.Youtube {
		b.WriteString(buttonYoutube())
		b.WriteString("\n")
	}

	if services.Weather {
		b.WriteString(buttonWeather())
		b.WriteString("\n")
	}

	b.WriteString(buttonHelp())
	b.WriteString("\n")

	return b.String()
}

// BuildKeyboard creates the main Telegram keyboard.
func BuildKeyboard(services EnabledServices) tgbotapi.ReplyKeyboardMarkup {
	return BuildMainKeyboard(services)
}

func BuildMainKeyboard(services EnabledServices) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton

	if services.Youtube {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttonYoutube()),
		))
	}

	if services.Weather {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttonWeather()),
		))
	}

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(buttonHelp()),
	))

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	return keyboard
}

func BuildModuleKeyboard(buttons ...string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton

	for i := 0; i < len(buttons); i += 2 {
		row := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(buttons[i]),
		}

		if i+1 < len(buttons) {
			row = append(row, tgbotapi.NewKeyboardButton(buttons[i+1]))
		}

		rows = append(rows, row)
	}

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(buttonBack()),
	))

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	return keyboard
}

// BuildYoutubeKeyboard returns a keyboard for the youtube module.
func BuildYoutubeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return BuildModuleKeyboard(
		buttonYoutubeAdd(),
		buttonYoutubeList(),
		buttonYoutubeRemove(),
		buttonYoutubeCheck(),
	)
}

// BuildWeatherKeyboard returns a keyboard for the weather module.
func BuildWeatherKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return BuildModuleKeyboard(
		buttonWeatherToday(),
		buttonWeatherTomorrow(),
		buttonWeatherSetLocation(),
	)
}
