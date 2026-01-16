package services

import (
	"fmt"
	"groupie-tracker/localisation"
	"groupie-tracker/models"
	"strings"
)

type Coordinates struct {
	City    string
	Country string
	Lat     float64
	Lon     float64
}

// Récupère les coordonnées pour chaque lieu d’un artiste
func GetArtistCoordinates(artist models.Artist) ([]Coordinates, error) {
	var coords []Coordinates

	// Récupérer les lieux de l'artiste depuis l'API
	location, err := Fetch[models.Location](artist.Locations)
	if err != nil {
		return nil, err
	}

	for _, locStr := range location.Locations {
		// Analyser la chaîne "Ville,Pays"
		parts := strings.Split(locStr, ",")
		if len(parts) != 2 {
			continue // Ignorer les formats invalides
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

// Géocode une adresse unique
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
