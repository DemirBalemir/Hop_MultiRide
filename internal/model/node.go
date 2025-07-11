package model

type Scooter struct {
	ID        int     `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Battery   int     `json:"battery"`
	Elevation float64 `json:"elevation,omitempty"` // omit if not set
}
