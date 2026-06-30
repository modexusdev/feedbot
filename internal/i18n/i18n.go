// modexusBot/internal/i18n/i18n.go
package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const DefaultLanguage = "de"

var (
	currentLanguage = DefaultLanguage
	translations    = map[string]map[string]string{}
	mu              sync.RWMutex
)

func Init() error {
	basePath := "internal/i18n/locales"

	for _, code := range []string{"de", "en"} {
		filePath := filepath.Join(basePath, code+".json")

		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read language file %s: %w", filePath, err)
		}

		var messages map[string]string
		if err := json.Unmarshal(data, &messages); err != nil {
			return fmt.Errorf("parse language file %s: %w", filePath, err)
		}

		translations[code] = messages
	}

	return nil
}

func SetLanguage(code string) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := translations[code]; ok {
		currentLanguage = code
		return
	}

	currentLanguage = DefaultLanguage
}

func GetLanguage() string {
	mu.RLock()
	defer mu.RUnlock()

	return currentLanguage
}

func Lang() func(string) string {
	return func(key string) string {
		return T(key)
	}
}

func T(key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if msg, ok := translations[currentLanguage][key]; ok {
		return msg
	}

	if msg, ok := translations[DefaultLanguage][key]; ok {
		return msg
	}

	return key
}
