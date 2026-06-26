// modexusBot/internal/weather/helper.go
package weather

import (
	"fmt"
	"strings"
)

func (loc GeoLocation) Flag() string {
	switch strings.ToLower(strings.TrimSpace(loc.Country)) {

	case "deutschland", "germany":
		return "🇩🇪"

	case "vereinigte staaten von amerika", "united states", "united states of america", "usa":
		return "🇺🇸"

	case "kanada", "canada":
		return "🇨🇦"

	case "vereinigtes königreich", "united kingdom", "england", "großbritannien", "great britain":
		return "🇬🇧"

	case "irland", "ireland":
		return "🇮🇪"

	case "frankreich", "france":
		return "🇫🇷"

	case "spanien", "spain":
		return "🇪🇸"

	case "italien", "italy":
		return "🇮🇹"

	case "portugal":
		return "🇵🇹"

	case "niederlande", "netherlands":
		return "🇳🇱"

	case "belgien", "belgium":
		return "🇧🇪"

	case "luxemburg", "luxembourg":
		return "🇱🇺"

	case "schweiz", "switzerland":
		return "🇨🇭"

	case "österreich", "austria":
		return "🇦🇹"

	case "polen", "poland":
		return "🇵🇱"

	case "tschechien", "czech republic", "czechia":
		return "🇨🇿"

	case "slowakei", "slovakia":
		return "🇸🇰"

	case "ungarn", "hungary":
		return "🇭🇺"

	case "slowenien", "slovenia":
		return "🇸🇮"

	case "kroatien", "croatia":
		return "🇭🇷"

	case "rumänien", "romania":
		return "🇷🇴"

	case "bulgarien", "bulgaria":
		return "🇧🇬"

	case "griechenland", "greece":
		return "🇬🇷"

	case "dänemark", "denmark":
		return "🇩🇰"

	case "schweden", "sweden":
		return "🇸🇪"

	case "norwegen", "norway":
		return "🇳🇴"

	case "finnland", "finland":
		return "🇫🇮"

	case "island", "iceland":
		return "🇮🇸"

	case "estland", "estonia":
		return "🇪🇪"

	case "lettland", "latvia":
		return "🇱🇻"

	case "litauen", "lithuania":
		return "🇱🇹"

	case "ukraine":
		return "🇺🇦"

	case "russland", "russia":
		return "🇷🇺"

	case "türkei", "turkey":
		return "🇹🇷"

	case "china":
		return "🇨🇳"

	case "japan":
		return "🇯🇵"

	case "südkorea", "south korea":
		return "🇰🇷"

	case "indien", "india":
		return "🇮🇳"

	case "australien", "australia":
		return "🇦🇺"

	case "neuseeland", "new zealand":
		return "🇳🇿"

	case "mexiko", "mexico":
		return "🇲🇽"

	case "brasilien", "brazil":
		return "🇧🇷"

	case "argentinien", "argentina":
		return "🇦🇷"

	case "südafrika", "south africa":
		return "🇿🇦"

	default:
		return "🌍"
	}
}
func (loc GeoLocation) Subtitle() string {
	switch {
	case loc.Admin1 != "" && loc.Country != "":
		return fmt.Sprintf("%s • %s", loc.Admin1, loc.Country)

	case loc.Admin1 != "":
		return loc.Admin1

	case loc.Country != "":
		return loc.Country

	default:
		return ""
	}
}
