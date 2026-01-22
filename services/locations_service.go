package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"io/ioutil"
	"net/http"
	"strings"
)

type LocationResponse struct {
	Index []models.Location `json:"index"`
}

func GetLocations() ([]models.Location, error) {
	url := "https://groupietrackers.herokuapp.com/api/locations"
	return Fetch[[]models.Location](url)
}

func GetArtistLocations(artistID int) ([]string, error) {
	rel, err := FetchRelationByID(artistID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération relation: %v", err)
	}

	locationsMap := make(map[string]bool)
	for location := range rel.DatesLocations {
		location = strings.TrimSpace(location)
		if location != "" {
			locationsMap[location] = true
		}
	}

	var locations []string
	for loc := range locationsMap {
		locations = append(locations, loc)
	}

	return locations, nil
}

func GetLocationsByID(id int) (*models.Location, error) {
	url := fmt.Sprintf("%s/locations/%d", APIBaseURL, id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur requête: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur HTTP: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture: %v", err)
	}

	var location models.Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing: %v", err)
	}

	return &location, nil
}
