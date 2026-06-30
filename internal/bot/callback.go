// modexusBot/internal/bot/callback.go
package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/config"
)

// handleCallback processes Telegram inline keyboard callbacks.
func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) {
	userID := fmt.Sprintf("%d", callback.From.ID)

	if !config.IsAllowedUser(userID, b.config.AllowedUserIDs) {
		return
	}

	if callback.Message == nil {
		return
	}

	chatID := callback.Message.Chat.ID

	answer := tgbotapi.NewCallback(callback.ID, "")
	if _, err := b.api.Request(answer); err != nil {
		log.Printf("failed to answer callback: %v", err)
	}

	// Language callbacks
	if b.handleLanguageCallback(chatID, callback.Data) {
		return
	}
	if b.handleWeatherCallback(chatID, callback.Message.MessageID, callback.Data) {
		return
	}

	switch callback.Data {
	case "youtube_add_yes":
		b.handleYoutubeAddConfirm(chatID)

	case "youtube_add_no":
		b.handleYoutubeAddCancel(chatID)
	}
}
