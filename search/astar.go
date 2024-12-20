package search

import (
	"math"
)

// AStar performs an A* search from the given start cell. Parameters of the problem are defined by passed in options.
//
// coord.World example:
//
//	grid := coord.Loadv2(lines, coord.Dense())
//	basePath := search.AStar(grid.Find('S')[0],
//		search.Grid(grid, false, '.', 'S', 'E'),
//		search.DistanceHeuristic(),
//		search.Goal(grid.Find('E')[0]))
//
// note!!! it's important to ensure that the start and end squares are considered walkable.
func AStar[Cell comparable](start Cell, opts ...Option[Cell]) []Cell {
	config := NewConfig[Cell]()
	for _, opt := range opts {
		opt(config)
	}

	if config.Neighbors == nil {
		panic("A* requires a neighbors function")
	}

	if config.Goals == nil || config.IsGoal == nil {
		panic("A* requires Goals and IsGoal")
	}

	openSet := map[Cell]bool{start: true}
	cameFrom := map[Cell]Cell{}
	gScore := map[Cell]int{
		start: 0,
	}
	fScore := map[Cell]int{
		start: config.Heuristic(start, config.Goals()),
	}

	var bestGScore = math.MaxInt

	var current Cell
	for len(openSet) > 0 {
		var curFScore = math.MaxInt
		first := true

		// get the item in the open set with the lowest fScore
		for c := range openSet {
			fs, ok := fScore[c]
			if !ok {
				fs = math.MaxInt
			}

			if first || fs < curFScore {
				first = false
				current = c
				curFScore = fs
			}
		}
		delete(openSet, current)

		for _, cb := range config.Callbacks {
			cb(openSet, cameFrom, gScore, fScore, current)
		}

		if config.IsGoal(current) {
			if gScore[current] < bestGScore {
				bestGScore = gScore[current]
			}

			break
		}

		curGS := gScore[current]
		for _, neighbor := range config.Neighbors(current) {
			neighGS, ok := gScore[neighbor]
			tentativeGScore := curGS + config.Cost(current, neighbor)
			if (!ok || tentativeGScore < neighGS) && tentativeGScore < config.MaxDistance {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = tentativeGScore + config.Heuristic(neighbor, config.Goals())
				openSet[neighbor] = true
			}
		}
	}

	if !config.IsGoal(current) {
		return nil
	}

	ret := []Cell{current}
	cursor := current
	for cursor != start {
		cursor = cameFrom[cursor]
		ret = append(ret, cursor)
	}

	for i := 0; i < len(ret)/2; i++ {
		ret[i], ret[len(ret)-1-i] = ret[len(ret)-1-i], ret[i]
	}

	return ret
}
