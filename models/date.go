package models

// Date contient toutes les dates de concerts d'un artiste
type Date struct {
	ID    int
	Dates []string `json:"dates"` // Liste des dates au format "DD-MM-YYYY"
}
