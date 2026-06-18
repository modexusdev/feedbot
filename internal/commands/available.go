// modexusBot internal/commands/available.go
package commands

type AvailableCommand struct {
	Name        string
	Description string
	Emoji       string
	Service     string
}

func AvailableCommands(services EnabledServices) []AvailableCommand {
	commands := []AvailableCommand{
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

	if services.Youtube {
		commands = append(commands, AvailableCommand{
			Name:        "youtube add",
			Description: "Manage YouTube channels",
			Emoji:       "🎥",
			Service:     "youtube",
		})
	}

	return commands
}
