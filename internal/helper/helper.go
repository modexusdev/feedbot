// modexusBot/internal/helper/helper.go
package helper

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GenerateID(prefix string) string {
	rand.Seed(time.Now().UnixNano())

	id := rand.Intn(90000000) + 10000000
	return fmt.Sprintf("%s_%d", strings.ToLower(prefix), id)
}

func DownloadFile(fileURL string, savePath string) error {
	if fileURL == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	return err
}
