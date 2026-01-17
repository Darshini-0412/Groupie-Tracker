package ui

import (
	"groupie-tracker/models"
)

type FavoritesManager struct {
	favorites map[int]bool
}

func NewFavoritesManager() *FavoritesManager {
	return &FavoritesManager{
		favorites: make(map[int]bool),
	}
}

func (f *FavoritesManager) Add(artistID int) {
	f.favorites[artistID] = true
}

func (f *FavoritesManager) Remove(artistID int) {
	delete(f.favorites, artistID)
}

func (f *FavoritesManager) IsFavorite(artistID int) bool {
	return f.favorites[artistID]
}

func (f *FavoritesManager) Toggle(artistID int) {
	if f.IsFavorite(artistID) {
		f.Remove(artistID)
	} else {
		f.Add(artistID)
	}
}

func (f *FavoritesManager) GetFavorites(allArtists []models.Artist) []models.Artist {
	var favs []models.Artist
	for _, artist := range allArtists {
		if f.IsFavorite(artist.ID) {
			favs = append(favs, artist)
		}
	}
	return favs
}
