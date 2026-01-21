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

// RenderArtistDetail affiche la page détails d'un artiste
func RenderArtistDetail(artistName string, w *AppWindow) *fyne.Container {
	var artist models.Artist
	found := false

	// Chercher l'artiste par son nom
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

	// Lieu sélectionné sur la carte
	var selectedLocation string

	// Bouton retour
	backBtn := widget.NewButton("← Retour", func() {
		w.ShowArtistList()
	})

	// Image de l'artiste
	img := loadImageFromURL(artist.Image)
	img.SetMinSize(fyne.NewSize(300, 300))

	// Titre (nom de l'artiste)
	title := canvas.NewText(artist.Name, textGray)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	// En-tête de la page
	header := container.NewVBox(
		container.NewHBox(backBtn),
		container.NewCenter(img),
		container.NewCenter(title),
	)

	// Créer les cartes
	mapCard, refreshMap := makeMapCard(artist, &selectedLocation)
	concertsCard := makeConcertsCard(artist, &selectedLocation, refreshMap)

	// Assembler tout le contenu
	content := container.NewVBox(
		header,
		makeDetailInfoCard(artist),
		makeDetailMembersCard(artist),
		concertsCard,
		mapCard,
	)

	return content
}

// makeMapCard crée la carte interactive des concerts
func makeMapCard(artist models.Artist, selected *string) (*fyne.Container, func()) {

	bg := canvas.NewRectangle(bgCard)

	titleLabel := canvas.NewText("🗺️ Carte des concerts", textGray)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = 20

	// Récupérer les lieux de concerts
	relation, err := services.FetchRelationByID(artist.ID)
	if err != nil {
		errorLabel := widget.NewLabel("Impossible de charger la carte")
		return container.NewStack(bg, errorLabel), func() {}
	}

	// Séparer concerts passés et futurs
	past, future := services.SplitPastFutureConcerts(*relation)

	mapBox := container.NewVBox()

	// Fonction pour rafraîchir la carte
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

	// Récupérer les concerts
	relation, err := services.FetchRelationByID(artist.ID)
	if err != nil {
		return container.NewStack(bg, widget.NewLabel("Erreur concerts"))
	}

	list := container.NewVBox()

	// Afficher chaque lieu avec ses dates
	for location, dates := range relation.DatesLocations {
		loc := location

		// Bouton cliquable pour chaque lieu
		btn := widget.NewButton("📍 "+loc, func() {
			*selected = loc
			refreshMap() // Rafraîchir la carte
		})
		btn.Importance = widget.LowImportance
		list.Add(btn)

		// Ajouter les dates du concert
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

// makeDetailInfoCard crée la carte d'informations générales
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

// makeDetailMembersCard crée la carte des membres du groupe
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

// loadImageFromURL télécharge et affiche une image depuis une URL
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
