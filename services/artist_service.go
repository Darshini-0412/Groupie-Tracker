package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"io/ioutil"
	"net/http"
)

// FetchArtistByID récupère un artiste spécifique par son ID
func FetchArtistByID(id int) (*models.Artist, error) {
	url := fmt.Sprintf("%s/artists/%d", APIBaseURL, id)

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
		return nil, fmt.Errorf("erreur lors de la lecture: %v", err)
	}

	var artist models.Artist
	err = json.Unmarshal(body, &artist)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing JSON: %v", err)
	}

	return &artist, nil
}
