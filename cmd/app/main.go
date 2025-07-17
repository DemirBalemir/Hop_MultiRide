package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/algorithm"
	"github.com/demirbalemir/hop/Hop_MultiRide/internal/data"
	"github.com/demirbalemir/hop/Hop_MultiRide/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// Step 0: Load .env and get API key
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		log.Fatalf("GOOGLE_MAPS_API_KEY not set")
	}

	//FLAGS
	genScooters := flag.Bool("generate", false, "Generate scooter data")
	addElevation := flag.Bool("elevation", false, "Add elevation to scooters.json")
	flag.Parse()

	if *genScooters {
		err := data.GenerateScooters(0, 20, "internal/data/scooters.json")
		if err != nil {
			log.Fatalf("Failed to generate scooters: %v", err)
		}
		fmt.Println("Scooters generated.")
		return
	}

	if *addElevation {
		err := service.AddElevationToScooters("internal/data/scooters.json", apiKey)
		if err != nil {
			log.Fatalf("Failed to add elevation: %v", err)
		}
		fmt.Println("Elevation updated.")
		return
	}

	// Step 2: Load scooters from JSON
	scooters, err := data.LoadScooters("internal/data/scooters.json")
	if err != nil {
		log.Fatalf("Failed to load scooters: %v", err)
	}

	// Step 3: Build graph using OSRM
	graph := service.BuildGraph(scooters)

	// Example routing inputs
	startNode := 3
	targetLat := 39.9479
	targetLon := 33.0440
	if _, ok := graph.Nodes[startNode]; !ok {
		log.Fatalf("Start node %d not found — invalid start node", startNode)
	}

	// Step 4: Run optimal pathfinding
	result := algorithm.FindOptimalPath(graph, startNode, targetLat, targetLon)

	// Step 5: Print results
	fmt.Println("Path:", result.Path)
	fmt.Printf("Total time: %.2f seconds\n", result.TimeSoFar)
	fmt.Printf("Number of scooter switches: %d\n", result.SwitchCount)
}
