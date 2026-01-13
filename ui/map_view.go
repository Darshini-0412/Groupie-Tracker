package ui

import (
	"fmt"
	"groupie-tracker/services"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func RenderMap(locations []string) *fyne.Container {
	title := widget.NewLabel("Carte des concerts")
	title.TextStyle = fyne.TextStyle{Bold: true}

	mapArea := canvas.NewRectangle(color.RGBA{R: 200, G: 220, B: 240, A: 255})
	mapArea.SetMinSize(fyne.NewSize(700, 500))

	list := container.NewVBox()

	for _, loc := range locations {
		coords, err := services.GeocodeAddress(loc)
		if err == nil {
			marker := widget.NewLabel(fmt.Sprintf("📍 %s (%.2f, %.2f)", loc, coords.Lat, coords.Lon))
			list.Add(marker)
		}
	}

	scroll := container.NewVScroll(list)

	return container.NewBorder(
		title,
		nil, nil, nil,
		container.NewHBox(mapArea, scroll),
	)
}
