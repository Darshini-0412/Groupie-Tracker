package services

import (
	"groupie-tracker/models"
	"sync"
)

// ArtistEnriched = Artist + ses lieux + ses dates
// Ça évite de refaire des appels API à chaque fois
type ArtistEnriched struct {
	Artist       models.Artist
	Locations    []string
	ConcertDates []string
}

// EnrichArtists enrichit TOUS les artistes en parallèle (goroutines)
// Ça accélère énormément le chargement initial
func EnrichArtists(artists []models.Artist) []ArtistEnriched {
	enriched := make([]ArtistEnriched, len(artists))
	var wg sync.WaitGroup
	var mu sync.Mutex // Pour éviter les race conditions

	for i, artist := range artists {
		wg.Add(1)
		go func(idx int, a models.Artist) {
			defer wg.Done()

			// Récupération des infos supplémentaires
			locations, _ := GetArtistLocations(a.ID)
			dates, _ := GetArtistConcertDates(a.ID)

			// Protection de l'écriture dans le slice partagé
			mu.Lock()
			enriched[idx] = ArtistEnriched{
				Artist:       a,
				Locations:    locations,
				ConcertDates: dates,
			}
			mu.Unlock()
		}(i, artist)
	}

	wg.Wait() // Attendre que toutes les goroutines finissent
	return enriched
}

// EnrichArtist enrichit UN SEUL artiste (version synchrone)
func EnrichArtist(artist models.Artist) ArtistEnriched {
	locations, _ := GetArtistLocations(artist.ID)
	dates, _ := GetArtistConcertDates(artist.ID)

	return ArtistEnriched{
		Artist:       artist,
		Locations:    locations,
		ConcertDates: dates,
	}
}
