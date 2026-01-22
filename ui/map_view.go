package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"math"
	"net/http"
	"strings"

	"groupie-tracker/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// RenderMap affiche la carte avec les lieux des concerts
func RenderMap(past []string, future []string, selected *string) *fyne.Container {

	countryLabel := widget.NewLabel("🌍 Aucun pays sélectionné")
	countryLabel.Alignment = fyne.TextAlignCenter
	countryLabel.TextStyle = fyne.TextStyle{Bold: true}

	mapContainer := container.NewMax()
	mapContainer.Add(tryLoadMap(selected, future, past))

	// la liste des lieux à gauche
	list := widget.NewList(
		func() int {
			return len(future) + len(past)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("📍"),
				widget.NewLabel("Location"),
				widget.NewLabel("→"),
				widget.NewLabel("Lat, Lon"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			var location string
			var isPast bool

			if id < len(future) {
				location = future[id]
				isPast = false
			} else {
				location = past[id-len(future)]
				isPast = true
			}

			cleanLoc := cleanAddress(location)
			coords, err := services.GeocodeAddress(cleanLoc)
			coordsText := "Chargement..."
			if err == nil {
				coordsText = fmt.Sprintf("%.4f, %.4f", coords.Lat, coords.Lon)
			} else {
				coordsText = "Non trouvé"
			}

			icon := "🔴"
			if isPast {
				icon = "⚫"
			}
			if selected != nil && *selected == cleanLoc {
				icon = "🔵"
			}

			box := obj.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(icon)
			box.Objects[1].(*widget.Label).SetText(location)
			box.Objects[3].(*widget.Label).SetText(coordsText)
		},
	)

	// quand on clique sur un lieu
	list.OnSelected = func(id widget.ListItemID) {
		var location string

		if id < len(future) {
			location = future[id]
		} else {
			location = past[id-len(future)]
		}

		clean := cleanAddress(location)
		*selected = clean

		coords, err := services.GeocodeAddress(clean)
		if err == nil {
			if coords.Country != "" {
				countryLabel.SetText("🌍 Pays : " + coords.Country)
			} else {
				countryLabel.SetText("🌍 Lieu : " + coords.City)
			}
		} else {
			countryLabel.SetText("🌍 Erreur pour: " + clean)
		}

		// recharger la carte
		newMap := tryLoadMap(selected, future, past)
		mapContainer.RemoveAll()
		mapContainer.Add(newMap)
		mapContainer.Refresh()
	}

	info := widget.NewLabel("🗺️ Géolocalisation des concerts")
	info.TextStyle = fyne.TextStyle{Bold: true}

	legend := container.NewVBox(
		widget.NewLabel("🔴 Concerts à venir: "+fmt.Sprintf("%d", len(future))),
		widget.NewLabel("⚫ Concerts passés: "+fmt.Sprintf("%d", len(past))),
		widget.NewLabel("📍 Marqueur rouge = lieu sélectionné"),
	)

	return container.NewBorder(
		container.NewVBox(info, widget.NewSeparator(), countryLabel),
		container.NewVBox(widget.NewSeparator(), legend),
		nil,
		nil,
		container.NewVSplit(
			mapContainer,
			list,
		),
	)
}

func cleanAddress(address string) string {
	clean := strings.ReplaceAll(address, "-", " ")
	clean = strings.ReplaceAll(clean, "_", " ")
	clean = strings.TrimSpace(clean)
	clean = strings.Join(strings.Fields(clean), " ")
	return clean
}

// 🔥 AJOUT — conversion lat/lon → pixel dans l’image
func latLonToPixel(lat, lon float64, zoom int, centerTileX, centerTileY int) (int, int) {
	tileX, tileY := latLonToTile(lat, lon, zoom)

	dx := tileX - centerTileX
	dy := tileY - centerTileY

	px := (dx+1)*256 + 128
	py := (dy+1)*256 + 128

	return px, py
}

// tryLoadMap charge la carte avec un marqueur
func tryLoadMap(selected *string, future, past []string) fyne.CanvasObject {
	centerLat, centerLon, zoom := 20.0, 0.0, 2

	var selectedCoords *services.Coordinates

	// si un lieu est sélectionné, on centre dessus
	if selected != nil && *selected != "" {
		coords, err := services.GeocodeAddress(*selected)
		if err == nil {
			centerLat, centerLon, zoom = coords.Lat, coords.Lon, 8
			selectedCoords = &coords
		}
	}

	// calculer quelle tuile au centre
	centerTileX, centerTileY := latLonToTile(centerLat, centerLon, zoom)

	mapWidth := 768
	mapHeight := 768

	mapImg := image.NewRGBA(image.Rect(0, 0, mapWidth, mapHeight))
	draw.Draw(mapImg, mapImg.Bounds(), &image.Uniform{color.RGBA{200, 220, 240, 255}}, image.Point{}, draw.Src)

	tilesLoaded := 0

	// télécharger 9 tuiles (3x3)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			tx := centerTileX + dx
			ty := centerTileY + dy

			url := fmt.Sprintf("https://tile.openstreetmap.org/%d/%d/%d.png", zoom, tx, ty)

			tile := downloadImage(url)
			if tile != nil {
				destX := (dx + 1) * 256
				destY := (dy + 1) * 256
				destRect := image.Rect(destX, destY, destX+256, destY+256)

				draw.Draw(mapImg, destRect, tile, image.Point{}, draw.Src)
				tilesLoaded++
			}
		}
	}

	// si rien a chargé
	if tilesLoaded == 0 {
		placeholder := canvas.NewRectangle(color.RGBA{230, 230, 230, 255})
		placeholder.SetMinSize(fyne.NewSize(700, 700))

		errorMsg := "🗺️ Carte pas dispo\n(vérifier internet)"
		if selected != nil && *selected != "" {
			errorMsg = fmt.Sprintf("🗺️ Marche pas pour:\n%s", *selected)
		}

		label := widget.NewLabel(errorMsg)
		label.Alignment = fyne.TextAlignCenter

		return container.NewStack(placeholder, container.NewCenter(label))
	}

	// 🔥 AJOUT — dessiner tous les concerts
	for _, loc := range append(future, past...) {
		clean := cleanAddress(loc)
		coords, err := services.GeocodeAddress(clean)
		if err != nil {
			continue
		}

		px, py := latLonToPixel(coords.Lat, coords.Lon, zoom, centerTileX, centerTileY)
		drawBigMarker(mapImg, px, py)
	}

	// 🔵 marqueur du lieu sélectionné (au centre)
	if selectedCoords != nil {
		markerX := mapWidth / 2
		markerY := mapHeight / 2
		drawBigMarker(mapImg, markerX, markerY)
	}

	canvasImg := canvas.NewImageFromImage(mapImg)
	canvasImg.FillMode = canvas.ImageFillContain
	canvasImg.SetMinSize(fyne.NewSize(700, 700))

	var infoText string
	if selectedCoords != nil {
		infoText = fmt.Sprintf("📍 %s, %s\n📊 Lat: %.4f, Lon: %.4f | 🔍 Zoom: %d",
			selectedCoords.City, selectedCoords.Country, selectedCoords.Lat, selectedCoords.Lon, zoom)
	} else {
		infoText = fmt.Sprintf("🌍 Vue mondiale | 🔍 Zoom: %d", zoom)
	}

	infoLabel := widget.NewLabel(infoText)
	infoLabel.Alignment = fyne.TextAlignCenter
	infoLabel.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewBorder(
		container.NewPadded(infoLabel),
		nil, nil, nil,
		container.NewCenter(canvasImg),
	)
}

