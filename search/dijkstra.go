package search

import (
	"github.com/asymmetricia/aoc23/pqueue"
)

// Dijkstra returns the shortest path including the start Cell; or nil if there
// is no path.
func Dijkstra[Cell comparable](start Cell, opts ...Option[Cell]) []Cell {
	config := NewConfig[Cell]()
	for _, opt := range opts {
		opt(config)
	}

	dist := map[Cell]int{}
	dist[start] = 0
	q := &pqueue.PQueue[Cell]{}
	q.AddWithPriority(start, 0)
	prev := map[Cell]Cell{}

	for q.Head != nil {
		u := q.Pop()

		for _, cb := range config.DCallbacks {
			cb(q, dist, prev, u)
		}

		if config.IsGoal(u) {
			path := []Cell{u}
			for u != start {
				u = prev[u]
				path = append(path, u)
			}
			for i := 0; i < len(path)/2; i++ {
				path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
			}
			return path
		}

		for _, neigh := range config.Neighbors(u) {
			tentativeDist := dist[u] + config.Cost(u, neigh)

			dv, ok := dist[neigh]
			if !ok || tentativeDist < dv {
				dist[neigh] = tentativeDist
				prev[neigh] = u
				q.AddWithPriority(neigh, tentativeDist+config.Heuristic(neigh, config.Goals()))
			}
		}
	}

	return nil
}
