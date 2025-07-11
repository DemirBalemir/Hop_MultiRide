package model

type Graph struct {
	Nodes []*Scooter
	Edges map[int]map[int]*Edge
}
