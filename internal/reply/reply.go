// modexusBot/internal/reply/reply.go
package reply

import (
	"fmt"
	"html"

	"github.com/modexusdev/feedbot/internal/storage"
)

func Format(emoji, text string) string {
	return "<b>🚀 FeedBot:</b>\n━━━━━━━━━━━━\n\n" + emoji + " " + text
}
func YoutubeFormat(text string) string {
	return "<b>🚀 FeedBot:</b>\n━━━━━━━━━━━━\n\n<b>🎥 YouTube</b>\n\n" + text
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
