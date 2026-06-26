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
	"github.com/modexusdev/feedbot/internal/weather"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	config   config.BotConfig
	services commands.EnabledServices
	// Youtube Service tracking
	waitingForYoutubeLink   map[int64]bool
	pendingYoutubeChannel   map[int64]storage.YoutubeChannel
	waitingForYoutubeRemove map[int64]bool
	// Weather services tracking
	waitingForWeatherLocation       map[int64]bool
	pendingWeatherLocations         map[int64][]weather.GeoLocation
	waitingForWeatherLocationNumber map[int64]bool
}

// New creates a new Telegram bot instance with the provided configuration.
func New(cfg config.BotConfig, services commands.EnabledServices) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:      api,
		config:   cfg,
		services: services,
		// Youtube Service tracking
		waitingForYoutubeLink:   make(map[int64]bool),
		pendingYoutubeChannel:   make(map[int64]storage.YoutubeChannel),
		waitingForYoutubeRemove: make(map[int64]bool),
		// Weather services tracking
		waitingForWeatherLocation:       make(map[int64]bool),
		pendingWeatherLocations:         make(map[int64][]weather.GeoLocation),
		waitingForWeatherLocationNumber: make(map[int64]bool),
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
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("failed to send message: %v", err)
	}
}

// Sends an automation message to the allowed user ID, if configured.
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

// Listens to the scheduler queue and sends automation messages to the allowed user ID, if configured.
func (b *Bot) ListenScheduler() {
	for msg := range scheduler.Queue {
		b.sendAutomation(msg.Text)
		// Prevent Telegram message bursts when multiple events arrive at once.
		time.Sleep(5 * time.Second)
	}
}
