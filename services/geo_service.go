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

// --- Fonction Fetch générique (corrige ton erreur "undefined: Fetch") ---
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

// --- Récupère les coordonnées pour chaque lieu d’un artiste ---
func GetArtistCoordinates(artist models.Artist) ([]Coordinates, error) {
	var coords []Coordinates

	// Récupérer les lieux de l'artiste depuis l'API
	location, err := Fetch[models.Location](artist.Locations)
	if err != nil {
		return nil, err
	}

	for _, locStr := range location.Locations {
		// Format attendu : "Ville, Pays"
		parts := strings.Split(locStr, ",")
		if len(parts) != 2 {
			continue
		}

		city := strings.TrimSpace(parts[0])
		country := strings.TrimSpace(parts[1])

		// Exemple : "Paris France"
		query := fmt.Sprintf("%s %s", city, country)

		result, err := localisation.SearchLocation(query)
		if err != nil {
			return nil, err
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
