package service

import (
	"log"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/model"
)

func BuildGraph(scooters []*model.Scooter) *model.Graph {
	edges := make(map[int]map[int]*model.Edge)

	for i, from := range scooters {
		edges[i] = make(map[int]*model.Edge)

		for j, to := range scooters {
			if i == j {
				continue
			}

			distance, duration, err := GetOSRMDistance(from.Longitude, from.Latitude, to.Longitude, to.Latitude)
			if err != nil {
				log.Printf("Error fetching route from %d to %d: %v", i, j, err)
				continue
			}

			elevationDiff := to.Elevation - from.Elevation

			// Calculate scooter's current effective range
			maxRangeKm := 35.0 * (float64(from.Battery) / 100.0)
			elevationPenalty := elevationDiff / 10.0 // lose 1km per 10m gain
			effectiveRange := maxRangeKm - elevationPenalty

			// Only connect if within effective range
			if distance/1000.0 > effectiveRange {
				continue
			}

			edges[i][j] = &model.Edge{
				FromID:        i,
				ToID:          j,
				DistanceM:     distance,
				DurationSec:   duration,
				ElevationDiff: elevationDiff,
			}
		}
	}

	return &model.Graph{
		Nodes: scooters,
		Edges: edges,
	}
}
