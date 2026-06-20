// modexusBot/internal/bot/bot.go
package bot

import (
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/commands"
	"github.com/modexusdev/feedbot/internal/config"
	"github.com/modexusdev/feedbot/internal/scheduler"
	"github.com/modexusdev/feedbot/internal/storage"
)

type Bot struct {
	api                     *tgbotapi.BotAPI
	config                  config.BotConfig
	services                commands.EnabledServices
	waitingForYoutubeLink   map[int64]bool
	pendingYoutubeChannel   map[int64]storage.YoutubeChannel
	waitingForYoutubeRemove map[int64]bool
}

func New(cfg config.BotConfig, services commands.EnabledServices) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:                     api,
		config:                  cfg,
		services:                services,
		waitingForYoutubeLink:   make(map[int64]bool),
		pendingYoutubeChannel:   make(map[int64]storage.YoutubeChannel),
		waitingForYoutubeRemove: make(map[int64]bool),
	}, nil
}

func (b *Bot) Run() {
	log.Printf("Bot started: @%s", b.api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.api.GetUpdatesChan(updateConfig)

	for update := range updates {
		b.handleUpdate(update)
	}
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("failed to send message: %v", err)
	}
}

func (b *Bot) sendAutomation(text string) {
	if len(b.config.AllowedUserIDs) == 0 {
		log.Println("no allowed user id found for automation message")
		return
	}

	chatID, err := strconv.ParseInt(b.config.AllowedUserIDs[0], 10, 64)
	if err != nil {
		log.Printf("invalid automation chat id: %v", err)
		return
	}

	b.sendMessage(chatID, text)
}

func (b *Bot) ListenScheduler() {
	for msg := range scheduler.Queue {
		b.sendAutomation(msg.Text)

		time.Sleep(15 * time.Second)
	}
}
