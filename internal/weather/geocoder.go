// modexusBot/internal/weather/geocoder.go
package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GeoLocation struct {
	Name      string
	Country   string
	Admin1    string
	Latitude  float64
	Longitude float64
}

type nominatimResult struct {
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Address     struct {
		City    string `json:"city"`
		Town    string `json:"town"`
		Village string `json:"village"`
		State   string `json:"state"`
		Country string `json:"country"`
	} `json:"address"`
}

// GeocodeCity returns the geolocation for the given city name
func GeocodeCity(city string) ([]GeoLocation, error) {
	query := normalizeCityQuery(city)

	apiURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?q=%s&format=json&addressdetails=1&limit=10&accept-language=de",
		url.QueryEscape(query),
	)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "modexusBot/1.0")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoder api returned status: %d", resp.StatusCode)
	}

	var raw []nominatimResult

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return nil, fmt.Errorf("city not found")
	}

	locations := make([]GeoLocation, 0, len(raw))

	for _, r := range raw {
		name := r.Address.City
		if name == "" {
			name = r.Address.Town
		}
		if name == "" {
			name = r.Address.Village
		}
		if name == "" {
			name = r.DisplayName
		}

		lat, lon, err := parseLatLon(r.Lat, r.Lon)
		if err != nil {
			continue
		}

		locations = append(locations, GeoLocation{
			Name:      name,
			Country:   r.Address.Country,
			Admin1:    r.Address.State,
			Latitude:  lat,
			Longitude: lon,
		})
	}

	if len(locations) == 0 {
		return nil, fmt.Errorf("city not found")
	}

	return locations, nil
}

// normalizeCityQuery returns the normalized city query string
func normalizeCityQuery(city string) string {
	city = strings.TrimSpace(city)
	city = strings.ReplaceAll(city, "_", " ")
	city = strings.ReplaceAll(city, "-", " ")

	parts := strings.Fields(city)
	return strings.Join(parts, " ")
}

// parseLatLon parses the latitude and longitude values from the given strings
func parseLatLon(latValue, lonValue string) (float64, float64, error) {
	lat, err := strconv.ParseFloat(latValue, 64)
	if err != nil {
		return 0, 0, err
	}

	lon, err := strconv.ParseFloat(lonValue, 64)
	if err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}

// FormatLocationList returns the formatted location list for the given locations
func FormatLocationList(locations []GeoLocation) string {
	var b strings.Builder

	b.WriteString(" <b>Standort auswählen</b>\n")
	b.WriteString("Sende die passende Nummer.\n\n")

	for i, loc := range locations {
		fmt.Fprintf(
			&b,
			"<b>%d.</b> %s %s\n",
			i+1,
			loc.Flag(),
			loc.Title(),
		)

		if subtitle := loc.Subtitle(); subtitle != "" {
			b.WriteString("   " + subtitle + "\n")
		}

		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}

// FormatSelectedLocation returns the formatted selected location for the given location
func FormatSelectedLocation(loc GeoLocation) string {
	var b strings.Builder

	b.WriteString("🌍 <b>Standort gespeichert</b>\n\n")
	fmt.Fprintf(&b, "%s %s\n", loc.Flag(), loc.Title())

	if subtitle := loc.Subtitle(); subtitle != "" {
		b.WriteString(subtitle)
	}

	return b.String()
}

func (loc GeoLocation) Title() string {
	if loc.Name != "" {
		return loc.Name
	}

	return "Unbekannter Standort"
}
