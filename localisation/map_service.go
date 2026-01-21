package localisation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Location struct {
	City    string
	Country string
	Lat     float64
	Lon     float64
}

type nominatimResponse struct {
	DisplayName string  `json:"display_name"`
	Lat         float64 `json:"lat,string"`
	Lon         float64 `json:"lon,string"`
}

var cache = make(map[string]Location)

func SearchLocation(query string) (Location, error) {

	if loc, ok := cache[query]; ok {
		return loc, nil
	}

	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("limit", "1")

	time.Sleep(1 * time.Second)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	var results []nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return Location{}, err
	}

	if len(results) == 0 {
		return Location{}, fmt.Errorf("no results found for query: %s", query)
	}

	parts := strings.Split(results[0].DisplayName, ",")
	city := strings.TrimSpace(parts[0])
	country := strings.TrimSpace(parts[len(parts)-1])

	loc := Location{
		City:    city,
		Country: country,
		Lat:     results[0].Lat,
		Lon:     results[0].Lon,
	}

	cache[query] = loc

	return loc, nil
}
