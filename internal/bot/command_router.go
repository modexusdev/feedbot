// modexusBot/internal/bot/command_router.go
package bot

import (
	"github.com/modexusdev/feedbot/internal/commands"
	"github.com/modexusdev/feedbot/internal/youtube"
)

// handleCommandAction routes commands to the appropriate service handler.
func (b *Bot) handleCommandAction(chatID int64, cmd commands.Command) bool {
	switch cmd.Name {
	case "youtube":
		return b.handleYoutubeCommand(chatID, cmd)

	default:
		return false
	}
}

// handleYoutubeCommand processes all YouTube-related command actions.
func (b *Bot) handleYoutubeCommand(chatID int64, cmd commands.Command) bool {
	switch cmd.Action {
	case "check":
		go youtube.CheckAllChannels()
		return false
	case "add":
		b.waitingForYoutubeLink[chatID] = true
		return false

	case "list":
		b.handleYoutubeList(chatID)
		return true

	case "remove":
		b.handleYoutubeRemoveStart(chatID)
		return true
	}

	return false
}
