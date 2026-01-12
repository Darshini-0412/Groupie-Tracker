package services

import "groupie-tracker/models"

func GetRelations() ([]models.Relation, error) {
	url := "http://groupietrackers.herokuapp.com/api/relation"
	return Fetch[[]models.Relation](url)
}
