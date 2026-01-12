package main

import (
	"groupie-tracker/services"
	"groupie-tracker/ui"
	"log"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	artists, err := services.GetArtists()
	if err != nil {
		log.Fatal(err)
	}
	ui.ShowArtislList(w, artists)
	w.ShowAndRun()
}
