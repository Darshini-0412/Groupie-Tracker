package services

import (
	"groupie-tracker/models"
	"strings"
)

// FilterByCreationDate filtre les artistes par année de création
func FilterByCreationDate(artists []models.Artist, min, max int) []models.Artist {
	var filtered []models.Artist
	for _, a := range artists {
		if a.CreationDate >= min && a.CreationDate <= max {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// FilterByMemberCount filtre par nombre de membres
func FilterByMemberCount(artists []models.Artist, min, max int) []models.Artist {
	var filtered []models.Artist
	for _, a := range artists {
		memberCount := len(a.Members)
		if memberCount >= min && memberCount <= max {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// FilterByFirstAlbumDate filtre par date du premier album
// Les dates sont au format "DD-MM-YYYY" donc on peut comparer directement
func FilterByFirstAlbumDate(artists []models.Artist, minDate, maxDate string) []models.Artist {
	var filtered []models.Artist
	for _, a := range artists {
		if a.FirstAlbum >= minDate && a.FirstAlbum <= maxDate {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// FilterByLocations filtre les artistes qui ont joué dans certains lieux
func FilterByLocations(artists []models.Artist, selectedLocations []string) []models.Artist {
	if len(selectedLocations) == 0 {
		return artists // Pas de filtre = tout
	}

	var filtered []models.Artist

	for _, artist := range artists {
		// Récupération de la relation pour avoir les lieux
		relation, err := FetchRelationByID(artist.ID)
		if err != nil {
			continue
		}

		hasMatchingLocation := false
		for location := range relation.DatesLocations {
			locationLower := strings.ToLower(location)

			// Check si un des lieux sélectionnés match
			for _, selected := range selectedLocations {
				selectedLower := strings.ToLower(selected)
				if strings.Contains(locationLower, selectedLower) {
					hasMatchingLocation = true
					break
				}
			}

			if hasMatchingLocation {
				break
			}
		}

		if hasMatchingLocation {
			filtered = append(filtered, artist)
		}
	}

	return filtered
}

// GetAllUniqueLocations récupère tous les lieux uniques de tous les artistes
// Pratique pour remplir un menu déroulant de filtres
func GetAllUniqueLocations(artists []models.Artist) []string {
	locationSet := make(map[string]bool)

	for _, artist := range artists {
		relation, err := FetchRelationByID(artist.ID)
		if err != nil {
			continue
		}

		for location := range relation.DatesLocations {
			cleanLocation := formatLocationName(location)
			locationSet[cleanLocation] = true
		}
	}

	// Conversion map → slice
	var locations []string
	for loc := range locationSet {
		locations = append(locations, loc)
	}

	return locations
}

// formatLocationName nettoie les noms de lieux
// Ex: "los_angeles-usa" → "Los Angeles, USA"
func formatLocationName(location string) string {
	// Remplacer underscores par espaces
	location = strings.ReplaceAll(location, "_", " ")

	// Séparer ville et pays (séparés par "-")
	parts := strings.Split(location, "-")

	if len(parts) == 2 {
		city := strings.Title(strings.TrimSpace(parts[0]))
		country := strings.ToUpper(strings.TrimSpace(parts[1]))
		return city + ", " + country
	}

	// Sinon juste capitaliser
	return strings.Title(location)
}
