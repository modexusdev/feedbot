// modexusBot internal/commands/commands.go
package commands

import "github.com/modexusdev/feedbot/internal/reply"

type Command struct {
	Name      string
	Action    string
	Args      []string
	IsCommand bool
}

func Handle(cmd Command, services EnabledServices) string {
	switch cmd.Name {
	case "hello", "hi", "hey":
		return "Hello Mo 👋"
	}

	if !cmd.IsCommand {
		return ""
	}

	switch cmd.Name {

	case "help":
		return reply.Format("📚", BuildHelpText(services))

	case "youtube":
		if !services.Youtube {
			return reply.Format("❌", "YouTube service is disabled")
		}

		switch cmd.Action {
		case "add":
			return reply.Format("🎥", "Send me a YouTube channel link or handle.")
		case "list":
			return reply.Format("🎥", "List all YouTube channels.")
		case "remove":
			return reply.Format("🎥", "Remove a YouTube channel.")

		default:
			return reply.Format(
				"🎥",
				"To add a YouTube channel:\n\nyoutube add",
			)
		}

	default:
		return reply.Format("❌", "Unknown command")
	}
}
