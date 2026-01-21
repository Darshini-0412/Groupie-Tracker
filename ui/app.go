package ui

import (
	"groupie-tracker/models"
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

// AppWindow = fenêtre principale de l'application
type AppWindow struct {
	App             fyne.App
	Window          fyne.Window
	AllArtists      []models.Artist
	EnrichedArtists []services.ArtistEnriched
	Favorites       *FavoritesManager
}

// NewApp crée et configure la fenêtre principale
func NewApp(artists []models.Artist) *AppWindow {
	myApp := app.New()

	// Créer fenêtre 1400x900 centrée
	myWindow := myApp.NewWindow("Groupie Tracker")
	myWindow.Resize(fyne.NewSize(1400, 900))
	myWindow.CenterOnScreen()

	// Enrichir les artistes avec lieux et dates
	enrichedArtists := services.EnrichArtists(artists)

	return &AppWindow{
		App:             myApp,
		Window:          myWindow,
		AllArtists:      artists,
		EnrichedArtists: enrichedArtists,
		Favorites:       NewFavoritesManager(),
	}
}

// ShowArtistList affiche la liste complète
func (w *AppWindow) ShowArtistList() {
	content := RenderArtistList(w.AllArtists, w)
	w.Window.SetContent(content)
}

// ShowFilteredArtistList affiche les résultats filtrés
func (w *AppWindow) ShowFilteredArtistList(filteredArtists []models.Artist) {
	content := RenderArtistList(filteredArtists, w)
	w.Window.SetContent(content)
}

// ShowArtistDetail affiche les détails d'un artiste
func (w *AppWindow) ShowArtistDetail(artistName string) {
	content := RenderArtistDetail(artistName, w)
	scroll := container.NewVScroll(content)
	w.Window.SetContent(scroll)
}

// Run lance l'application
func (w *AppWindow) Run() {
	w.ShowArtistList()
	w.Window.ShowAndRun()
}
