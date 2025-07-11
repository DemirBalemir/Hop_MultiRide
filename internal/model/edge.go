package model

type Edge struct {
	FromID                int
	ToID                  int
	DistanceM             float64
	DurationSec           float64
	ElevationDiff         float64
	EstimatedBatteryUsage float64
}
