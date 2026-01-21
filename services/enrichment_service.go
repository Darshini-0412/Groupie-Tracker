package services

import (
	"groupie-tracker/models"
	"sync"
)

type ArtistEnriched struct {
	Artist       models.Artist
	Locations    []string
	ConcertDates []string
}

func EnrichArtists(artists []models.Artist) []ArtistEnriched {
	enriched := make([]ArtistEnriched, len(artists))
	var wg sync.WaitGroup
	var mu sync.Mutex

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
