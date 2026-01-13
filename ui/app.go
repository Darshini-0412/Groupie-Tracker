package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type AppWindow struct {
	App    fyne.App
	Window fyne.Window
}

func NewApp() *AppWindow {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.CenterOnScreen()

	return &AppWindow{
		App:    myApp,
		Window: myWindow,
	}
}

func (w *AppWindow) ShowArtistList(artists interface{}) {
	content := RenderArtistList(artists, w)
	w.Window.SetContent(content)
}

func (w *AppWindow) ShowArtistDetail(artist interface{}) {
	content := RenderArtistDetail(artist, w)
	w.Window.SetContent(content)
}

func (w *AppWindow) Run() {
	w.Window.ShowAndRun()
}
