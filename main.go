package main

import (
	"groupie-tracker/services"
	"groupie-tracker/ui"
	"log"
)

func main() {
	// On récupère tous les artistes depuis l'API
	artists, err := services.FetchArtists()
	if err != nil {
		log.Fatal("Erreur chargement:", err)
	}

	// On démarre l'application avec la liste complète
	app := ui.NewApp(artists)
	app.Run()
}
