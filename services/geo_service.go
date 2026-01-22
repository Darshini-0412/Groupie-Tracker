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

// Fetch fonction générique pour récupérer du JSON
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

// GetArtistCoordinates récupère les coordonnées GPS de tous les concerts
func GetArtistCoordinates(artist models.Artist) ([]Coordinates, error) {
	var coords []Coordinates

	rel, err := FetchRelationByID(artist.ID)
	if err != nil {
		return nil, err
	}

	past, future := SplitPastFutureConcerts(*rel)
	allLocations := append(past, future...)

	locationSet := make(map[string]bool)

	for _, locStr := range allLocations {
		if locationSet[locStr] {
			continue
		}
		locationSet[locStr] = true

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

func GeocodeAddress(address string) (Coordinates, error) {
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
func cleanAddressForGeocoding(address string) string {
	clean := strings.ReplaceAll(address, "-", " ")
	clean = strings.ReplaceAll(clean, "_", " ")

	clean = strings.TrimSpace(clean)
	clean = strings.Join(strings.Fields(clean), " ")

	parts := strings.Split(clean, " ")
	for i, part := range parts {
		if len(part) > 0 {
			if len(part) <= 3 && strings.ToUpper(part) == part {
				parts[i] = strings.ToUpper(part)
			} else {
				parts[i] = strings.Title(strings.ToLower(part))
			}
		}
	}

	return strings.Join(parts, " ")
}
