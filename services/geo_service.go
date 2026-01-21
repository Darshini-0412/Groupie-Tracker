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

func GetArtistCoordinates(artist models.Artist) ([]Coordinates, error) {
	var coords []Coordinates

	rel, err := FetchRelationByID(artist.ID)
	if err != nil {
		return nil, err
	}

	past, future := SplitPastFutureConcerts(*rel)
	allLocations := append(past, future...)

	// Utiliser un map pour éviter les doublons
	locationSet := make(map[string]bool)

	for _, locStr := range allLocations {
		if locationSet[locStr] {
			continue
		}
		locationSet[locStr] = true

		// Nettoyer l'adresse
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

// GeocodeAddress géocode une adresse en coordonnées
func GeocodeAddress(address string) (Coordinates, error) {
	// Nettoyer l'adresse avant géocodage
	clean := cleanAddressForGeocoding(address)

	// Appel API avec adresse nettoyée
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

// cleanAddressForGeocoding nettoie une adresse pour le géocodage
func cleanAddressForGeocoding(address string) string {
	// Remplacer tirets et underscores par espaces
	clean := strings.ReplaceAll(address, "-", " ")
	clean = strings.ReplaceAll(clean, "_", " ")

	// Supprimer espaces multiples
	clean = strings.TrimSpace(clean)
	clean = strings.Join(strings.Fields(clean), " ")

	// Capitaliser correctement pour une meilleure reconnaissance
	parts := strings.Split(clean, " ")
	for i, part := range parts {
		if len(part) > 0 {
			// Garder les codes pays en majuscules (ex: USA, UK)
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
