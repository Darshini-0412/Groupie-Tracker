package ui

import (
	"groupie-tracker/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type AppWindow struct {
	App        fyne.App
	Window     fyne.Window
	AllArtists []models.Artist
}

func NewApp(artists []models.Artist) *AppWindow {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.CenterOnScreen()

	return &AppWindow{
		App:        myApp,
		Window:     myWindow,
		AllArtists: artists,
	}
}

func (w *AppWindow) ShowArtistList() {
	content := RenderArtistList(w.AllArtists, w)
	scroll := container.NewVScroll(content)
	w.Window.SetContent(scroll)
}

func (w *AppWindow) ShowFilteredArtistList(filteredArtists []models.Artist) {
	content := RenderArtistList(filteredArtists, w)
	scroll := container.NewVScroll(content)
	w.Window.SetContent(scroll)
}

func (w *AppWindow) ShowArtistDetail(artistName string) {
	content := RenderArtistDetail(artistName, w)
	scroll := container.NewVScroll(content)
	w.Window.SetContent(scroll)
}

func (w *AppWindow) Run() {
	w.Window.ShowAndRun()
}
