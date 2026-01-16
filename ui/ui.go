package ui

import (
	"fmt"
	"groupie-tracker/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func ShowArtistList(w fyne.Window, artists []models.Artist) {
	list := widget.NewList(
		func() int {
			return len(artists)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("%s - %s", artists[i].Name, artists[i].Image))
		},
	)
	w.SetContent(list)
}
