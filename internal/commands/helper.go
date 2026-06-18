// modexusBot/internal/commands/helper.go
package commands

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Parse(text string) Command {
	text = strings.TrimSpace(text)

	cmd := Command{}

	if strings.HasPrefix(text, "@") {
		cmd.IsCommand = true
	}

	text = strings.ToLower(text)
	text = strings.TrimPrefix(text, "@")

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
func BuildHelpText(services EnabledServices) string {
	var b strings.Builder

	b.WriteString("Available commands\n\n")

	for _, cmd := range AvailableCommands(services) {
		b.WriteString("• ")
		b.WriteString(cmd.Name)
		b.WriteString("\n")
	}

	return b.String()
}

func BuildKeyboard(services EnabledServices) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton

	for _, cmd := range AvailableCommands(services) {
		button := tgbotapi.NewKeyboardButton("@" + cmd.Name)
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(button))
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	return keyboard
}
