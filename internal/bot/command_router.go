package bot

import (
	"github.com/modexusdev/feedbot/internal/commands"
)

// handleCommandAction routes commands to the appropriate service handler.
func (b *Bot) handleCommandAction(chatID int64, cmd commands.Command) bool {
	switch cmd.Name {
	case "youtube":
		return b.handleYoutubeCommand(chatID, cmd)

	case "weather":
		return b.handleWeatherCommand(chatID, cmd)

	default:
		return false
	}
}
