package models

// Artist représente un artiste ou groupe de musique
type Artist struct {
	ID           int      `json:"id"`           // Identifiant unique
	Image        string   `json:"image"`        // URL de la photo
	Name         string   `json:"name"`         // Nom du groupe/artiste
	Members      []string `json:"members"`      // Liste des membres
	CreationDate int      `json:"creationDate"` // Année de création
	FirstAlbum   string   `json:"firstAlbum"`   // Date du premier album
	Locations    string   `json:"locations"`    // URL vers les lieux de concerts
	ConcertDates string   `json:"concertDates"` // URL vers les dates
	Relations    string   `json:"relations"`    // URL vers les relations dates/lieux
}
