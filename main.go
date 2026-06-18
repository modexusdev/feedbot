package main

import (
	"log"

	"github.com/modexusdev/feedbot/internal/bot"
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

	app, err := bot.New(cfg, services)
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
