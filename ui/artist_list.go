package ui

import (
	"fmt"
	"groupie-tracker/models"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Couleurs globales de l'interface
var (
	bgCard   = color.RGBA{R: 10, G: 10, B: 40, A: 255}    // Bleu foncé pour les cartes
	textGray = color.RGBA{R: 220, G: 220, B: 220, A: 255} // Gris clair pour le texte
)

// RenderArtistList affiche la page principale avec la liste des artistes
func RenderArtistList(artists []models.Artist, w *AppWindow) *fyne.Container {
	// Titre principal
	title := canvas.NewText("ＧRØUPIE TЯΛCKΞR", nil)
	title.TextSize = 32
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// Barre de recherche
	searchBar := CreateSmartSearchBar(w)

	// En-tête (titre + recherche)
	header := container.NewVBox(
		container.NewCenter(title),
		container.NewPadded(searchBar),
		widget.NewSeparator(),
	)

	// Si aucun artiste trouvé
	if len(artists) == 0 {
		emptyMsg := widget.NewLabel("Aucun artiste trouvé...")
		emptyMsg.Alignment = fyne.TextAlignCenter
		return container.NewBorder(header, nil, nil, nil, container.NewCenter(emptyMsg))
	}

	var scrollContent *fyne.Container

	// Si on affiche TOUS les artistes → créer des sections
	if len(artists) == len(w.AllArtists) {
		sections := []fyne.CanvasObject{}

		// Section Favoris (si il y en a)
		favoriteArtists := w.Favorites.GetFavorites(w.AllArtists)
		if len(favoriteArtists) > 0 {
			favSection := createSection("❤️ Ma Sélection", favoriteArtists, w)
			sections = append(sections, favSection, widget.NewSeparator())
		}

		// Section artistes populaires (0-10)
		popularSection := createSection("🔥 Artistes les plus écoutés", getArtistsRange(artists, 0, 10), w)
		sections = append(sections, popularSection, widget.NewSeparator())

		// Section récents (10-15)
		recentSection := createSection("🕐 Récemment écoutés", getArtistsRange(artists, 10, 15), w)
		sections = append(sections, recentSection, widget.NewSeparator())

		// Section suggestions (15-20)
		suggestionsSection := createSection("💡 Suggestions", getArtistsRange(artists, 15, 20), w)
		sections = append(sections, suggestionsSection, widget.NewSeparator())

		// Section tous les artistes
		allSection := createSection("📋 Tous les artistes", artists, w)
		sections = append(sections, allSection)

		scrollContent = container.NewVBox(sections...)
	} else {
		// Résultats de recherche/filtre
		resultSection := createSection(fmt.Sprintf("🔍 Résultats (%d)", len(artists)), artists, w)
		scrollContent = resultSection
	}

	// Rendre scrollable
	content := container.NewVScroll(scrollContent)

	// Panneau de filtres à gauche
	filtersPanel := CreateFiltersPanel(w.AllArtists, w)
	separator := canvas.NewRectangle(color.RGBA{R: 100, G: 100, B: 100, A: 255})
	separator.SetMinSize(fyne.NewSize(2, 0))

	return container.NewBorder(header, nil, container.NewHBox(filtersPanel, separator), nil, content)
}

// getArtistsRange extrait une plage d'artistes (ex: 0-10)
func getArtistsRange(artists []models.Artist, start, end int) []models.Artist {
	if start >= len(artists) {
		return []models.Artist{}
	}
	if end > len(artists) {
		end = len(artists)
	}
	return artists[start:end]
}

// createSection crée une section avec titre + grille de cartes
func createSection(sectionTitle string, artists []models.Artist, w *AppWindow) *fyne.Container {
	titleText := canvas.NewText(sectionTitle, nil)
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 20

	// Grille responsive de cartes (250x400 chacune)
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

// makeRealArtistCard crée une carte pour un artiste
func makeRealArtistCard(artist models.Artist, w *AppWindow) *fyne.Container {
	// Image de l'artiste
	img := loadImageFromURL(artist.Image)
	img.SetMinSize(fyne.NewSize(230, 230))

	// Nom de l'artiste
	name := widget.NewLabel(artist.Name)
	name.TextStyle = fyne.TextStyle{Bold: true}
	name.Alignment = fyne.TextAlignCenter
	name.Wrapping = fyne.TextWrapWord

	// Infos (nombre membres + année)
	membersCount := len(artist.Members)
	info := widget.NewLabel(fmt.Sprintf("👥 %d membres | 📅 %d", membersCount, artist.CreationDate))
	info.Alignment = fyne.TextAlignCenter

	// Bouton favoris (🤍 ou ❤️)
	favoriteIcon := "🤍"
	if w.Favorites.IsFavorite(artist.ID) {
		favoriteIcon = "❤️"
	}

	favoriteBtn := widget.NewButton(favoriteIcon, func() {
		w.Favorites.Toggle(artist.ID)
		w.ShowArtistList() // Rafraîchir l'affichage
	})

	// Bouton détails
	detailBtn := widget.NewButton("Voir détails", func() {
		w.ShowArtistDetail(artist.Name)
	})

	// Grille 2 boutons côte à côte
	buttons := container.NewGridWithColumns(2, favoriteBtn, detailBtn)

	// Assembler la carte
	card := container.NewVBox(img, name, info, buttons)
	return card
}
