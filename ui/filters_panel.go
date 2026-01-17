package ui

import (
	"fmt"
	"groupie-tracker/models"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateFiltersPanel(allArtists []models.Artist, w *AppWindow) *fyne.Container {
	titleText := canvas.NewText("🎯 Filtres", nil)
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 18

	yearLabel := widget.NewLabel("Année de création:")
	yearLabel.TextStyle = fyne.TextStyle{Bold: true}

	minYearEntry := widget.NewEntry()
	minYearEntry.SetPlaceHolder("Min")
	maxYearEntry := widget.NewEntry()
	maxYearEntry.SetPlaceHolder("Max")

	membersLabel := widget.NewLabel("Nombre de membres:")
	membersLabel.TextStyle = fyne.TextStyle{Bold: true}

	minMembersEntry := widget.NewEntry()
	minMembersEntry.SetPlaceHolder("Min")
	maxMembersEntry := widget.NewEntry()
	maxMembersEntry.SetPlaceHolder("Max")

	albumLabel := widget.NewLabel("Année premier album:")
	albumLabel.TextStyle = fyne.TextStyle{Bold: true}

	minAlbumEntry := widget.NewEntry()
	minAlbumEntry.SetPlaceHolder("Min")
	maxAlbumEntry := widget.NewEntry()
	maxAlbumEntry.SetPlaceHolder("Max")

	applyBtn := widget.NewButton("✓ Appliquer", func() {
		filtered := allArtists

		if minYearEntry.Text != "" && maxYearEntry.Text != "" {
			minYear, _ := strconv.Atoi(minYearEntry.Text)
			maxYear, _ := strconv.Atoi(maxYearEntry.Text)
			filtered = filterByCreationYear(filtered, minYear, maxYear)
		}

		if minMembersEntry.Text != "" && maxMembersEntry.Text != "" {
			minMembers, _ := strconv.Atoi(minMembersEntry.Text)
			maxMembers, _ := strconv.Atoi(maxMembersEntry.Text)
			filtered = filterByMembers(filtered, minMembers, maxMembers)
		}

		if minAlbumEntry.Text != "" && maxAlbumEntry.Text != "" {
			minAlbum, _ := strconv.Atoi(minAlbumEntry.Text)
			maxAlbum, _ := strconv.Atoi(maxAlbumEntry.Text)
			filtered = filterByFirstAlbum(filtered, minAlbum, maxAlbum)
		}

		w.ShowFilteredArtistList(filtered)
	})

	resetBtn := widget.NewButton("↺ Réinitialiser", func() {
		minYearEntry.SetText("")
		maxYearEntry.SetText("")
		minMembersEntry.SetText("")
		maxMembersEntry.SetText("")
		minAlbumEntry.SetText("")
		maxAlbumEntry.SetText("")
		w.ShowArtistList()
	})

	return container.NewVBox(
		titleText,
		widget.NewSeparator(),
		yearLabel,
		container.NewGridWithColumns(2, minYearEntry, maxYearEntry),
		widget.NewSeparator(),
		membersLabel,
		container.NewGridWithColumns(2, minMembersEntry, maxMembersEntry),
		widget.NewSeparator(),
		albumLabel,
		container.NewGridWithColumns(2, minAlbumEntry, maxAlbumEntry),
		widget.NewSeparator(),
		applyBtn,
		resetBtn,
	)
}

func filterByCreationYear(artists []models.Artist, min, max int) []models.Artist {
	var result []models.Artist
	for _, artist := range artists {
		if artist.CreationDate >= min && artist.CreationDate <= max {
			result = append(result, artist)
		}
	}
	return result
}

func filterByMembers(artists []models.Artist, min, max int) []models.Artist {
	var result []models.Artist
	for _, artist := range artists {
		count := len(artist.Members)
		if count >= min && count <= max {
			result = append(result, artist)
		}
	}
	return result
}

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

func extractYear(date string) (int, error) {
	if len(date) >= 4 {
		yearStr := date[len(date)-4:]
		return strconv.Atoi(yearStr)
	}
	return 0, fmt.Errorf("invalid date format")
}
