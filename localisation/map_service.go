// localisation/map_service.go
package localisation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

var (
	cache        = make(map[string]Location)
	lastRequest  time.Time
	requestMutex sync.Mutex
)

// SearchLocation transforme une adresse en coordonnées GPS
func SearchLocation(query string) (Location, error) {
	// vérifier si on l'a déjà cherché
	if loc, ok := cache[query]; ok {
		return loc, nil
	}

	// respecter le rate limit de 1 req/sec
	requestMutex.Lock()
	timeSinceLastRequest := time.Since(lastRequest)
	if timeSinceLastRequest < time.Second {
		time.Sleep(time.Second - timeSinceLastRequest)
	}
	lastRequest = time.Now()
	requestMutex.Unlock()

	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("limit", "1")

	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return Location{}, err
	}

	// headers requis par Nominatim
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	// si rate limited (429), attendre et réessayer
	if resp.StatusCode == 429 {
		time.Sleep(2 * time.Second)
		return SearchLocation(query) // retry
	}

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

	// sauvegarder dans le cache
	cache[query] = loc
	return loc, nil
}
