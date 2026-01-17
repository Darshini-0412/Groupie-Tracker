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

	backBtn := widget.NewButton("← Retour", func() {
		w.ShowArtistList()
	})

	favoriteIcon := "🤍"
	if w.Favorites.IsFavorite(artist.ID) {
		favoriteIcon = "❤️"
	}

	favoriteBtn := widget.NewButton(favoriteIcon+" Ma Sélection", func() {
		w.Favorites.Toggle(artist.ID)
		w.ShowArtistDetail(artistName)
	})

	img := loadImageFromURL(artist.Image)
	img.SetMinSize(fyne.NewSize(300, 300))

	title := canvas.NewText(artist.Name, textGray)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	header := container.NewVBox(
		container.NewHBox(backBtn, favoriteBtn),
		container.NewCenter(img),
		container.NewCenter(title),
	)

	infoCard := makeDetailInfoCard(artist)
	membersCard := makeDetailMembersCard(artist)
	concertsCard := makeConcertsCard(artist)

	content := container.NewVBox(
		header,
		infoCard,
		membersCard,
		concertsCard,
	)

	return content
}

func makeConcertsCard(artist models.Artist) *fyne.Container {
	bg := canvas.NewRectangle(bgCard)

	titleLabel := canvas.NewText("🎤 Dates et Lieux des Concerts", textGray)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = 20

	relation, err := services.FetchRelationByID(artist.ID)
	if err != nil {
		errorLabel := widget.NewLabel("Impossible de charger les concerts")
		return container.NewStack(bg, container.NewPadded(container.NewVBox(titleLabel, errorLabel)))
	}

	concertsList := container.NewVBox()

	if len(relation.DatesLocations) == 0 {
		noConcerts := widget.NewLabel("Aucun concert programmé")
		concertsList.Add(noConcerts)
	} else {
		for location, dates := range relation.DatesLocations {
			locationLabel := canvas.NewText("📍 "+location, textGray)
			locationLabel.TextStyle = fyne.TextStyle{Bold: true}
			concertsList.Add(locationLabel)

			for _, date := range dates {
				dateLabel := widget.NewLabel("   📅 " + date)
				concertsList.Add(dateLabel)
			}
			concertsList.Add(widget.NewSeparator())
		}
	}

	scroll := container.NewVScroll(concertsList)
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
	bg.SetMinSize(fyne.NewSize(0, 150))

	titleLabel := canvas.NewText("📋 Informations", textGray)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = 20

	yearLabel := canvas.NewText(fmt.Sprintf("📅 Année de création: %d", artist.CreationDate), textGray)
	albumLabel := canvas.NewText(fmt.Sprintf("💿 Premier album: %s", artist.FirstAlbum), textGray)
	membersCountLabel := canvas.NewText(fmt.Sprintf("👥 Nombre de membres: %d", len(artist.Members)), textGray)

	info := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		yearLabel,
		albumLabel,
		membersCountLabel,
	)

	return container.NewStack(bg, container.NewPadded(info))
}

func makeDetailMembersCard(artist models.Artist) *fyne.Container {
	bg := canvas.NewRectangle(bgCard)

	titleLabel := canvas.NewText("👥 Membres du groupe", textGray)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = 20

	membersList := container.NewVBox()
	for _, member := range artist.Members {
		memberText := canvas.NewText("• "+member, textGray)
		membersList.Add(memberText)
	}

	content := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		membersList,
	)

	return container.NewStack(bg, container.NewPadded(content))
}

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
