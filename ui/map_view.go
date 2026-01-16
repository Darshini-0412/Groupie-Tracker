package ui

import (
	"image/color"

	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Convertit latitude/longitude en position X/Y sur la carte
func convertLatLonToXY(lat, lon float64, width, height float32) (float32, float32) {
	x := float32((lon+180.0)/360.0) * width
	y := float32((90.0-lat)/180.0) * height
	return x, y
}

func RenderMap(past []string, future []string) *fyne.Container {
	title := widget.NewLabel("Carte des concerts")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// --- Image de carte ---
	mapImg := canvas.NewImageFromFile("assets/world_map.png")
	mapWidth := float32(700)
	mapHeight := float32(500)
	mapImg.SetMinSize(fyne.NewSize(mapWidth, mapHeight))
	mapImg.FillMode = canvas.ImageFillContain

	// --- Conteneur pour les marqueurs ---
	var markers []fyne.CanvasObject

	// --- Fonction interne pour ajouter un point ---
	addMarker := func(loc string, col color.NRGBA) {
		coords, err := services.GeocodeAddress(loc)
		if err != nil {
			return
		}

		x, y := convertLatLonToXY(coords.Lat, coords.Lon, mapWidth, mapHeight)

		dot := canvas.NewCircle(col)
		dot.Resize(fyne.NewSize(10, 10))
		dot.Move(fyne.NewPos(x, y))

		markers = append(markers, dot)
	}

	// 🔵 Concerts passés
	for _, loc := range past {
		addMarker(loc, color.NRGBA{R: 0, G: 100, B: 255, A: 255})
	}

	// 🔴 Concerts futurs
	for _, loc := range future {
		addMarker(loc, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	}

	// --- Superposition carte + marqueurs ---
	mapLayer := []fyne.CanvasObject{mapImg}
	mapLayer = append(mapLayer, markers...)

	mapContainer := container.NewWithoutLayout(mapLayer...)

	// --- Liste des lieux à droite ---
	list := container.NewVBox()

	list.Add(widget.NewLabel("Concerts passés :"))
	for _, loc := range past {
		list.Add(widget.NewLabel("🔵 " + loc))
	}

	list.Add(widget.NewLabel("\nConcerts futurs :"))
	for _, loc := range future {
		list.Add(widget.NewLabel("🔴 " + loc))
	}

	scroll := container.NewVScroll(list)

	// --- Layout final ---
	return container.NewBorder(
		title,
		nil, nil, nil,
		container.NewHBox(mapContainer, scroll),
	)
}
