package main

import (
	"groupie-tracker/services"
	"groupie-tracker/ui"
	"log"
)

func main() {
	artists, err := services.FetchArtists()
	if err != nil {
		log.Fatal("Erreur chargement:", err)
	}

	app := ui.NewApp(artists)
	app.Run()
}
