// modexusBot/internal/i18n/i18n.go
package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Language struct {
	Code string
	Name string
	Flag string
}

const DefaultLanguage = "en"

var AvailableLanguages = []Language{
	{Code: "de", Name: "Deutsch", Flag: "🇩🇪"},
	{Code: "en", Name: "English", Flag: "🇬🇧"},
}

var (
	currentLanguage = DefaultLanguage
	translations    = map[string]map[string]string{}
	mu              sync.RWMutex
)

func Init() error {
	basePath := "internal/i18n/locales"

	for _, lang := range AvailableLanguages {
		filePath := filepath.Join(basePath, lang.Code+".json")

		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read language file %s: %w", filePath, err)
		}

		var messages map[string]string
		if err := json.Unmarshal(data, &messages); err != nil {
			return fmt.Errorf("parse language file %s: %w", filePath, err)
		}

		translations[lang.Code] = messages
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

func IsSupported(code string) bool {
	mu.RLock()
	defer mu.RUnlock()

	_, ok := translations[code]
	return ok
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
func GetAvailableLanguages() []Language {
	return AvailableLanguages
}
func GetLanguageName(code string) string {
	for _, lang := range AvailableLanguages {
		if lang.Code == code {
			return lang.Name
		}
	}

	return code
}
