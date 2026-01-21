package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"io/ioutil"
	"net/http"
)

const APIBaseURL = "https://groupietrackers.herokuapp.com/api"

// FetchArtists récupère la liste complète des artistes
func FetchArtists() ([]models.Artist, error) {
	url := APIBaseURL + "/artists"

	// On fait la requête HTTP
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête API: %v", err)
	}
	defer resp.Body.Close()

	// Vérification du code de réponse
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur HTTP: %d", resp.StatusCode)
	}

	// Lecture de la réponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse: %v", err)
	}

	// Conversion JSON → struct Go
	var artists []models.Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du parsing JSON: %v", err)
	}

	return artists, nil
}
