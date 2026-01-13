package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func RenderArtistList(artists interface{}, w *AppWindow) *fyne.Container {
	title := widget.NewLabel("Liste des Artistes")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(280, 350)))

	for i := 0; i < 20; i++ {
		card := makeArtistCard(i, w)
		grid.Add(card)
	}

	scroll := container.NewVScroll(grid)

	header := container.NewVBox(title, widget.NewSeparator())

	return container.NewBorder(header, nil, nil, nil, scroll)
}

func makeArtistCard(index int, w *AppWindow) *fyne.Container {
	img := canvas.NewRectangle(color.RGBA{R: 100, G: 150, B: 200, A: 255})
	img.SetMinSize(fyne.NewSize(260, 260))

	name := widget.NewLabel(fmt.Sprintf("Artiste %d", index+1))
	name.TextStyle = fyne.TextStyle{Bold: true}
	name.Alignment = fyne.TextAlignCenter

	info := widget.NewLabel("Membres: 4 | Année: 2000")
	info.Alignment = fyne.TextAlignCenter

	btn := widget.NewButton("Voir détails", func() {
		w.ShowArtistDetail(index)
	})

	card := container.NewVBox(img, name, info, btn)

	return container.NewPadded(card)
}
