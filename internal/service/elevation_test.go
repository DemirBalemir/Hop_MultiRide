package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestAddElevationToScooters(t *testing.T) {
	tests := []struct {
		name          string
		scooters      []model.Scooter
		mockElevFunc  func(float64, float64, string) (float64, error)
		expectError   bool
		expectedValue float64
	}{
		{
			name: "successful elevation update",
			scooters: []model.Scooter{
				{ID: 1, Latitude: 39.9, Longitude: 32.8},
				{ID: 2, Latitude: 41.0, Longitude: 29.0},
			},
			mockElevFunc: func(lat, lon float64, key string) (float64, error) {
				return 111.0, nil
			},
			expectError:   false,
			expectedValue: 111.0,
		},
		{
			name: "elevation API returns error for one scooter",
			scooters: []model.Scooter{
				{ID: 3, Latitude: 0, Longitude: 0}, // fails
				{ID: 4, Latitude: 1, Longitude: 1}, // succeeds
			},
			mockElevFunc: func(lat, lon float64, key string) (float64, error) {
				if lat == 0 && lon == 0 {
					return 0, assert.AnError
				}
				return 222.0, nil
			},
			expectError:   false,
			expectedValue: 222.0, // only second one should be updated
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "scooters.json")

			// Write input scooters
			data, _ := json.MarshalIndent(tt.scooters, "", "  ")
			err := os.WriteFile(tmpFile, data, 0644)
			assert.NoError(t, err)

			// Mock the Get_Elev function
			original := Get_Elev
			defer func() { Get_Elev = original }()
			Get_Elev = tt.mockElevFunc

			// Set dummy API key
			t.Setenv("GOOGLE_MAPS_API_KEY", "testkey")

			// Run the actual function
			err = AddElevationToScooters(tmpFile, "testkey")
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Read result and verify
			updatedData, _ := os.ReadFile(tmpFile)
			var updated []model.Scooter
			_ = json.Unmarshal(updatedData, &updated)

			for _, scooter := range updated {
				if scooter.Latitude == 0 && scooter.Longitude == 0 {
					assert.Equal(t, 0.0, scooter.Elevation) // should be unchanged
				} else {
					assert.Equal(t, tt.expectedValue, scooter.Elevation)
				}
			}
		})
	}
}
