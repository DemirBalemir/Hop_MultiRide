package main

//docker run -t -i -p 5000:5000 -v "${PWD}:/data" ghcr.io/project-osrm/osrm-backend osrm-routed --algorithm mld /data/turkey-latest.osrm

import (
	"fmt"
	"log"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/algorithm"
	"github.com/demirbalemir/hop/Hop_MultiRide/internal/data"
	"github.com/demirbalemir/hop/Hop_MultiRide/internal/service"
)

func main() {
	// Step 1: Load scooters from JSON
	scooters, err := data.LoadScooters("internal/data/scooters.json")
	if err != nil {
		log.Fatalf("Failed to load scooters: %v", err)
	}

	// Step 3: Build graph (distance + duration) using OSRM
	graph := service.BuildGraph(scooters)
	startNode := 17
	targetLat := 39.9479
	targetLon := 33.0440

	result := algorithm.FindOptimalPath(graph, startNode, targetLat, targetLon)

	fmt.Println("Path:", result.Path)
	fmt.Printf("Total time: %.2f seconds\n", result.TimeSoFar)
	fmt.Printf("Number of scooter switches: %d\n", result.SwitchCount)

}
