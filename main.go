package main

import (
	"groupie-tracker/services"
	"groupie-tracker/ui"
	"log"
)

func main() {
	// récupérer les artistes depuis l'API
	artists, err := services.FetchArtists()
	if err != nil {
		log.Fatal("Erreur chargement:", err)
	}

	// lancer l'appli
	app := ui.NewApp(artists)
	app.Run()
}
