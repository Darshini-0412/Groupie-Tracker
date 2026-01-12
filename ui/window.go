package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AppWindow struct {
	App    fyne.App
	Window fyne.Window
}

func SetupWindow() *AppWindow {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.CenterOnScreen()

	welcome := widget.NewLabel("Bienvenue dans Groupie Tracker")
	welcome.Alignment = fyne.TextAlignCenter
	
	content := container.NewCenter(welcome)
	myWindow.SetContent(content)

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