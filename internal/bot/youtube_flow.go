// modexusBot/internal/bot/youtube_flow.go
package bot

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/commands"
	"github.com/modexusdev/feedbot/internal/reply"
	"github.com/modexusdev/feedbot/internal/storage"
	"github.com/modexusdev/feedbot/internal/youtube"
)

// handleYoutubeLink extracts channel data from a submitted YouTube link
// and asks the user to confirm before saving it.
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
		photo.ParseMode = tgbotapi.ModeHTML
		photo.ReplyMarkup = markup

		if _, err := b.api.Send(photo); err != nil {
			log.Printf("failed to send photo: %v", err)
		}

		return
	}

	msg := tgbotapi.NewMessage(chatID, reply.YoutubeAddFormat(channel))
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = markup

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("failed to send message: %v", err)
	}
}

// handleYoutubeAddConfirm saves the pending YouTube channel after user confirmation.
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

// handleYoutubeAddCancel cancels the pending YouTube channel addition.
func (b *Bot) handleYoutubeAddCancel(chatID int64) {
	delete(b.pendingYoutubeChannel, chatID)

	b.sendMessage(chatID, reply.YoutubeFormat("❌ YouTube channel was not added."))
}

// handleYoutubeList sends all saved YouTube channels to the user.
func (b *Bot) handleYoutubeList(chatID int64) {
	channels, err := storage.GetYoutubeChannels()
	if err != nil {
		b.sendMessage(chatID, reply.YoutubeFormat("❌ Could not load channels."))
		return
	}

	b.sendMessage(chatID, reply.YoutubeListFormat(channels))
}

// handleYoutubeRemoveStart initiates the process of removing a YouTube channel.
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

// handleYoutubeRemoveNumber removes the selected YouTube channel based on user input.
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

// handleYoutubeCommand processes all YouTube-related command actions.
func (b *Bot) handleYoutubeCommand(chatID int64, cmd commands.Command) bool {
	if cmd.Action == "" {
		b.sendYoutubeMenu(chatID)
		return true
	}

	switch cmd.Action {
	case "check":
		go youtube.CheckAllChannels()
		b.sendMessage(chatID, reply.YoutubeFormat("YouTube check started."))
		return true

	case "add":
		b.waitingForYoutubeLink[chatID] = true
		b.sendMessage(chatID, reply.YoutubeFormat("Send me a YouTube channel link or handle."))
		return true

	case "list":
		b.handleYoutubeList(chatID)
		return true

	case "remove":
		b.handleYoutubeRemoveStart(chatID)
		return true
	}

	return false
}
func (b *Bot) sendYoutubeMenu(chatID int64) {
	msg := tgbotapi.NewMessage(
		chatID,
		reply.YoutubeFormat("Choose a YouTube action."),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = commands.BuildYoutubeKeyboard()

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("failed to send youtube menu: %v", err)
	}
}
