package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func RenderArtistDetail(artist interface{}, w fyne.Window) *fyne.Container {
	backBtn := widget.NewButton("← Retour", func() {})

	title := widget.NewLabel("Détails de l'artiste")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	img := canvas.NewRectangle(color.RGBA{R: 150, G: 150, B: 150, A: 255})
	img.SetMinSize(fyne.NewSize(350, 350))

	info := makeInfoSection()
	concerts := makeConcertSection()

	top := container.NewHBox(img, info)
	scroll := container.NewVScroll(container.NewVBox(top, widget.NewSeparator(), concerts))

	return container.NewBorder(container.NewVBox(backBtn, title, widget.NewSeparator()), nil, nil, nil, scroll)
}

func makeInfoSection() *fyne.Container {
	name := widget.NewLabel("Queen")
	name.TextStyle = fyne.TextStyle{Bold: true}

	creation := widget.NewLabel("Création: 1970")
	album := widget.NewLabel("Premier album: 1973")

	membersTitle := widget.NewLabel("Membres:")
	membersTitle.TextStyle = fyne.TextStyle{Bold: true}

	members := widget.NewLabel("Freddie Mercury\nBrian May\nRoger Taylor\nJohn Deacon")

	return container.NewVBox(name, widget.NewSeparator(), creation, album, widget.NewSeparator(), membersTitle, members)
}

func makeConcertSection() *fyne.Container {
	title := widget.NewLabel("Dates et lieux des concerts")
	title.TextStyle = fyne.TextStyle{Bold: true}

	list := container.NewVBox()

	locations := []string{"Paris, France", "London, UK", "New York, USA", "Tokyo, Japan"}
	dates := []string{"15-01-2024", "20-02-2024", "10-03-2024", "05-04-2024"}

	for i := 0; i < len(locations); i++ {
		list.Add(makeConcertCard(locations[i], dates[i]))
	}

	return container.NewVBox(title, widget.NewSeparator(), list)
}

func makeConcertCard(location string, date string) *fyne.Container {
	loc := widget.NewLabel(fmt.Sprintf("📍 %s", location))
	loc.TextStyle = fyne.TextStyle{Bold: true}

	dt := widget.NewLabel(fmt.Sprintf("📅 %s", date))

	return container.NewVBox(loc, dt, widget.NewSeparator())
}
