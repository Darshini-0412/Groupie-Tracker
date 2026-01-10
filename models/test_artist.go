package models
import (
	"encoding/json"
	"testing"
)
func TestArtistParsing(t *testing.T) {
    jsonData := `{

		"id": 1,
		"name": "Queen",
		"image": "https://example.com/queen.jpg",
		"members": ["Freddie Mercury", "Brian May"],
		"creationDate": 1970,
		"firstAlbum": "1973-07-13",
		"locations": "Lhttps://api/locations/1",
		"concertDates": "https://api/dates/1",
		"relations": "https://api/dates/1"
	}`
	var a Artist
	err := json.Unmarshal([]byte(jsonData), &a)
	if err != nil {
		t.Errorf("Failed to parse JSON: %v", err)
	}
	if a.Name != "Queen" {
		t.Errorf("Expected name 'Queen', got '%s'", a.Name)
	}