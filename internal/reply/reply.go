// modexusBot/internal/reply/reply.go
package reply

import (
	"fmt"
	"html"
	"strings"

	"github.com/modexusdev/feedbot/internal/storage"
)

func Format(emoji, text string) string {
	return "<b>🚀 FeedBot:</b>\n━━━━━━━━━━━━\n\n" + emoji + " " + text
}
func YoutubeFormat(text string) string {
	return "<b>🚀 FeedBot:</b>\n━━━━━━━━━━━━\n\n<b>🎥 YouTube:</b>\n\n" + text
}

func YoutubeAddFormat(channel storage.YoutubeChannel) string {
	return YoutubeFormat(
		fmt.Sprintf(
			"Do you want to add this YouTube channel?\n\n<b>Name:</b> %s\n<b>Handle:</b> %s",
			html.EscapeString(channel.Name),
			html.EscapeString(channel.Handle),
		),
	)
}
func YoutubeAlreadyAddedFormat(channel storage.YoutubeChannel) string {
	return YoutubeFormat(
		fmt.Sprintf(
			"⚠️ This YouTube channel has already been added.\n\n<b>Name:</b> %s\n<b>Handle:</b> %s",
			html.EscapeString(channel.Name),
			html.EscapeString(channel.Handle),
		),
	)
}
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

func AutomationFormat(sourceEmoji, sourceName, text string) string {
	return "<b>🚀 FeedBot • Automation 🤖</b>\n━━━━━━━━━━━━\n\n<b>" +
		sourceEmoji + " " + sourceName + ":</b>\n\n" +
		text
}
