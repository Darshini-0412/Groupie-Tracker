package ui

import (
	"groupie-tracker/models"
	"groupie-tracker/services"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func RenderFiltersPanel(artists []models.Artist, w *AppWindow) *fyne.Container {
	title := widget.NewLabel("Filtres")
	title.TextStyle = fyne.TextStyle{Bold: true}

	minYear := widget.NewEntry()
	minYear.SetPlaceHolder("Min")
	maxYear := widget.NewEntry()
	maxYear.SetPlaceHolder("Max")

	minMembers := widget.NewEntry()
	minMembers.SetPlaceHolder("Min")
	maxMembers := widget.NewEntry()
	maxMembers.SetPlaceHolder("Max")

	applyBtn := widget.NewButton("Appliquer", func() {
		filtered := artists

		if minYear.Text != "" && maxYear.Text != "" {
			min, _ := strconv.Atoi(minYear.Text)
			max, _ := strconv.Atoi(maxYear.Text)
			filtered = services.FilterByCreationDate(filtered, min, max)
		}

		if minMembers.Text != "" && maxMembers.Text != "" {
			min, _ := strconv.Atoi(minMembers.Text)
			max, _ := strconv.Atoi(maxMembers.Text)
			filtered = services.FilterByMemberCount(filtered, min, max)
		}

		w.ShowFilteredArtistList(filtered)
	})

	resetBtn := widget.NewButton("Reset", func() {
		minYear.SetText("")
		maxYear.SetText("")
		minMembers.SetText("")
		maxMembers.SetText("")
		w.ShowArtistList()
	})

	return container.NewVBox(
		title,
		widget.NewLabel("Année:"),
		container.NewHBox(minYear, maxYear),
		widget.NewLabel("Membres:"),
		container.NewHBox(minMembers, maxMembers),
		applyBtn,
		resetBtn,
	)
}
