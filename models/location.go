package models

type Location struct {
	ID        int
	Locations []string `json:"locations"`
}
