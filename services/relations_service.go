package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Relations struct {
	Index []Relation `json:"index"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func FetchRelations() (*Relations, error) {
	url := APIBaseURL + "/relation"

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

	var relations Relations
	err = json.Unmarshal(body, &relations)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing: %v", err)
	}

	return &relations, nil
}

func FetchRelationByID(id int) (*Relation, error) {
	url := fmt.Sprintf("%s/relation/%d", APIBaseURL, id)

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

	var relation Relation
	err = json.Unmarshal(body, &relation)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing: %v", err)
	}

	return &relation, nil
}

// SplitPastFutureConcerts sépare les concerts passés et futurs
func SplitPastFutureConcerts(rel Relation) (past []string, future []string) {
	now := time.Now()

	for location, dates := range rel.DatesLocations {
		for _, d := range dates {
			date, err := time.Parse("2006-01-02", d)
			if err != nil {
				continue
			}

			if date.Before(now) {
				past = append(past, location)
			} else {
				future = append(future, location)
			}
		}
	}

	return past, future
}
