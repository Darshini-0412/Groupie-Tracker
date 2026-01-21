package models

// Location stocke tous les lieux de concerts
type Location struct {
	ID        int
	Locations []string `json:"locations"` // Ex: "paris-france", "london-uk"
}
