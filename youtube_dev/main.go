package main

import (
	"fmt"
	"strings"

	"github.com/modexusdev/feedbot/internal/youtube"
)

func main() {
	for {
		var link string

		fmt.Print("Gib einen Youtube Link ein (exit zum Beenden): ")
		fmt.Scanln(&link)

		if strings.EqualFold(link, "exit") {
			fmt.Println("Programm beendet.")
			break
		}

		err := youtube.ExtractYoutubeLink(link)
		if err != nil {
			fmt.Println("❌ Fehler:", err)
			continue
		}

		fmt.Println("✅ Kanal gespeichert")
	}
}
