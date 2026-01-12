package services

import "groupie-tracker/models"

func GetDates() ([]models.Dates, error) {
	url := "http://groupietrackers.herokuapp.com/api/dates"
	return Fetch[[]models.Dates](url)
}
