package services

import "groupie-tracker/models"

func GetArtists() ([]models.Artist, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	return Fetch[[]models.Artist](url)
}
