package services

import (
	"groupie-tracker/models"
	"sync"
)

// ArtistEnriched c'est un artiste avec ses lieux et dates déjà chargés
type ArtistEnriched struct {
	Artist       models.Artist
	Locations    []string
	ConcertDates []string
}

// EnrichArtists charge tout en parallèle pour aller plus vite
func EnrichArtists(artists []models.Artist) []ArtistEnriched {
	enriched := make([]ArtistEnriched, len(artists))
	var wg sync.WaitGroup // pour attendre que tout finisse
	var mu sync.Mutex     // pour éviter les bugs quand plusieurs trucs écrivent en même temps

	for i, artist := range artists {
		wg.Add(1)
		go func(idx int, a models.Artist) {
			defer wg.Done()

			locations, _ := GetArtistLocations(a.ID)
			dates, _ := GetArtistConcertDates(a.ID)

			mu.Lock()
			enriched[idx] = ArtistEnriched{
				Artist:       a,
				Locations:    locations,
				ConcertDates: dates,
			}
			mu.Unlock()
		}(i, artist)
	}

	wg.Wait()
	return enriched
}

func EnrichArtist(artist models.Artist) ArtistEnriched {
	locations, _ := GetArtistLocations(artist.ID)
	dates, _ := GetArtistConcertDates(artist.ID)

	return ArtistEnriched{
		Artist:       artist,
		Locations:    locations,
		ConcertDates: dates,
	}
}
