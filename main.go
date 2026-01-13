package main

import (
	"groupie-tracker/ui"
)

func main() {
	app := ui.NewApp()
	app.ShowArtistList(nil)
	app.Run()
}
