package ui

import (
	"groupie-tracker/models"
)

// FavoritesManager gère les artistes favoris de l'utilisateur
type FavoritesManager struct {
	favorites map[int]bool // Clé = ID artiste, Valeur = true si favori
}

// NewFavoritesManager crée un nouveau gestionnaire de favoris
func NewFavoritesManager() *FavoritesManager {
	return &FavoritesManager{
		favorites: make(map[int]bool),
	}
}

// Add ajoute un artiste aux favoris
func (f *FavoritesManager) Add(artistID int) {
	f.favorites[artistID] = true
}

// Remove retire un artiste des favoris
func (f *FavoritesManager) Remove(artistID int) {
	delete(f.favorites, artistID)
}

// IsFavorite vérifie si un artiste est dans les favoris
func (f *FavoritesManager) IsFavorite(artistID int) bool {
	return f.favorites[artistID]
}

// Toggle ajoute ou retire un artiste des favoris (🤍 ↔ ❤️)
func (f *FavoritesManager) Toggle(artistID int) {
	if f.IsFavorite(artistID) {
		f.Remove(artistID)
	} else {
		f.Add(artistID)
	}
}

// GetFavorites retourne la liste des artistes favoris
func (f *FavoritesManager) GetFavorites(allArtists []models.Artist) []models.Artist {
	var favs []models.Artist
	for _, artist := range allArtists {
		if f.IsFavorite(artist.ID) {
			favs = append(favs, artist)
		}
	}
	return favs
}
