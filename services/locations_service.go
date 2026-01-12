package services

import "groupie-tracker/models"

func GetLocations() ([]models.Locations, error) {
	url := "http://groupietrackers.herokuapp.com/api/locations"
	return Fetch[[]models.Locations](url)
}
