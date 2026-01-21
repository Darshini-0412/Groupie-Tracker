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

// SearchSuggestion = une suggestion de recherche
type SearchSuggestion struct {
	Text   string        // Texte affiché (ex: "Queen (1970)")
	Type   string        // Type de recherche (Artiste, Membre, Année, etc.)
	Artist models.Artist // Artiste correspondant
}

// CreateSmartSearchBar crée la barre de recherche intelligente avec suggestions
func CreateSmartSearchBar(w *AppWindow) *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("🔍 Rechercher...")

	// Conteneur pour les suggestions
	suggestionsBox := container.NewVBox()
	suggestionsScroll := container.NewVScroll(suggestionsBox)
	suggestionsScroll.SetMinSize(fyne.NewSize(0, 200))
	suggestionsScroll.Hide()

	// Quand l'utilisateur tape dans la barre de recherche
	searchEntry.OnChanged = func(query string) {
		query = strings.TrimSpace(query)

		// Minimum 2 caractères pour déclencher la recherche
		if len(query) < 2 {
			suggestionsBox.Objects = nil
			suggestionsScroll.Hide()
			return
		}

		// Recherche avancée (artiste, membre, lieu, date)
		suggestions := services.SearchArtistsWithLocations(w.EnrichedArtists, query)
		suggestionsBox.Objects = nil

		// Limiter à 8 suggestions maximum
		if len(suggestions) > 8 {
			suggestions = suggestions[:8]
		}

		// Créer un bouton pour chaque suggestion
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

		// Afficher les suggestions si on en a
		if len(suggestions) > 0 {
			suggestionsScroll.Show()
		}
		suggestionsBox.Refresh()
	}

	bg := canvas.NewRectangle(bgCard)
	return container.NewStack(bg, container.NewPadded(container.NewVBox(searchEntry, suggestionsScroll)))
}

// getIcon retourne l'icône selon le type de recherche
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

// SearchArtists recherche basique d'artistes
func SearchArtists(artists []models.Artist, query string) []SearchSuggestion {
	if query == "" {
		return nil
	}

	q := strings.ToLower(query)
	var results []SearchSuggestion
	seen := make(map[int]bool) // Éviter les doublons

	for _, a := range artists {
		if seen[a.ID] {
			continue
		}

		// Recherche par nom d'artiste
		if strings.Contains(strings.ToLower(a.Name), q) {
			results = append(results, SearchSuggestion{a.Name, "Artiste", a})
			seen[a.ID] = true
			continue
		}

		// Recherche par nom de membre
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

		// Recherche par année de création
		if strings.Contains(strconv.Itoa(a.CreationDate), query) {
			results = append(results, SearchSuggestion{fmt.Sprintf("%s (%d)", a.Name, a.CreationDate), "Année", a})
			seen[a.ID] = true
			continue
		}

		// Recherche par album
		if strings.Contains(strings.ToLower(a.FirstAlbum), q) {
			results = append(results, SearchSuggestion{fmt.Sprintf("%s (%s)", a.Name, a.FirstAlbum), "Album", a})
			seen[a.ID] = true
		}
	}

	return results
}
