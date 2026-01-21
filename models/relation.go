package models

// Relation fait le lien entre les lieux et les dates
// C'est la structure la plus importante pour afficher les concerts
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"dates_locations"` // Clé = lieu, Valeur = dates
}
