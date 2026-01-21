package localisation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Location = lieu avec coordonnées GPS
type Location struct {
	City    string
	Country string
	Lat     float64
	Lon     float64
}

// Réponse API OpenStreetMap
type nominatimResponse struct {
	DisplayName string  `json:"display_name"`
	Lat         float64 `json:"lat,string"`
	Lon         float64 `json:"lon,string"`
}

// Cache des résultats
var cache = make(map[string]Location)

// SearchLocation transforme une adresse en coordonnées GPS
func SearchLocation(query string) (Location, error) {

	// Vérifier cache
	if loc, ok := cache[query]; ok {
		return loc, nil
	}

	// Préparer URL API
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("limit", "1")

	time.Sleep(1 * time.Second) // Règle API

	// Appeler API
	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	// Décoder JSON
	var results []nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return Location{}, err
	}

	if len(results) == 0 {
		return Location{}, fmt.Errorf("no results found for query: %s", query)
	}

	// Extraire ville et pays
	parts := strings.Split(results[0].DisplayName, ",")
	city := strings.TrimSpace(parts[0])
	country := strings.TrimSpace(parts[len(parts)-1])

	// Créer Location
	loc := Location{
		City:    city,
		Country: country,
		Lat:     results[0].Lat,
		Lon:     results[0].Lon,
	}

	cache[query] = loc
	return loc, nil
}
