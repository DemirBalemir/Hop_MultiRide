package algorithm

import (
	"errors"
	"testing"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/model"
)

func TestFindOptimalPath(t *testing.T) {
	tests := []struct {
		name           string
		graph          *model.Graph
		startNode      int
		targetLat      float64
		targetLon      float64
		mockDistanceFn func(float64, float64, float64, float64) (float64, float64, error)
		expectNil      bool
		expectPath     []int
		expectSwitches int
	}{
		{
			name: "Directly reachable from start",
			graph: &model.Graph{
				Nodes: []*model.Scooter{
					{ID: 0, Latitude: 1.0, Longitude: 1.0, Battery: 100},
				},
				Edges: map[int]map[int]*model.Edge{},
			},
			startNode: 0,
			targetLat: 1.1,
			targetLon: 1.1,
			mockDistanceFn: func(fromLon, fromLat, toLon, toLat float64) (float64, float64, error) {
				return 20000, 300, nil // reachable
			},
			expectNil:      false,
			expectPath:     []int{0, -1},
			expectSwitches: 0,
		},
		{
			name: "One hop then reachable",
			graph: &model.Graph{
				Nodes: []*model.Scooter{
					{ID: 0, Latitude: 1.0, Longitude: 1.0, Battery: 100},
					{ID: 1, Latitude: 1.0, Longitude: 1.1, Battery: 100},
				},
				Edges: map[int]map[int]*model.Edge{
					0: {
						1: &model.Edge{DistanceM: 1000, DurationSec: 60},
					},
					1: {},
				},
			},
			startNode: 0,
			targetLat: 1.2,
			targetLon: 1.2,
			mockDistanceFn: func(fromLon, fromLat, toLon, toLat float64) (float64, float64, error) {
				if fromLat == 1.0 && fromLon == 1.0 {
					return 50000, 1000, nil // not reachable directly
				}
				return 20000, 300, nil // reachable from node 1
			},
			expectNil:      false,
			expectPath:     []int{0, 1, -1},
			expectSwitches: 1,
		},
		{
			name: "No path available",
			graph: &model.Graph{
				Nodes: []*model.Scooter{
					{ID: 0, Latitude: 1.0, Longitude: 1.0, Battery: 10},
				},
				Edges: map[int]map[int]*model.Edge{},
			},
			startNode: 0,
			targetLat: 2.0,
			targetLon: 2.0,
			mockDistanceFn: func(fromLon, fromLat, toLon, toLat float64) (float64, float64, error) {
				return 100000, 5000, nil // too far
			},
			expectNil: true,
		},
		{
			name: "Distance API error",
			graph: &model.Graph{
				Nodes: []*model.Scooter{
					{ID: 0, Latitude: 1.0, Longitude: 1.0, Battery: 100},
				},
				Edges: map[int]map[int]*model.Edge{},
			},
			startNode: 0,
			targetLat: 1.1,
			targetLon: 1.1,
			mockDistanceFn: func(fromLon, fromLat, toLon, toLat float64) (float64, float64, error) {
				return 0, 0, errors.New("API failed")
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDistance = tt.mockDistanceFn

			result := FindOptimalPath(tt.graph, tt.startNode, tt.targetLat, tt.targetLon)

			if tt.expectNil {
				if result != nil {
					t.Errorf("Expected nil, got %+v", result)
				}
			} else {
				if result == nil {
					t.Fatal("Expected a valid path but got nil")
				}
				if len(result.Path) != len(tt.expectPath) {
					t.Errorf("Expected path %v, got %v", tt.expectPath, result.Path)
				}
				for i := range tt.expectPath {
					if result.Path[i] != tt.expectPath[i] {
						t.Errorf("Expected path %v, got %v", tt.expectPath, result.Path)
						break
					}
				}
				if result.SwitchCount != tt.expectSwitches {
					t.Errorf("Expected %d switches, got %d", tt.expectSwitches, result.SwitchCount)
				}
			}
		})
	}
}
