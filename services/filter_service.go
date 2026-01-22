package services

import (
	"groupie-tracker/models"
	"strings"
)

func FilterByCreationDate(artists []models.Artist, min, max int) []models.Artist {
	var filtered []models.Artist
	for _, a := range artists {
		if a.CreationDate >= min && a.CreationDate <= max {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

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

func FilterByFirstAlbumDate(artists []models.Artist, minDate, maxDate string) []models.Artist {
	var filtered []models.Artist
	for _, a := range artists {
		if a.FirstAlbum >= minDate && a.FirstAlbum <= maxDate {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// FilterByLocations garde seulement les artistes qui ont joué dans certains endroits
func FilterByLocations(artists []models.Artist, selectedLocations []string) []models.Artist {
	if len(selectedLocations) == 0 {
		return artists
	}

	var filtered []models.Artist

	for _, artist := range artists {
		relation, err := FetchRelationByID(artist.ID)
		if err != nil {
			continue
		}

		hasMatchingLocation := false
		for location := range relation.DatesLocations {
			locationLower := strings.ToLower(location)

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

// GetAllUniqueLocations récupère tous les lieux différents
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

	var locations []string
	for loc := range locationSet {
		locations = append(locations, loc)
	}

	return locations
}

// formatLocationName nettoie les noms de lieux
func formatLocationName(location string) string {
	location = strings.ReplaceAll(location, "_", " ")

	parts := strings.Split(location, "-")

	if len(parts) == 2 {
		city := strings.Title(strings.TrimSpace(parts[0]))
		country := strings.ToUpper(strings.TrimSpace(parts[1]))
		return city + ", " + country
	}

	return strings.Title(location)
}
