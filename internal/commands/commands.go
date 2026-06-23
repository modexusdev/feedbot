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

	case "help":
		return reply.Format("📚", BuildHelpText(services))
	}

	return ""
}
