package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/localisation"
	"groupie-tracker/models"
	"net/http"
	"strings"
)

// Coordinates représente les coordonnées GPS d'un lieu
type Coordinates struct {
	City    string
	Country string
	Lat     float64 // Latitude
	Lon     float64 // Longitude
}

// Fetch est une fonction générique pour récupérer des données JSON
// Le [T any] c'est la "magie" des generics de Go 1.18+
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

// GetArtistCoordinates récupère les coordonnées GPS de tous les concerts d'un artiste
func GetArtistCoordinates(artist models.Artist) ([]Coordinates, error) {
	var coords []Coordinates

	rel, err := FetchRelationByID(artist.ID)
	if err != nil {
		return nil, err
	}

	// Récupération des lieux passés ET futurs
	past, future := SplitPastFutureConcerts(*rel)
	allLocations := append(past, future...)

	// Map pour éviter de géocoder 2 fois le même lieu
	locationSet := make(map[string]bool)

	for _, locStr := range allLocations {
		if locationSet[locStr] {
			continue // Déjà traité
		}
		locationSet[locStr] = true

		// Nettoyage de l'adresse avant géocodage
		clean := cleanAddressForGeocoding(locStr)

		result, err := localisation.SearchLocation(clean)
		if err != nil {
			fmt.Printf("Erreur géocodage pour %s: %v\n", clean, err)
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

// GeocodeAddress convertit une adresse en coordonnées GPS
func GeocodeAddress(address string) (Coordinates, error) {
	// Nettoyage avant l'appel API
	clean := cleanAddressForGeocoding(address)

	result, err := localisation.SearchLocation(clean)
	if err != nil {
		return Coordinates{}, fmt.Errorf("erreur géocodage pour '%s': %v", clean, err)
	}

	return Coordinates{
		City:    result.City,
		Country: result.Country,
		Lat:     result.Lat,
		Lon:     result.Lon,
	}, nil
}

// cleanAddressForGeocoding prépare une adresse pour le géocodage
// Ex: "los_angeles-usa" → "Los Angeles USA"
func cleanAddressForGeocoding(address string) string {
	// Remplacer tirets et underscores par espaces
	clean := strings.ReplaceAll(address, "-", " ")
	clean = strings.ReplaceAll(clean, "_", " ")

	// Supprimer espaces multiples
	clean = strings.TrimSpace(clean)
	clean = strings.Join(strings.Fields(clean), " ")

	// Capitalisation intelligente
	parts := strings.Split(clean, " ")
	for i, part := range parts {
		if len(part) > 0 {
			// Garder codes pays en majuscules (USA, UK, etc.)
			if len(part) <= 3 && strings.ToUpper(part) == part {
				parts[i] = strings.ToUpper(part)
			} else {
				// Capitaliser la première lettre
				parts[i] = strings.Title(strings.ToLower(part))
			}
		}
	}

	return strings.Join(parts, " ")
}
