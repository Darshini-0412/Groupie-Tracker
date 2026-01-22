package models

// Relation fait le lien entre les lieux et les dates
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"dates_locations"`
}
