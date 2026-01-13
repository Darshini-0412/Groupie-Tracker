package ui

import (
	"groupie-tracker/models"
	"groupie-tracker/services"

	"fyne.io/fyne/v2/widget"
)

func RenderSearchBar(artists []models.Artist, w *AppWindow) *widget.Entry {
	search := widget.NewEntry()
	search.SetPlaceHolder("Rechercher...")

	search.OnChanged = func(query string) {
		results := services.SearchArtists(artists, query)
		w.ShowArtistList(results)
	}

	return search
}
