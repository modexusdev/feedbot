// modexusBot/internal/commands/helper.go
package commands

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ButtonYoutube       = "🎥 YouTube"
	ButtonHelp          = "📚 Help"
	ButtonYoutubeAdd    = "➕ Add"
	ButtonYoutubeList   = "📋 List"
	ButtonYoutubeRemove = "➖ Remove"
	ButtonYoutubeCheck  = "🔄 Check"
	ButtonBack          = "🔙 Back"
)

// NormalizeKeyboardText converts keyboard button text into real commands.
func NormalizeKeyboardText(text string) string {
	text = strings.TrimSpace(text)

	switch text {
	case ButtonYoutube:
		return "#youtube"
	case ButtonHelp:
		return "#help"
	case ButtonYoutubeAdd:
		return "#youtube add"
	case ButtonYoutubeList:
		return "#youtube list"
	case ButtonYoutubeRemove:
		return "#youtube remove"
	case ButtonYoutubeCheck:
		return "#youtube check"
	case ButtonBack:
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

	b.WriteString("Choose a service\n\n")

	if services.Youtube {
		b.WriteString("🎥 YouTube\n")
	}

	b.WriteString("📚 Help\n")

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
			tgbotapi.NewKeyboardButton(ButtonYoutube),
		))
	}

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ButtonHelp),
	))

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	return keyboard
}

type KeyboardButtonConfig struct {
	Text    string
	Command string
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
		tgbotapi.NewKeyboardButton(ButtonBack),
	))

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	return keyboard
}

func BuildYoutubeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return BuildModuleKeyboard(
		ButtonYoutubeAdd,
		ButtonYoutubeList,
		ButtonYoutubeRemove,
		ButtonYoutubeCheck,
	)
}

// Example for later:
// const (
// 	ButtonWeather         = "🌦 Weather"
// 	ButtonWeatherToday    = "📍 Today"
// 	ButtonWeatherTomorrow = "📅 Tomorrow"
// 	ButtonWeatherWarning  = "⚠️ Warnings"
// )

// func BuildWeatherKeyboard() tgbotapi.ReplyKeyboardMarkup {
// 	return BuildModuleKeyboard(
// 		ButtonWeatherToday,
// 		ButtonWeatherTomorrow,
// 		ButtonWeatherWarning,
// 	)
// }
