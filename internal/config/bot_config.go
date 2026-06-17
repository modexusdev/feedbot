package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type BotConfig struct {
	Token          string
	AllowedUserIDs []string
}

func Load() BotConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	allowedUserIDs := os.Getenv("ALLOWED_USER_IDS")
	if allowedUserIDs == "" {
		log.Fatal("ALLOWED_USER_IDS is not set")
	}

	return BotConfig{
		Token:          token,
		AllowedUserIDs: strings.Split(allowedUserIDs, ","),
	}
}

func IsAllowedUser(userID string, allowedIDs []string) bool {
	for _, id := range allowedIDs {
		if userID == id {
			return true
		}
	}

	return false
}
