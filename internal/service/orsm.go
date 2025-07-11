package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OSRMResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
	} `json:"routes"`
}

func GetOSRMDistance(fromLng, fromLat, toLng, toLat float64) (float64, float64, error) {
	url := fmt.Sprintf(
		"http://localhost:5000/route/v1/driving/%.6f,%.6f;%.6f,%.6f?overview=false",
		fromLng, fromLat, toLng, toLat,
	)

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("OSRM request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read OSRM response: %w", err)
	}

	var data OSRMResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, 0, fmt.Errorf("failed to parse OSRM JSON: %w", err)
	}

	if len(data.Routes) == 0 {
		return 0, 0, fmt.Errorf("no routes returned")
	}

	return data.Routes[0].Distance, data.Routes[0].Duration, nil
}
