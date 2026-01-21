package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"io/ioutil"
	"net/http"
)

type DatesResponse struct {
	Index []models.Date `json:"index"`
}

func GetDates() ([]models.Date, error) {
	url := "http://groupietrackers.herokuapp.com/api/dates"
	return Fetch[[]models.Date](url)
}

func GetArtistConcertDates(artistID int) ([]string, error) {
	rel, err := FetchRelationByID(artistID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération relation: %v", err)
	}
	datesMap := make(map[string]bool)
	for _, dates := range rel.DatesLocations {
		for _, date := range dates {
			if date != "" {
				datesMap[date] = true
			}
		}
	}

	var allDates []string
	for date := range datesMap {
		allDates = append(allDates, date)
	}

	return allDates, nil
}

func GetDatesByID(id int) (*models.Date, error) {
	url := fmt.Sprintf("%s/dates/%d", APIBaseURL, id)

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

	var date models.Date
	err = json.Unmarshal(body, &date)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing: %v", err)
	}

	return &date, nil
}
