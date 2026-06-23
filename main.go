// modexusBot/main.go
package main

import (
	"log"

	"github.com/modexusdev/feedbot/internal/bot"
	"github.com/modexusdev/feedbot/internal/commands"
	"github.com/modexusdev/feedbot/internal/config"
	"github.com/modexusdev/feedbot/internal/tracker"
)

func main() {
	// Load application configuration from environment variables.
	cfg := config.Load()

	// Define which services are enabled for this bot instance.
	services := commands.EnabledServices{
		Youtube: true,
		Weather: true,
		Github:  false,
		RSS:     false,
	}
	// Create a new bot instance with the loaded configuration and enabled services.
	app, err := bot.New(cfg, services)
	if err != nil {
		log.Fatal(err)
	}

	// Start background workers for scheduled messages and content tracking.
	go app.ListenScheduler()
	go tracker.Watch()

	// Start the Telegram bot update loop.
	app.Run()
}
