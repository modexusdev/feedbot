// modexusBot internal/commands/service.go
package commands

// EnabledServices defines which FeedBot modules are available.
type EnabledServices struct {
	Youtube bool
	Weather bool
	Github  bool
	RSS     bool
}
