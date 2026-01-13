package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type AppWindow struct {
	App    fyne.App
	Window fyne.Window
}

func NewAppWindow() *AppWindow {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")
	return &AppWindow{
		App:    a,
		Window: w,
	}
}
func (aw *AppWindow) Show() {
	aw.Window.ShowAndRun()
}
