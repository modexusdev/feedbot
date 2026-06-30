// modexusBot/internal/bot/language_flow.go
package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/i18n"
	"github.com/modexusdev/feedbot/internal/reply"
	"github.com/modexusdev/feedbot/internal/storage"
)

func (b *Bot) sendLanguageMenu(chatID int64) {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, lang := range i18n.GetAvailableLanguages() {
		text := fmt.Sprintf("%s %s", lang.Flag, lang.Name)
		data := "language_set_" + lang.Code

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, data),
		))
	}

	msg := tgbotapi.NewMessage(
		chatID,
		reply.Format("🌍", "Language", i18n.T("language.choose")),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("failed to send language menu: %v", err)
	}
}

func (b *Bot) handleLanguageCallback(chatID int64, data string) bool {
	if !strings.HasPrefix(data, "language_set_") {
		return false
	}

	code := strings.TrimPrefix(data, "language_set_")

	if !i18n.IsSupported(code) {
		b.sendMessage(chatID, reply.Format("❌", "Language", i18n.T("language.unsupported")))
		return true
	}

	if err := storage.SaveLanguage(code); err != nil {
		b.sendMessage(chatID, reply.Format("❌", "Language", i18n.T("language.save_error")))
		return true
	}

	i18n.SetLanguage(code)

	b.sendMessage(
		chatID,
		reply.Format(
			"🌍",
			"Language",
			fmt.Sprintf(
				i18n.T("language.changed"),
				i18n.GetLanguageName(code),
			),
		),
	)

	return true
}
