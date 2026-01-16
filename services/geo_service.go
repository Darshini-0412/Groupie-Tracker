package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/localisation"
	"groupie-tracker/models"
	"net/http"
	"strings"
)

type Coordinates struct {
	City    string
	Country string
	Lat     float64
	Lon     float64
}

// --- Fonction Fetch générique ---
func Fetch[T any](url string) (T, error) {
	var result T

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// --- Récupère les coordonnées pour un artiste ---
func GetArtistCoordinates(artist models.Artist) ([]Coordinates, error) {
	var coords []Coordinates

	// 1. Récupérer les relations de l'artiste
	rel, err := FetchRelationByID(artist.ID)
	if err != nil {
		return nil, err
	}

	// 2. Séparer concerts passés / futurs
	past, future := SplitPastFutureConcerts(*rel)

	// 3. Fusionner les deux listes (ou choisir l’une des deux)
	allLocations := append(past, future...)

	// 4. Géolocaliser chaque lieu
	for _, locStr := range allLocations {
		// Format API : "Paris-France"
		parts := strings.Split(locStr, "-")
		if len(parts) != 2 {
			continue
		}

		city := strings.TrimSpace(parts[0])
		country := strings.TrimSpace(parts[1])

		query := fmt.Sprintf("%s %s", city, country)

		result, err := localisation.SearchLocation(query)
		if err != nil {
			continue
		}

		coords = append(coords, Coordinates{
			City:    result.City,
			Country: result.Country,
			Lat:     result.Lat,
			Lon:     result.Lon,
		})
	}

	return coords, nil
}

// --- Géocode une adresse unique ---
func GeocodeAddress(address string) (Coordinates, error) {
	result, err := localisation.SearchLocation(address)
	if err != nil {
		return Coordinates{}, err
	}

	return Coordinates{
		City:    result.City,
		Country: result.Country,
		Lat:     result.Lat,
		Lon:     result.Lon,
	}, nil
}
