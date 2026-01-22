package ui

import (
	"groupie-tracker/models"
	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
)

type AppWindow struct {
	App             fyne.App
	Window          fyne.Window
	AllArtists      []models.Artist
	EnrichedArtists []services.ArtistEnriched
	Favorites       *FavoritesManager
}

func NewApp(artists []models.Artist) *AppWindow {
	myApp := app.New()

	myWindow := myApp.NewWindow("Groupie Tracker")
	myWindow.Resize(fyne.NewSize(1400, 900))
	myWindow.CenterOnScreen()

	enrichedArtists := services.EnrichArtists(artists)

	return &AppWindow{
		App:             myApp,
		Window:          myWindow,
		AllArtists:      artists,
		EnrichedArtists: enrichedArtists,
		Favorites:       NewFavoritesManager(),
	}
}

func (w *AppWindow) ShowArtistList() {
	content := RenderArtistList(w.AllArtists, w)
	w.Window.SetContent(content)
}

func (w *AppWindow) ShowFilteredArtistList(filteredArtists []models.Artist) {
	content := RenderArtistList(filteredArtists, w)
	w.Window.SetContent(content)
}

func (w *AppWindow) ShowArtistDetail(artistName string) {
	content := RenderArtistDetail(artistName, w)
	scroll := container.NewVScroll(content)
	w.Window.SetContent(scroll)
}

func (w *AppWindow) Run() {
	w.SetupShortcuts()
	w.ShowArtistList()
	w.Window.ShowAndRun()
}

// Configuration des raccourcis clavier pour l'application
func (w *AppWindow) SetupShortcuts() {
	w.Window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyH,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		w.ShowArtistList()
	})

	w.Window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyEscape,
	}, func(shortcut fyne.Shortcut) {
		w.ShowArtistList()
	})

	w.Window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyR,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		w.ShowArtistList()
	})

	w.Window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		w.App.Quit()
	})
}
