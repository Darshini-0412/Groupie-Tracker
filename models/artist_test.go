package models

import (
	"encoding/json"
	"testing"
)

// Test pour vérifier qu'on parse bien le JSON de l'API
func TestArtistParsing(t *testing.T) {
	// JSON d'exemple qui ressemble à ce que renvoie l'API
	jsonData := `{
		"id": 1,
		"name": "Queen",
		"image": "https://example.com/queen.jpg",
		"members": ["Freddie Mercury", "Brian May"],
		"creationDate": 1970,
		"firstAlbum": "1973-07-13",
		"locations": "https://api/locations/1",
		"concertDates": "https://api/dates/1",
		"relations": "https://api/relations/1"
	}`

	var a Artist
	err := json.Unmarshal([]byte(jsonData), &a)
	if err != nil {
		t.Errorf("Failed to parse JSON: %v", err)
	}

	// On vérifie que le nom est bien "Queen"
	if a.Name != "Queen" {
		t.Errorf("Expected name 'Queen', got '%s'", a.Name)
	}
}
