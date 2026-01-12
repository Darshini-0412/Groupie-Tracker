package main

import (
    "fmt"
    "groupie-tracker/ui" // ton module pour l’UI
)

func main() {
    fmt.Println("Hello Groupie Tracker!")

    window := ui.SetupWindow()
    window.ShowArtistList(nil)
    window.Run()
}
