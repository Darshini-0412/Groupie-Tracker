package ui

import (
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// RenderArtistDetail affiche la page de détails d'un artiste
func RenderArtistDetail(artistName string, w *AppWindow) *fyne.Container {
	var artist models.Artist
	found := false

	for _, a := range w.AllArtists {
		if a.Name == artistName {
			artist = a
			found = true
			break
		}
	}

	if !found {
		return container.NewVBox(widget.NewLabel("Artiste non trouvé"))
	}

	var selectedLocation string

	backBtn := widget.NewButton("← Retour", func() {
		w.ShowArtistList()
	})

	img := loadImageFromURL(artist.Image)
	img.SetMinSize(fyne.NewSize(300, 300))

	title := canvas.NewText(artist.Name, textGray)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	header := container.NewVBox(
		container.NewHBox(backBtn),
		container.NewCenter(img),
		container.NewCenter(title),
	)

	mapCard, refreshMap := makeMapCard(artist, &selectedLocation)
	concertsCard := makeConcertsCard(artist, &selectedLocation, refreshMap)

	content := container.NewVBox(
		header,
		makeDetailInfoCard(artist),
		makeDetailMembersCard(artist),
		concertsCard,
		mapCard,
	)

	return content
}

// makeMapCard crée la carte interactive
func makeMapCard(artist models.Artist, selected *string) (*fyne.Container, func()) {

	bg := canvas.NewRectangle(bgCard)

	titleLabel := canvas.NewText("🗺️ Carte des concerts", textGray)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = 20

	relation, err := services.FetchRelationByID(artist.ID)
	if err != nil {
		errorLabel := widget.NewLabel("Impossible de charger la carte")
		return container.NewStack(bg, errorLabel), func() {}
	}

	past, future := services.SplitPastFutureConcerts(*relation)

	mapBox := container.NewVBox()

	refresh := func() {
		mapBox.Objects = nil
		mapBox.Objects = []fyne.CanvasObject{
			RenderMap(past, future, selected),
		}
		mapBox.Refresh()
	}

	refresh()

	content := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		mapBox,
	)

	return container.NewStack(bg, container.NewPadded(content)), refresh
}

// makeConcertsCard crée la liste des concerts
func makeConcertsCard(artist models.Artist, selected *string, refreshMap func()) *fyne.Container {

	bg := canvas.NewRectangle(bgCard)

	titleLabel := canvas.NewText("🎤 Dates et Lieux des Concerts", textGray)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = 20

	relation, err := services.FetchRelationByID(artist.ID)
	if err != nil {
		return container.NewStack(bg, widget.NewLabel("Erreur concerts"))
	}

	list := container.NewVBox()

	for location, dates := range relation.DatesLocations {
		loc := location

		btn := widget.NewButton("📍 "+loc, func() {
			*selected = loc
			refreshMap()
		})
		btn.Importance = widget.LowImportance
		list.Add(btn)

		for _, date := range dates {
			list.Add(widget.NewLabel("   📅 " + date))
		}

		list.Add(widget.NewSeparator())
	}

	scroll := container.NewVScroll(list)
	scroll.SetMinSize(fyne.NewSize(0, 300))

	content := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		scroll,
	)

	return container.NewStack(bg, container.NewPadded(content))
}

func makeDetailInfoCard(artist models.Artist) *fyne.Container {
	bg := canvas.NewRectangle(bgCard)

	titleLabel := canvas.NewText("📋 Informations", textGray)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = 20

	info := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		canvas.NewText(fmt.Sprintf("📅 Année: %d", artist.CreationDate), textGray),
		canvas.NewText(fmt.Sprintf("💿 Premier album: %s", artist.FirstAlbum), textGray),
		canvas.NewText(fmt.Sprintf("👥 Membres: %d", len(artist.Members)), textGray),
	)

	return container.NewStack(bg, container.NewPadded(info))
}

func makeDetailMembersCard(artist models.Artist) *fyne.Container {
	bg := canvas.NewRectangle(bgCard)

	titleLabel := canvas.NewText("👥 Membres du groupe", textGray)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = 20

	list := container.NewVBox()
	for _, m := range artist.Members {
		list.Add(canvas.NewText("• "+m, textGray))
	}

	return container.NewStack(
		bg,
		container.NewPadded(container.NewVBox(titleLabel, widget.NewSeparator(), list)),
	)
}

// loadImageFromURL télécharge une image depuis une URL
func loadImageFromURL(url string) *canvas.Image {
	resp, err := http.Get(url)
	if err != nil {
		return canvas.NewImageFromImage(nil)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return canvas.NewImageFromImage(nil)
	}

	img := &canvas.Image{Resource: fyne.NewStaticResource(url, data)}
	img.FillMode = canvas.ImageFillContain
	return img
}
