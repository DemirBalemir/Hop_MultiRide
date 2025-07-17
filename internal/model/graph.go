package model

type Graph struct {
	Nodes map[int]*Scooter
	Edges map[int]map[int]*Edge
}
