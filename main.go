package main

import (
	"groupie-tracker/services"
	"groupie-tracker/ui"
	"log"
)

func main() {
	artists, err := services.FetchArtists()
	if err != nil {
		log.Fatal("Erreur lors du chargement des artistes:", err)
	}

	app := ui.NewApp(artists)
	app.ShowArtistList()
	app.Run()
}
