package service

import (
	"log"
	"net/http"
	"sync"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/model"
)

func BuildGraph(scooters []*model.Scooter) *model.Graph {
	nodes := make(map[int]*model.Scooter)
	for _, s := range scooters {
		nodes[s.ID] = s
	}

	edges := make(map[int]map[int]*model.Edge)
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, 20) // 20 seems fine

	for _, from := range scooters {
		wg.Add(1)
		go func(from *model.Scooter) {
			defer wg.Done()
			sem <- struct{}{}        //  Acquire slot
			defer func() { <-sem }() //  Release slot

			localEdges := make(map[int]*model.Edge)

			for _, to := range scooters {
				if from.ID == to.ID {
					continue
				}

				distance, duration, err := GetOSRMDistance(http.DefaultClient, from.Longitude, from.Latitude, to.Longitude, to.Latitude)
				if err != nil {
					log.Printf("Error fetching route from %d to %d: %v", from.ID, to.ID, err)
					continue
				}

				elevationDiff := to.Elevation - from.Elevation
				maxRangeKm := 35.0 * (float64(from.Battery) / 100.0)
				elevationPenalty := elevationDiff / 10.0
				effectiveRange := maxRangeKm - elevationPenalty

				if distance/1000.0 > effectiveRange {
					continue
				}

				localEdges[to.ID] = &model.Edge{
					FromID:        from.ID,
					ToID:          to.ID,
					DistanceM:     distance,
					DurationSec:   duration,
					ElevationDiff: elevationDiff,
				}
			}

			mu.Lock()
			edges[from.ID] = localEdges
			mu.Unlock()
		}(from)
	}

	wg.Wait()

	return &model.Graph{
		Nodes: nodes,
		Edges: edges,
	}
}
