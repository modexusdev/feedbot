package commands

import (
	"github.com/modexusdev/feedbot/internal/reply"
)

type Command struct {
	Name      string
	Action    string
	Args      []string
	IsCommand bool
}

func Handle(cmd Command) string {

	switch cmd.Name {

	case "hello", "hi", "hey":
		return "Hello Mo 👋"
	}

	if !cmd.IsCommand {
		return ""
	}

	switch cmd.Name {

	case "ping":
		return reply.Format("🏓", "Pong")

	case "help":
		return reply.Format("📚", BuildHelpText())

	default:
		return "Unknown command"
	}
}
