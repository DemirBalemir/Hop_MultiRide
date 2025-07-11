package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/model"
)

var Get_Elev = GetElevation

func AddElevationToScooters(filepath string, apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("GOOGLE_MAPS_API_KEY not set")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var scooters []model.Scooter
	if err := json.Unmarshal(data, &scooters); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	for i, scooter := range scooters {
		elevation, err := Get_Elev(scooter.Latitude, scooter.Longitude, apiKey) // <-- ✅ USE MOCKABLE VERSION
		if err != nil {
			log.Printf("failed to get elevation for scooter %d: %v", scooter.ID, err)
			continue
		}
		scooters[i].Elevation = elevation
	}

	updated, err := json.MarshalIndent(scooters, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated scooters: %w", err)
	}

	if err := os.WriteFile(filepath, updated, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func GetElevation(lat, lon float64, apiKey string) (float64, error) {
	url := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/elevation/json?locations=%f,%f&key=%s",
		lat, lon, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	var res struct {
		Results []struct {
			Elevation float64 `json:"elevation"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(res.Results) == 0 {
		return 0, fmt.Errorf("no elevation results")
	}

	return res.Results[0].Elevation, nil
}
