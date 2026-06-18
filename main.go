// modexusBot/main.go
package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modexusdev/feedbot/internal/commands"
	"github.com/modexusdev/feedbot/internal/config"
)

func main() {
	cfg := config.Load()

	services := commands.EnabledServices{
		Youtube: true,
		Github:  false,
		RSS:     false,
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot started: @%s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := fmt.Sprintf("%d", update.Message.From.ID)

		if !config.IsAllowedUser(userID, cfg.AllowedUserIDs) {
			log.Printf(
				"Unauthorized access: %d (%s)",
				update.Message.From.ID,
				update.Message.From.UserName,
			)
			continue
		}

		log.Printf(
			"User: %s | Message: %s",
			update.Message.From.UserName,
			update.Message.Text,
		)

		cmd := commands.Parse(update.Message.Text)
		response := commands.Handle(cmd, services)

		if response == "" {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ParseMode = "HTML"

		if cmd.Name == "help" {
			log.Println("Keyboard attached")
			msg.ReplyMarkup = commands.BuildKeyboard(services)
		}
		log.Printf("Command: %s", cmd.Name)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	}
}
