package ui

import (
	"fmt"
	"groupie-tracker/models"
	"image/color"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	bgCard   = color.RGBA{R: 10, G: 10, B: 40, A: 255}
	blueSpot = color.RGBA{R: 30, G: 144, B: 255, A: 255}
	textGray = color.RGBA{R: 220, G: 220, B: 220, A: 255}
)

func RenderArtistList(artists []models.Artist, w *AppWindow) *fyne.Container {
	title := canvas.NewText("GROUPIE TRACKER", textGray)
	title.TextSize = 32
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	searchBar := CreateSmartSearchBar(w)

	header := container.NewVBox(
		container.NewCenter(title),
		container.NewPadded(searchBar),
		widget.NewSeparator(),
	)

	if len(artists) == 0 {
		emptyMsg := widget.NewLabel("Aucun artiste trouvé...")
		emptyMsg.Alignment = fyne.TextAlignCenter
		return container.NewBorder(header, nil, nil, nil, container.NewCenter(emptyMsg))
	}

	var content fyne.CanvasObject

	if len(artists) == len(w.AllArtists) {
		popularSection := createSection("🔥 Artistes les plus écoutés", getRange(artists, 0, 10), w)
		recentSection := createSection("🕐 Récemment écoutés", getRange(artists, 10, 15), w)
		suggestionsSection := createSection("💡 Suggestions", getRange(artists, 15, 20), w)
		allSection := createSection("📋 Tous les artistes", artists, w)

		mainContent := container.NewVBox(
			popularSection,
			widget.NewSeparator(),
			recentSection,
			widget.NewSeparator(),
			suggestionsSection,
			widget.NewSeparator(),
			allSection,
		)
		content = container.NewVScroll(mainContent)
	} else {
		resultSection := createSection(fmt.Sprintf("🔍 Résultats (%d)", len(artists)), artists, w)
		content = container.NewVScroll(resultSection)
	}

	filtersPanel := RenderFiltersPanel(w.AllArtists, w)
	separator := canvas.NewRectangle(color.RGBA{R: 100, G: 100, B: 100, A: 255})
	separator.SetMinSize(fyne.NewSize(2, 0))

	return container.NewBorder(header, nil, container.NewHBox(filtersPanel, separator), nil, content)
}

func getRange(artists []models.Artist, start, end int) []models.Artist {
	if start >= len(artists) {
		return []models.Artist{}
	}
	if end > len(artists) {
		end = len(artists)
	}
	return artists[start:end]
}

func createSection(title string, artists []models.Artist, w *AppWindow) *fyne.Container {
	titleText := canvas.NewText(title, textGray)
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 20

	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(250, 380)))

	for _, artist := range artists {
		grid.Add(makeArtistCard(artist, w))
	}

	return container.NewVBox(container.NewPadded(titleText), grid)
}

func makeArtistCard(artist models.Artist, w *AppWindow) *fyne.Container {
	img := loadImage(artist.Image)
	img.SetMinSize(fyne.NewSize(230, 230))

	name := widget.NewLabel(artist.Name)
	name.TextStyle = fyne.TextStyle{Bold: true}
	name.Alignment = fyne.TextAlignCenter
	name.Wrapping = fyne.TextWrapWord

	info := widget.NewLabel(fmt.Sprintf("👥 %d membres | 📅 %d", len(artist.Members), artist.CreationDate))
	info.Alignment = fyne.TextAlignCenter

	btn := widget.NewButton("Voir détails", func() {
		w.ShowArtistDetail(artist.Name)
	})

	return container.NewVBox(img, name, info, btn)
}

func loadImage(url string) *canvas.Image {
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
