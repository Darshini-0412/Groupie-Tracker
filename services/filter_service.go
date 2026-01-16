package services

import (
	"groupie-tracker/models"
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
