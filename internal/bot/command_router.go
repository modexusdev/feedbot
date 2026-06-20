// modexusBot/internal/bot/command_router.go
package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/commands"
)

func (b *Bot) handleCommandAction(chatID int64, cmd commands.Command) bool {
	switch cmd.Name {
	case "youtube":
		return b.handleYoutubeCommand(chatID, cmd)

	default:
		return false
	}
}

func (b *Bot) handleYoutubeCommand(chatID int64, cmd commands.Command) bool {
	switch cmd.Action {
	case "add":
		b.waitingForYoutubeLink[chatID] = true
		return false

	case "list":
		b.handleYoutubeList(chatID)
		return true

	case "remove":
		b.handleYoutubeRemoveStart(chatID)
		return true
	}

	return false
}

func (b *Bot) sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	b.api.Send(msg)
}
