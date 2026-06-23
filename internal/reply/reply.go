// modexusBot/internal/reply/reply.go
package reply

import (
	"fmt"
	"html"
	"strings"

	"github.com/modexusdev/feedbot/internal/storage"
)

// Format creates a standard FeedBot message.
func Format(emoji, text string) string {
	return "<b>🚀 FeedBot:</b>\n━━━━━━━━━━━━\n\n" + emoji + " " + text
}

// YoutubeFormat creates a formatted YouTube message.
func YoutubeFormat(text string) string {
	return Format("🎥", "<b>YouTube:</b>\n\n"+text)

}

// YoutubeAddFormat creates a confirmation message for adding a channel.
func YoutubeAddFormat(channel storage.YoutubeChannel) string {
	return YoutubeFormat(
		fmt.Sprintf(
			"Do you want to add this YouTube channel?\n\n<b>Name:</b> %s\n<b>Handle:</b> %s",
			html.EscapeString(channel.Name),
			html.EscapeString(channel.Handle),
		),
	)
}

// YoutubeAlreadyAddedFormat creates a message for channels that already exist.
func YoutubeAlreadyAddedFormat(channel storage.YoutubeChannel) string {
	return YoutubeFormat(
		fmt.Sprintf(
			"⚠️ This YouTube channel has already been added.\n\n<b>Name:</b> %s\n<b>Handle:</b> %s",
			html.EscapeString(channel.Name),
			html.EscapeString(channel.Handle),
		),
	)
}

// YoutubeListFormat creates a formatted list of saved channels.
func YoutubeListFormat(channels []storage.YoutubeChannel) string {
	if len(channels) == 0 {
		return YoutubeFormat("No YouTube channels saved.")
	}

	var text strings.Builder

	text.WriteString("<b>YouTube Channel List:</b>\n\n")

	for i, channel := range channels {
		text.WriteString(fmt.Sprintf(
			"%d. %s\n",
			i+1,
			html.EscapeString(channel.Name),
		))
	}

	return YoutubeFormat(text.String())
}

func WeatherFormat(text string) string {
	return Format("🌦️", "<b>Weather:</b>\n\n"+text)
}

// AutomationFormat creates a formatted automation notification message.
func AutomationFormat(sourceEmoji, sourceName, text string) string {
	return Format("🚀", "<b>Automation:</b>\n\n"+text)
}