// drawBigMarker dessine un gros point rouge
func drawBigMarker(img *image.RGBA, x, y int) {
	markerColor := color.RGBA{255, 0, 0, 255}
	outlineColor := color.RGBA{255, 255, 255, 255}
	shadowColor := color.RGBA{0, 0, 0, 100}
	centerColor := color.RGBA{255, 255, 255, 255}

	shadowRadius := 18
	outlineRadius := 16
	markerRadius := 14
	centerRadius := 6

	// ombre
	for dy := -shadowRadius; dy <= shadowRadius; dy++ {
		for dx := -shadowRadius; dx <= shadowRadius; dx++ {
			if dx*dx+dy*dy <= shadowRadius*shadowRadius {
				px, py := x+dx+2, y+dy+2
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, shadowColor)
				}
			}
		}
	}

	// contour blanc
	for dy := -outlineRadius; dy <= outlineRadius; dy++ {
		for dx := -outlineRadius; dx <= outlineRadius; dx++ {
			if dx*dx+dy*dy <= outlineRadius*outlineRadius {
				px, py := x+dx, y+dy
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, outlineColor)
				}
			}
		}
	}

	// cercle rouge
	for dy := -markerRadius; dy <= markerRadius; dy++ {
		for dx := -markerRadius; dx <= markerRadius; dx++ {
			if dx*dx+dy*dy <= markerRadius*markerRadius {
				px, py := x+dx, y+dy
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, markerColor)
				}
			}
		}
	}

	// point blanc au milieu
	for dy := -centerRadius; dy <= centerRadius; dy++ {
		for dx := -centerRadius; dx <= centerRadius; dx++ {
			if dx*dx+dy*dy <= centerRadius*centerRadius {
				px, py := x+dx, y+dy
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, centerColor)
				}
			}
		}
	}

	// croix noire
	crossSize := 10
	crossThickness := 2
	crossColor := color.RGBA{0, 0, 0, 255}

	for dx := -crossSize; dx <= crossSize; dx++ {
		for t := -crossThickness; t <= crossThickness; t++ {
			px, py := x+dx, y+t
			if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
				img.Set(px, py, crossColor)
			}
		}
	}

	for dy := -crossSize; dy <= crossSize; dy++ {
		for t := -crossThickness; t <= crossThickness; t++ {
			px, py := x+t, y+dy
			if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
				img.Set(px, py, crossColor)
			}
		}
	}
}

func downloadImage(url string) image.Image {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "GroupieTracker/1.0")
	req.Header.Set("Referer", "https://www.openstreetmap.org/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil
	}

	return img
}

// latLonToTile transforme GPS en numéro de tuile
func latLonToTile(lat, lon float64, zoom int) (int, int) {
	latRad := lat * math.Pi / 180.0
	n := math.Pow(2.0, float64(zoom))

	x := int((lon + 180.0) / 360.0 * n)
	y := int((1.0 - math.Log(math.Tan(latRad)+1/math.Cos(latRad))/math.Pi) / 2.0 * n)

	maxTile := int(n) - 1
	if x < 0 {
		x = 0
	}
	if x > maxTile {
		x = maxTile
	}
	if y < 0 {
		y = 0
	}
	if y > maxTile {
		y = maxTile
	}

	return x, y
}
