package ui

import (
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	bgCard   = color.RGBA{R: 10, G: 10, B: 40, A: 255}
	textGray = color.RGBA{R: 220, G: 220, B: 220, A: 255}
)

// Affiche la page principale avec tous les artistes
func RenderArtistList(artists []models.Artist, w *AppWindow) *fyne.Container {
	title := canvas.NewText("GROUPIE TRACKER", nil)
	title.TextSize = 32
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	searchBar := CreateSmartSearchBar(w)

	// Bouton d'aide pour afficher les raccourcis
	helpBtn := widget.NewButton("❓ Raccourcis", func() {
		showShortcutsHelp(w)
	})
	helpBtn.Importance = widget.LowImportance

	header := container.NewVBox(
		container.NewCenter(title),
		container.NewBorder(nil, nil, nil, helpBtn, container.NewPadded(searchBar)),
		widget.NewSeparator(),
	)

	if len(artists) == 0 {
		emptyMsg := widget.NewLabel("Aucun artiste trouvé...")
		emptyMsg.Alignment = fyne.TextAlignCenter
		return container.NewBorder(header, nil, nil, nil, container.NewCenter(emptyMsg))
	}

	var scrollContent *fyne.Container

	if len(artists) == len(w.AllArtists) {
		sections := []fyne.CanvasObject{}

		favoriteArtists := w.Favorites.GetFavorites(w.AllArtists)
		if len(favoriteArtists) > 0 {
			favSection := createSection("Ma sélection", favoriteArtists, w)
			sections = append(sections, favSection, widget.NewSeparator())
		}

		popularSection := createSection("Artistes les plus écoutés", getArtistsRange(artists, 0, 10), w)
		sections = append(sections, popularSection, widget.NewSeparator())

		recentSection := createSection("Récemment écoutés", getArtistsRange(artists, 10, 15), w)
		sections = append(sections, recentSection, widget.NewSeparator())

		suggestionsSection := createSection("Suggestions", getArtistsRange(artists, 15, 20), w)
		sections = append(sections, suggestionsSection, widget.NewSeparator())

		allSection := createSection("Tous les artistes", artists, w)
		sections = append(sections, allSection)

		scrollContent = container.NewVBox(sections...)
	} else {
		resultSection := createSection(fmt.Sprintf("Résultats (%d)", len(artists)), artists, w)
		scrollContent = resultSection
	}

	content := container.NewVScroll(scrollContent)

	filtersPanel := CreateFiltersPanel(w.AllArtists, w)
	separator := canvas.NewRectangle(color.RGBA{R: 100, G: 100, B: 100, A: 255})
	separator.SetMinSize(fyne.NewSize(2, 0))

	return container.NewBorder(header, nil, container.NewHBox(filtersPanel, separator), nil, content)
}

// Affiche une fenêtre popup avec les raccourcis clavier
func showShortcutsHelp(w *AppWindow) {
	helpText := `⌨️ RACCOURCIS CLAVIER DISPONIBLES

 Ctrl+H : Retour à l'accueil
  Echap : Retour en arrière
 Ctrl+R : Réinitialiser les filtres
 Ctrl+Q : Quitter l'application

💡 Astuce : Utilisez ces raccourcis pour naviguer plus rapidement !`

	dialog.ShowInformation("Raccourcis clavier", helpText, w.Window)
}

func getArtistsRange(artists []models.Artist, start, end int) []models.Artist {
	if start >= len(artists) {
		return []models.Artist{}
	}
	if end > len(artists) {
		end = len(artists)
	}
	return artists[start:end]
}

func createSection(sectionTitle string, artists []models.Artist, w *AppWindow) *fyne.Container {
	titleText := canvas.NewText(sectionTitle, nil)
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 20

	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(250, 400)))

	for _, artist := range artists {
		card := makeRealArtistCard(artist, w)
		grid.Add(card)
	}

	return container.NewVBox(
		container.NewPadded(titleText),
		grid,
	)
}

func makeRealArtistCard(artist models.Artist, w *AppWindow) *fyne.Container {
	img := loadImageFromURL(artist.Image)
	img.SetMinSize(fyne.NewSize(230, 230))

	name := widget.NewLabel(artist.Name)
	name.TextStyle = fyne.TextStyle{Bold: true}
	name.Alignment = fyne.TextAlignCenter
	name.Wrapping = fyne.TextWrapWord

	membersCount := len(artist.Members)
	info := widget.NewLabel(fmt.Sprintf("%d membres | %d", membersCount, artist.CreationDate))
	info.Alignment = fyne.TextAlignCenter

	favoriteIcon := "🤍"
	if w.Favorites.IsFavorite(artist.ID) {
		favoriteIcon = "❤️"
	}

	favoriteBtn := widget.NewButton(favoriteIcon, func() {
		w.Favorites.Toggle(artist.ID)
		w.ShowArtistList()
	})

	detailBtn := widget.NewButton("Voir détails", func() {
		w.ShowArtistDetail(artist.Name)
	})

	spotifyBtn := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
		spotifyURL := services.GenerateSpotifySearchURL(artist.Name)
		openURL(spotifyURL)
	})
	spotifyBtn.Importance = widget.MediumImportance

	buttons := container.NewGridWithColumns(3,
		favoriteBtn,
		detailBtn,
		spotifyBtn,
	)

	card := container.NewVBox(
		img,
		name,
		info,
		buttons,
	)

	return card
}
