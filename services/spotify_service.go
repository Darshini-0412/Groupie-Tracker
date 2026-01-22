package services

import (
	"fmt"
	"net/url"
	"strings"
)

// GenerateSpotifySearchURL crée un lien de recherche Spotify pour un artiste
func GenerateSpotifySearchURL(artistName string) string {

	cleanName := strings.TrimSpace(artistName)

	encoded := url.QueryEscape(cleanName)

	return fmt.Sprintf("https://open.spotify.com/search/%s", encoded)
}
