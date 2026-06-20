// modexusBot/internal/bot/youtube_flow.go
package bot

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/reply"
	"github.com/modexusdev/feedbot/internal/storage"
	"github.com/modexusdev/feedbot/internal/youtube"
)

func (b *Bot) handleYoutubeLink(chatID int64, text string) {
	b.waitingForYoutubeLink[chatID] = false

	channel, err := youtube.ExtractYoutubeChannel(text)
	if err != nil {
		b.sendMessage(chatID, reply.Format("❌", "Could not read YouTube channel."))
		return
	}
	if storage.YoutubeChannelExists(channel.Handle, channel.RSSURL) {
		b.sendMessage(
			chatID,
			reply.YoutubeAlreadyAddedFormat(channel),
		)
		return
	}
	b.pendingYoutubeChannel[chatID] = channel

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Yes", "youtube_add_yes"),
			tgbotapi.NewInlineKeyboardButtonData("❌ No", "youtube_add_no"),
		),
	)

	if channel.AvatarURL != "" {
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(channel.AvatarURL))
		photo.Caption = reply.YoutubeAddFormat(channel)
		photo.ParseMode = "HTML"
		photo.ReplyMarkup = markup

		if _, err := b.api.Send(photo); err != nil {
			log.Printf("failed to send photo: %v", err)
		}

		return
	}

	msg := tgbotapi.NewMessage(chatID, reply.YoutubeAddFormat(channel))
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = markup

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("failed to send message: %v", err)
	}
}

func (b *Bot) handleYoutubeAddConfirm(chatID int64) {
	channel, ok := b.pendingYoutubeChannel[chatID]
	if !ok {
		b.sendMessage(chatID, reply.Format("❌", "No pending YouTube channel found."))
		return
	}

	if _, err := storage.SaveYoutubeChannel(channel); err != nil {
		b.sendMessage(chatID, reply.Format("❌", "Could not save YouTube channel."))
		return
	}

	delete(b.pendingYoutubeChannel, chatID)

	b.sendMessage(chatID, reply.YoutubeFormat("✅ YouTube channel added."))
}
func (b *Bot) handleYoutubeAddCancel(chatID int64) {
	delete(b.pendingYoutubeChannel, chatID)

	b.sendMessage(chatID, reply.YoutubeFormat("❌ YouTube channel was not added."))
}
func (b *Bot) handleYoutubeList(chatID int64) {
	channels, err := storage.GetYoutubeChannels()
	if err != nil {
		b.sendMessage(chatID, reply.YoutubeFormat("❌ Could not load channels."))
		return
	}

	b.sendMessage(chatID, reply.YoutubeListFormat(channels))
}
func (b *Bot) handleYoutubeRemoveStart(chatID int64) {
	channels, err := storage.GetYoutubeChannels()
	if err != nil {
		b.sendMessage(chatID, reply.YoutubeFormat("❌ Could not load channels."))
		return
	}

	if len(channels) == 0 {
		b.sendMessage(chatID, reply.YoutubeFormat("No YouTube channels found."))
		return
	}

	b.waitingForYoutubeRemove[chatID] = true

	b.sendMessage(
		chatID,
		reply.YoutubeListFormat(channels)+"\n\nWrite the number you want to remove.",
	)
}

func (b *Bot) handleYoutubeRemoveNumber(chatID int64, text string) {
	b.waitingForYoutubeRemove[chatID] = false

	index, err := strconv.Atoi(text)
	if err != nil || index < 1 {
		b.sendMessage(chatID, reply.YoutubeFormat("❌ Invalid number."))
		return
	}

	channels, err := storage.GetYoutubeChannels()
	if err != nil {
		b.sendMessage(chatID, reply.YoutubeFormat("❌ Could not load channels."))
		return
	}

	if index > len(channels) {
		b.sendMessage(chatID, reply.YoutubeFormat("❌ Channel number not found."))
		return
	}

	channel := channels[index-1]

	if err := storage.DeleteYoutubeChannel(channel.ID); err != nil {
		b.sendMessage(chatID, reply.YoutubeFormat("❌ Could not remove channel."))
		return
	}

	b.sendMessage(chatID, reply.YoutubeFormat("✅ Removed:\n\n"+channel.Name))
}
