// modexusBot/internal/helper/helper.go
package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateID creates a random identifier with the provided prefix.
func GenerateID(prefix string) string {
	rand.Seed(time.Now().UnixNano())

	id := rand.Intn(90000000) + 10000000
	return fmt.Sprintf("%s_%d", strings.ToLower(prefix), id)
}
