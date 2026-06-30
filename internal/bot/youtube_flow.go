// modexusBot/internal/bot/youtube_flow.go
package bot

import (
	"fmt"
	"html"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/commands"
	"github.com/modexusdev/feedbot/internal/i18n"
	"github.com/modexusdev/feedbot/internal/reply"
	"github.com/modexusdev/feedbot/internal/storage"
	"github.com/modexusdev/feedbot/internal/youtube"
)

func (b *Bot) handleYoutubeLink(chatID int64, text string) {
	b.waitingForYoutubeLink[chatID] = false

	channel, err := youtube.ExtractYoutubeChannel(text)
	if err != nil {
		b.sendMessage(chatID, reply.Format("❌", i18n.T("youtube.read_channel_error")))
		return
	}

	if storage.YoutubeChannelExists(channel.Handle, channel.RSSURL) {
		b.sendMessage(chatID, reply.YoutubeAlreadyAddedFormat(channel))
		return
	}

	b.pendingYoutubeChannel[chatID] = channel

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(i18n.T("button.yes"), "youtube_add_yes"),
			tgbotapi.NewInlineKeyboardButtonData(i18n.T("button.no"), "youtube_add_no"),
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

func (b *Bot) handleYoutubeAddConfirm(chatID int64) {
	channel, ok := b.pendingYoutubeChannel[chatID]
	if !ok {
		b.sendMessage(chatID, reply.Format("❌", i18n.T("youtube.no_pending_channel")))
		return
	}

	if _, err := storage.SaveYoutubeChannel(channel); err != nil {
		b.sendMessage(chatID, reply.Format("❌", i18n.T("youtube.save_channel_error")))
		return
	}

	delete(b.pendingYoutubeChannel, chatID)

	b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.channel_added")))
}

func (b *Bot) handleYoutubeAddCancel(chatID int64) {
	delete(b.pendingYoutubeChannel, chatID)

	b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.channel_not_added")))
}

func (b *Bot) handleYoutubeList(chatID int64) {
	channels, err := storage.GetYoutubeChannels()
	if err != nil {
		b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.load_channels_error")))
		return
	}

	b.sendMessage(chatID, reply.YoutubeListFormat(channels))
}

func (b *Bot) handleYoutubeRemoveStart(chatID int64) {
	channels, err := storage.GetYoutubeChannels()
	if err != nil {
		b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.load_channels_error")))
		return
	}

	if len(channels) == 0 {
		b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.no_channels_found")))
		return
	}

	b.waitingForYoutubeRemove[chatID] = true

	b.sendMessage(
		chatID,
		reply.YoutubeListFormat(channels)+"\n\n"+i18n.T("youtube.remove_write_number"),
	)
}

func (b *Bot) handleYoutubeRemoveNumber(chatID int64, text string) {
	b.waitingForYoutubeRemove[chatID] = false

	index, err := strconv.Atoi(text)
	if err != nil || index < 1 {
		b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.invalid_number")))
		return
	}

	channels, err := storage.GetYoutubeChannels()
	if err != nil {
		b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.load_channels_error")))
		return
	}

	if index > len(channels) {
		b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.channel_number_not_found")))
		return
	}

	channel := channels[index-1]

	if err := storage.DeleteYoutubeChannel(channel.ID); err != nil {
		b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.remove_channel_error")))
		return
	}

	b.sendMessage(
		chatID,
		reply.YoutubeFormat(
			fmt.Sprintf(
				i18n.T("youtube.channel_removed"),
				html.EscapeString(channel.Name),
			),
		),
	)
}

func (b *Bot) handleYoutubeCommand(chatID int64, cmd commands.Command) bool {
	if cmd.Action == "" {
		b.sendYoutubeMenu(chatID)
		return true
	}

	switch cmd.Action {
	case "check":
		go youtube.CheckAllChannels()
		b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.check_started")))
		return true

	case "add":
		b.waitingForYoutubeLink[chatID] = true
		b.sendMessage(chatID, reply.YoutubeFormat(i18n.T("youtube.send_channel_link")))
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
		reply.YoutubeFormat(i18n.T("youtube.choose_action")),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = commands.BuildYoutubeKeyboard()

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("failed to send youtube menu: %v", err)
	}
}
