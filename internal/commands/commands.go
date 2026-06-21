// modexusBot internal/commands/commands.go
package commands

import "github.com/modexusdev/feedbot/internal/reply"

// Command represents a parsed user command.
type Command struct {
	Name      string
	Action    string
	Args      []string
	IsCommand bool
}

// Handle returns the response text for a parsed command.
func Handle(cmd Command, services EnabledServices) string {
	switch cmd.Name {
	case "hello", "hi", "hey":
		return "Hello 👋"
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
		case "check":
			return reply.Format("🎥", "YouTube check started.")
		case "add":
			return reply.Format("🎥", "Send me a YouTube channel link or handle.")
		case "list":
			return reply.Format("🎥", "List all YouTube channels.")
		case "remove":
			return reply.Format("🎥", "Remove a YouTube channel.")

		default:
			return reply.Format("🎥", "Choose a YouTube action.")
		}

	default:
		return reply.Format("❌", "Unknown command")
	}
}
