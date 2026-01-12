package services

import (
	"groupie-tracker/models"
	"strconv"
)

func FilterByCreationYear(artists []models.Artist, min, max int) []models.Artist {
	var result []models.Artist
	for _, a := range artists {
		if a.CreationDate >= min && a.CreationDate <= max {
			result = append(result, a)
		}
	}
	return result
}

func FilterByMemberCount(artists []models.Artist, count int) []models.Artist {
	var result []models.Artist
	for _, a := range artists {
		if len(a.Members) == count {
			result = append(result, a)
		}
	}
	return result
}

func FilterByFirstAlbumYear(artists []models.Artist, min, max int) []models.Artist {
	var result []models.Artist

	for _, a := range artists {
		yearStr := a.FirstAlbum[:4]
		year, _ := strconv.Atoi(yearStr)

		if year >= min && year <= max {
			result = append(result, a)
		}
	}

	return result
}
