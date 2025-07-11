package algorithm

import (
	"container/heap"
	"net/http"

	"github.com/demirbalemir/hop/Hop_MultiRide/internal/model"
	"github.com/demirbalemir/hop/Hop_MultiRide/internal/service"
)

type State struct {
	NodeID      int
	TimeSoFar   float64
	SwitchCount int
	Path        []int
	Index       int
}

var GetDistance = func(fromLon, fromLat, toLon, toLat float64) (float64, float64, error) {
	return service.GetOSRMDistance(http.DefaultClient, fromLon, fromLat, toLon, toLat)
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].SwitchCount == pq[j].SwitchCount {
		return pq[i].TimeSoFar < pq[j].TimeSoFar
	}
	return pq[i].SwitchCount < pq[j].SwitchCount
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i]; pq[i].Index, pq[j].Index = i, j }
func (pq *PriorityQueue) Push(x any) {
	item := x.(*State)
	item.Index = len(*pq)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() any {
	old := *pq
	item := old[len(old)-1]
	*pq = old[:len(old)-1]
	return item
}

func FindOptimalPath(graph *model.Graph, startNode int, targetLat, targetLon float64) *State {
	visited := make(map[int]bool)
	pq := &PriorityQueue{}
	heap.Init(pq)

	heap.Push(pq, &State{
		NodeID:      startNode,
		TimeSoFar:   0,
		SwitchCount: 0,
		Path:        []int{startNode},
	})

	for pq.Len() > 0 {
		curr := heap.Pop(pq).(*State)

		if visited[curr.NodeID] {
			continue
		}
		visited[curr.NodeID] = true

		from := graph.Nodes[curr.NodeID]
		maxRangeMeters := 35.0 * (float64(from.Battery) / 100.0) * 1000

		distToTarget, durToTarget, err := GetDistance(from.Longitude, from.Latitude, targetLon, targetLat)
		if err == nil && distToTarget <= maxRangeMeters {
			return &State{
				NodeID:      -1, // destination
				TimeSoFar:   curr.TimeSoFar + durToTarget,
				SwitchCount: curr.SwitchCount,
				Path:        append(curr.Path, -1),
			}
		}

		for neighborID, edge := range graph.Edges[curr.NodeID] {
			if visited[neighborID] {
				continue
			}

			if edge.DistanceM > maxRangeMeters {
				continue
			}

			newPath := append([]int{}, curr.Path...)
			newPath = append(newPath, neighborID)

			heap.Push(pq, &State{
				NodeID:      neighborID,
				TimeSoFar:   curr.TimeSoFar + edge.DurationSec,
				SwitchCount: curr.SwitchCount + 1,
				Path:        newPath,
			})
		}
	}

	return nil
}
