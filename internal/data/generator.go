package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/model"
)

func GenerateScooters(startID, count int, filePath string) error {

	latMin, latMax := 39.8, 40.2
	lonMin, lonMax := 32.5, 33.2
	elevationMin, elevationMax := 850.0, 1100.0

	var scooters []model.Scooter

	for i := 0; i < count; i++ {
		s := model.Scooter{
			ID:        startID + i,
			Latitude:  roundFloat(randFloat(latMin, latMax), 6),
			Longitude: roundFloat(randFloat(lonMin, lonMax), 6),
			Battery:   rand.Intn(101),
			Elevation: roundFloat(randFloat(elevationMin, elevationMax), 1),
		}
		scooters = append(scooters, s)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(scooters); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	fmt.Printf("Successfully generated %d scooters into %s\n", count, filePath)
	return nil
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func roundFloat(val float64, precision int) float64 {
	format := fmt.Sprintf("%%.%df", precision)
	str := fmt.Sprintf(format, val)
	var rounded float64
	fmt.Sscanf(str, "%f", &rounded)
	return rounded
}
