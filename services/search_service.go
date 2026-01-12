package services

import (
	"groupie-tracker/models"
	"strings"
)

func SearchArtists(artists []models.Artist, query string) []models.Artist {
	query = strings.ToLower(query)
	var result []models.Artist
	for _, a := range artists {
		if strings.Contains(strings.ToLower(a.Name), query) {
			result = append(result, a)
			continue
		}

		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), query) {
				result = append(result, a)
				break
			}
		}
	}

	return result
}
