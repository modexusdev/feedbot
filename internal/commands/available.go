// modexusBot internal/commands/available.go
package commands

import "github.com/modexusdev/feedbot/internal/i18n"

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
			Description: i18n.T("command.help.description"),
			Emoji:       "📚",
		},
	}

	if services.Youtube {
		commands = append(commands, AvailableCommand{
			Name:        "youtube check",
			Description: i18n.T("command.youtube_check.description"),
			Emoji:       "🎥",
			Service:     "youtube",
		})

		commands = append(commands, AvailableCommand{
			Name:        "youtube add",
			Description: i18n.T("command.youtube_add.description"),
			Emoji:       "🎥",
			Service:     "youtube",
		})

		commands = append(commands, AvailableCommand{
			Name:        "youtube list",
			Description: i18n.T("command.youtube_list.description"),
			Emoji:       "🎥",
			Service:     "youtube",
		})

		commands = append(commands, AvailableCommand{
			Name:        "youtube remove",
			Description: i18n.T("command.youtube_remove.description"),
			Emoji:       "🎥",
			Service:     "youtube",
		})

	}

	if services.Weather {
		commands = append(commands, AvailableCommand{
			Name:        "weather today",
			Description: i18n.T("command.weather_today.description"),
			Emoji:       "🌤",
			Service:     "weather",
		})

		commands = append(commands, AvailableCommand{
			Name:        "weather tomorrow",
			Description: i18n.T("command.weather_tomorrow.description"),
			Emoji:       "🌤",
			Service:     "weather",
		})

		commands = append(commands, AvailableCommand{
			Name:        "set location",
			Description: i18n.T("command.set_location.description"),
			Emoji:       "📍",
			Service:     "weather",
		})
	}

	return commands
}
