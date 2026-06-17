// modexusBot/internal/commands/available.go
package commands

type AvailableCommand struct {
	Name        string
	Description string
	Emoji       string
}

var AvailableCommands = []AvailableCommand{
	{
		Name:        "ping",
		Description: "Check if bot is online",
		Emoji:       "🏓",
	},
	{
		Name:        "help",
		Description: "Show available commands",
		Emoji:       "📚",
	},
}
