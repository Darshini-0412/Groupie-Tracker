package ui

import (
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SearchSuggestion struct {
	Text   string
	Type   string
	Artist models.Artist
}

// CreateSmartSearchBar crée la barre de recherche avec suggestions
func CreateSmartSearchBar(w *AppWindow) *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("🔍 Rechercher...")

	suggestionsBox := container.NewVBox()
	suggestionsScroll := container.NewVScroll(suggestionsBox)
	suggestionsScroll.SetMinSize(fyne.NewSize(0, 200))
	suggestionsScroll.Hide()

	// quand on tape dans la barre
	searchEntry.OnChanged = func(query string) {
		query = strings.TrimSpace(query)

		// faut au moins 2 lettres
		if len(query) < 2 {
			suggestionsBox.Objects = nil
			suggestionsScroll.Hide()
			return
		}

		suggestions := services.SearchArtistsWithLocations(w.EnrichedArtists, query)
		suggestionsBox.Objects = nil

		// max 8 suggestions
		if len(suggestions) > 8 {
			suggestions = suggestions[:8]
		}

		for _, s := range suggestions {
			suggestion := s
			btn := widget.NewButton(
				fmt.Sprintf("%s %s", getIcon(s.Type), s.Text),
				func() {
					w.ShowArtistDetail(suggestion.Artist.Name)
					searchEntry.SetText("")
					suggestionsScroll.Hide()
				},
			)
			suggestionsBox.Add(btn)
		}

		if len(suggestions) > 0 {
			suggestionsScroll.Show()
		}
		suggestionsBox.Refresh()
	}

	bg := canvas.NewRectangle(bgCard)
	return container.NewStack(bg, container.NewPadded(container.NewVBox(searchEntry, suggestionsScroll)))
}

func getIcon(t string) string {
	icons := map[string]string{
		"Artiste": "🎤",
		"Membre":  "👤",
		"Année":   "📅",
		"Album":   "💿",
		"Lieu":    "📍",
		"Date":    "🗓️",
	}
	if icon, ok := icons[t]; ok {
		return icon
	}
	return "✨"
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

		if strings.Contains(strconv.Itoa(a.CreationDate), query) {
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
