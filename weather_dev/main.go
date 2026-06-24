package main

import (
	"encoding/json"
	"fmt"

	"github.com/modexusdev/feedbot/internal/weather"
)

func main() {
	lat := 51.4828
	lon := 11.9697

	result, err := weather.FetchWeather(lat, lon, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
}
