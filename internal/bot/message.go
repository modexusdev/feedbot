// modexusBot/internal/bot/message.go

package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/commands"
	"github.com/modexusdev/feedbot/internal/config"
)

// handleUpdate processes incoming Telegram updates and routes them
// to the appropriate command, callback, or workflow handler.
func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		b.handleCallback(update.CallbackQuery)
		return
	}

	if update.Message == nil {
		return
	}

	userID := fmt.Sprintf("%d", update.Message.From.ID)

	if !config.IsAllowedUser(userID, b.config.AllowedUserIDs) {
		log.Printf(
			"Unauthorized access: %d (%s)",
			update.Message.From.ID,
			update.Message.From.UserName,
		)
		return
	}

	chatID := update.Message.Chat.ID
	text := commands.NormalizeKeyboardText(update.Message.Text)

	if b.waitingForYoutubeLink[chatID] {
		b.handleYoutubeLink(chatID, text)
		return
	}

	if b.waitingForYoutubeRemove[chatID] {
		b.handleYoutubeRemoveNumber(chatID, text)
		return
	}

	cmd := commands.Parse(text)
	log.Printf("Command executed: %s", cmd.Name)

	response := commands.Handle(cmd, b.services)

	if handled := b.handleCommandAction(chatID, cmd); handled {
		return
	}

	if response == "" {
		return
	}

	msg := tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = tgbotapi.ModeHTML

	if cmd.Name == "help" {
		msg.ReplyMarkup = commands.BuildMainKeyboard(b.services)
	}

	if cmd.Name == "youtube" && cmd.Action == "" {
		msg.ReplyMarkup = commands.BuildYoutubeKeyboard()
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("failed to send message: %v", err)
	}
}
