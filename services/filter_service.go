package services

import "groupie-tracker/models"

func FilterByCreationDate(artists []models.Artist, min, max int) []models.Artist {
	var result []models.Artist

	for _, artist := range artists {
		if artist.CreationDate >= min && artist.CreationDate <= max {
			result = append(result, artist)
		}
	}

	return result
}

func FilterByMemberCount(artists []models.Artist, min, max int) []models.Artist {
	var result []models.Artist

	for _, artist := range artists {
		count := len(artist.Members)
		if count >= min && count <= max {
			result = append(result, artist)
		}
	}

	return result
}
