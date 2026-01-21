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

// CreateFiltersPanel crée le panneau de filtres à gauche
func CreateFiltersPanel(allArtists []models.Artist, w *AppWindow) *fyne.Container {
	titleText := canvas.NewText("🎯 Filtres", nil)
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 18

	// Filtre: Année de création
	yearLabel := widget.NewLabel("Année de création:")
	yearLabel.TextStyle = fyne.TextStyle{Bold: true}

	minYearEntry := widget.NewEntry()
	minYearEntry.SetPlaceHolder("Min")
	maxYearEntry := widget.NewEntry()
	maxYearEntry.SetPlaceHolder("Max")

	// Filtre: Année premier album
	albumLabel := widget.NewLabel("Année premier album:")
	albumLabel.TextStyle = fyne.TextStyle{Bold: true}

	minAlbumEntry := widget.NewEntry()
	minAlbumEntry.SetPlaceHolder("Min")
	maxAlbumEntry := widget.NewEntry()
	maxAlbumEntry.SetPlaceHolder("Max")

	// Filtre: Lieu des concerts
	locationLabel := widget.NewLabel("Lieu des concerts:")
	locationLabel.TextStyle = fyne.TextStyle{Bold: true}

	locationEntry := widget.NewEntry()
	locationEntry.SetPlaceHolder("Ex: paris, usa, london...")

	// Filtre: Nombre de membres (1-8 checkboxes)
	membersLabel := widget.NewLabel("Nombre de membres:")
	membersLabel.TextStyle = fyne.TextStyle{Bold: true}

	memberChecks := make([]*widget.Check, 0)
	memberContainer := container.NewVBox()

	for i := 1; i <= 8; i++ {
		num := i
		check := widget.NewCheck(fmt.Sprintf("%d", num), nil)
		memberChecks = append(memberChecks, check)
		memberContainer.Add(check)
	}

	// Bouton Appliquer
	applyBtn := widget.NewButton("✓ Appliquer", func() {
		filtered := allArtists

		// Appliquer filtre année de création
		if minYearEntry.Text != "" && maxYearEntry.Text != "" {
			minYear, _ := strconv.Atoi(minYearEntry.Text)
			maxYear, _ := strconv.Atoi(maxYearEntry.Text)
			filtered = filterByCreationYear(filtered, minYear, maxYear)
		}

		// Appliquer filtre nombre de membres
		selectedMembers := []int{}
		for i, check := range memberChecks {
			if check.Checked {
				selectedMembers = append(selectedMembers, i+1)
			}
		}
		if len(selectedMembers) > 0 {
			filtered = filterByMembersCheckbox(filtered, selectedMembers)
		}

		// Appliquer filtre album
		if minAlbumEntry.Text != "" && maxAlbumEntry.Text != "" {
			minAlbum, _ := strconv.Atoi(minAlbumEntry.Text)
			maxAlbum, _ := strconv.Atoi(maxAlbumEntry.Text)
			filtered = filterByFirstAlbum(filtered, minAlbum, maxAlbum)
		}

		// Appliquer filtre lieu
		if locationEntry.Text != "" {
			locationQuery := strings.TrimSpace(locationEntry.Text)
			filtered = filterByLocationText(filtered, locationQuery)
		}

		// Afficher résultats filtrés
		w.ShowFilteredArtistList(filtered)
	})

	// Bouton Réinitialiser
	resetBtn := widget.NewButton("↺ Réinitialiser", func() {
		minYearEntry.SetText("")
		maxYearEntry.SetText("")
		minAlbumEntry.SetText("")
		maxAlbumEntry.SetText("")
		locationEntry.SetText("")

		for _, check := range memberChecks {
			check.SetChecked(false)
		}

		w.ShowArtistList()
	})

	return container.NewVBox(
		titleText,
		widget.NewSeparator(),
		yearLabel,
		container.NewGridWithColumns(2, minYearEntry, maxYearEntry),
		widget.NewSeparator(),
		albumLabel,
		container.NewGridWithColumns(2, minAlbumEntry, maxAlbumEntry),
		widget.NewSeparator(),
		locationLabel,
		locationEntry,
		widget.NewSeparator(),
		membersLabel,
		memberContainer,
		widget.NewSeparator(),
		applyBtn,
		resetBtn,
	)
}

// filterByCreationYear filtre par année de création (ex: 1990-2000)
func filterByCreationYear(artists []models.Artist, min, max int) []models.Artist {
	var result []models.Artist
	for _, artist := range artists {
		if artist.CreationDate >= min && artist.CreationDate <= max {
			result = append(result, artist)
		}
	}
	return result
}

// filterByMembersCheckbox filtre par nombre de membres sélectionnés
func filterByMembersCheckbox(artists []models.Artist, selectedCounts []int) []models.Artist {
	var result []models.Artist
	for _, artist := range artists {
		count := len(artist.Members)
		for _, selected := range selectedCounts {
			if count == selected {
				result = append(result, artist)
				break
			}
		}
	}
	return result
}

// filterByFirstAlbum filtre par année du premier album
func filterByFirstAlbum(artists []models.Artist, min, max int) []models.Artist {
	var result []models.Artist
	for _, artist := range artists {
		year, err := extractYear(artist.FirstAlbum)
		if err == nil && year >= min && year <= max {
			result = append(result, artist)
		}
	}
	return result
}

// filterByLocationText filtre par lieu de concert (ex: "paris", "usa")
func filterByLocationText(artists []models.Artist, query string) []models.Artist {
	query = strings.ToLower(strings.TrimSpace(query))
	var result []models.Artist

	for _, artist := range artists {
		relation, err := services.FetchRelationByID(artist.ID)
		if err != nil {
			continue
		}

		found := false
		for location := range relation.DatesLocations {
			locationLower := strings.ToLower(location)
			if strings.Contains(locationLower, query) {
				found = true
				break
			}
		}

		if found {
			result = append(result, artist)
		}
	}

	return result
}

// extractYear extrait l'année d'une date (ex: "01-02-2006" → 2006)
func extractYear(date string) (int, error) {
	if len(date) >= 4 {
		yearStr := date[len(date)-4:] // Les 4 derniers caractères
		return strconv.Atoi(yearStr)
	}
	return 0, fmt.Errorf("invalid date format")
}
