package service

import (
	"net/http"
	"testing"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/model"
	"github.com/stretchr/testify/assert"
)

// Mock the OSRM call
// in graph_builder_test.go
func mockDistance(client *http.Client, fromLon, fromLat, toLon, toLat float64) (float64, float64, error) {
	return 1000, 300, nil // 1 km, 5 min dummy values
}

func TestBuildGraph(t *testing.T) {
	// Override the real function
	GetOSRMDistance = mockDistance

	tests := []struct {
		name          string
		scooters      []*model.Scooter
		expectedNodes int
		expectedEdges int
	}{
		{
			name: "single scooter no edges",
			scooters: []*model.Scooter{
				{ID: 1, Latitude: 0, Longitude: 0, Battery: 100, Elevation: 0},
			},
			expectedNodes: 1,
			expectedEdges: 0,
		},
		{
			name: "two scooters with connection",
			scooters: []*model.Scooter{
				{ID: 1, Latitude: 0, Longitude: 0, Battery: 100, Elevation: 0},
				{ID: 2, Latitude: 0.01, Longitude: 0.01, Battery: 100, Elevation: 0},
			},
			expectedNodes: 2,
			expectedEdges: 2, // bi-directional connection
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			graph := BuildGraph(tc.scooters)

			assert.Equal(t, tc.expectedNodes, len(graph.Nodes))
			totalEdges := 0
			for _, edges := range graph.Edges {
				totalEdges += len(edges)
			}
			assert.Equal(t, tc.expectedEdges, totalEdges)
		})
	}
}
