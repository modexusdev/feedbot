// modexusBot internal/commands/available.go
package commands

type AvailableCommand struct {
	Name        string
	Description string
	Emoji       string
	Service     string
}

// AvailableCommands returns all commands that are enabled
// for the configured services.
func AvailableCommands(services EnabledServices) []AvailableCommand {
	commands := []AvailableCommand{

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
		commands = append(commands, AvailableCommand{
			Name:        "youtube list",
			Description: "List all YouTube channels",
			Emoji:       "🎥",
			Service:     "youtube",
		})
		commands = append(commands, AvailableCommand{
			Name:        "youtube remove",
			Description: "Remove a YouTube channel",
			Emoji:       "🎥",
			Service:     "youtube",
		})

	}

	return commands
}
