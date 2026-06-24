// modexusBot/weather_dev/main.go
package main

import (
	"fmt"

	"github.com/modexusdev/feedbot/internal/weather"
)

func main() {
	// lat := 51.4828
	// lon := 11.9697

	// result, err := weather.GetWeather(lat, lon, "Halle Saale ", 0)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// data, err := json.MarshalIndent(result, "", "  ")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(string(data))
	city := "Sankt"

	fmt.Println("Searching for:", city)
	fmt.Println()

	locations, err := weather.GeocodeCity(city)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Results for %q:\n\n", city)

	for i, loc := range locations {
		fmt.Printf(
			"%d) %s, %s, %s | %.5f, %.5f\n",
			i+1,
			loc.Name,
			loc.Admin1,
			loc.Country,
			loc.Latitude,
			loc.Longitude,
		)
	}
}
