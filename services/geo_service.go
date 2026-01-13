package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Coordinates struct {
	Lat float64
	Lon float64
}

func GeocodeAddress(address string) (Coordinates, error) {
	apiURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", address)
	params.Set("format", "json")
	params.Set("limit", "1")

	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return Coordinates{}, err
	}
	req.Header.Set("User-Agent", "groupie-tracker-app")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Coordinates{}, err
	}
	defer resp.Body.Close()

	var result []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Coordinates{}, err
	}

	if len(result) == 0 {
		return Coordinates{}, errors.New("address not found")
	}

	lat, _ := strconv.ParseFloat(result[0].Lat, 64)
	lon, _ := strconv.ParseFloat(result[0].Lon, 64)

	return Coordinates{Lat: lat, Lon: lon}, nil
}
