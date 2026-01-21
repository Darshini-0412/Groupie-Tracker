package services

import (
	"fmt"
	"groupie-tracker/models"
	"strconv"
	"strings"
)

type SearchSuggestion struct {
	Text   string
	Type   string
	Artist models.Artist
}

func SearchArtistsWithLocations(artistsEnriched []ArtistEnriched, query string) []SearchSuggestion {
	if query == "" {
		return nil
	}

	q := strings.ToLower(query)
	var results []SearchSuggestion
	seen := make(map[int]bool)

	for _, enriched := range artistsEnriched {
		a := enriched.Artist

		if !seen[a.ID] && strings.Contains(strings.ToLower(a.Name), q) {
			results = append(results, SearchSuggestion{a.Name, "Artiste", a})
			seen[a.ID] = true
			continue
		}

		if !seen[a.ID] {
			for _, m := range a.Members {
				if strings.Contains(strings.ToLower(m), q) {
					results = append(results, SearchSuggestion{fmt.Sprintf("%s → %s", m, a.Name), "Membre", a})
					seen[a.ID] = true
					break
				}
			}
		}

		if !seen[a.ID] && strings.Contains(strconv.Itoa(a.CreationDate), q) {
			results = append(results, SearchSuggestion{fmt.Sprintf("%s (%d)", a.Name, a.CreationDate), "Année", a})
			seen[a.ID] = true
			continue
		}

		if !seen[a.ID] && strings.Contains(strings.ToLower(a.FirstAlbum), q) {
			results = append(results, SearchSuggestion{fmt.Sprintf("%s (%s)", a.Name, a.FirstAlbum), "Album", a})
			seen[a.ID] = true
			continue
		}

		if !seen[a.ID] {
			for _, location := range enriched.Locations {
				if strings.Contains(strings.ToLower(location), q) {
					results = append(results, SearchSuggestion{fmt.Sprintf("%s (%s)", a.Name, location), "Lieu", a})
					seen[a.ID] = true
					break
				}
			}
		}

		if !seen[a.ID] {
			for _, date := range enriched.ConcertDates {
				if strings.Contains(date, q) {
					results = append(results, SearchSuggestion{fmt.Sprintf("%s (%s)", a.Name, date), "Date", a})
					seen[a.ID] = true
					break
				}
			}
		}
	}

	return results
}

func SearchArtists(artists []models.Artist, query string) []SearchSuggestion {
	if query == "" {
		return nil
	}

	q := strings.ToLower(query)
	var results []SearchSuggestion
	seen := make(map[int]bool)

	for _, a := range artists {
		if seen[a.ID] {
			continue
		}
		if strings.Contains(strings.ToLower(a.Name), q) {
			results = append(results, SearchSuggestion{a.Name, "Artiste", a})
			seen[a.ID] = true
			continue
		}
		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), q) {
				results = append(results, SearchSuggestion{fmt.Sprintf("%s → %s", m, a.Name), "Membre", a})
				seen[a.ID] = true
				break
			}
		}
		if seen[a.ID] {
			continue
		}
		if strings.Contains(strconv.Itoa(a.CreationDate), q) {
			results = append(results, SearchSuggestion{fmt.Sprintf("%s (%d)", a.Name, a.CreationDate), "Année", a})
			seen[a.ID] = true
			continue
		}
		if strings.Contains(strings.ToLower(a.FirstAlbum), q) {
			results = append(results, SearchSuggestion{fmt.Sprintf("%s (%s)", a.Name, a.FirstAlbum), "Album", a})
			seen[a.ID] = true
		}
	}

	return results
}
