package data

import (
	"encoding/json"
	"os"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/model"
)

func LoadScooters(path string) ([]*model.Scooter, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var scooters []*model.Scooter
	if err := json.NewDecoder(file).Decode(&scooters); err != nil {
		return nil, err
	}

	return scooters, nil
}
