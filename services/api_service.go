package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"io/ioutil"
	"net/http"
)

const APIBaseURL = "https://groupietrackers.herokuapp.com/api"

// FetchArtists récupère tous les artistes
func FetchArtists() ([]models.Artist, error) {
	url := APIBaseURL + "/artists"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur HTTP: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse: %v", err)
	}

	var artists []models.Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du parsing JSON: %v", err)
	}

	return artists, nil
}
