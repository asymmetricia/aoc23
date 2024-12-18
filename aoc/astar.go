package aoc

import (
	"github.com/asymmetricia/aoc23/coord"
	"github.com/asymmetricia/aoc23/search"
	"github.com/asymmetricia/aoc23/set"
)

// AStarGraph finds the path from start to end along the graph defined by edges
// returns from calling neighbors against each cell such that the path minimizes
// the total cost.
//
// If any callbacks are defined, they're called just before each time a cell is
// picked from the open set.
func AStarGraph[Cell comparable](
	start Cell,
	goal set.Set[Cell],
	neighbors func(a Cell) []Cell,
	cost func(a, b Cell) int,
	heuristic func(a Cell) int,
	callback ...search.CallbackFn[Cell],
) []Cell {
	return search.AStar(start,
		search.Callbacks(callback...),
		search.Cost(cost),
		search.GoalSet(goal),
		search.Heuristic(heuristic),
		search.Neighbors(neighbors),
	)
}

func AStarGrid(
	grid coord.World,
	start coord.Coord,
	goal set.Set[coord.Coord],
	cost func(from, to coord.Coord) int,
	heuristic func(from coord.Coord) int,
	diag bool,
	callback ...search.CallbackFn[coord.Coord]) []coord.Coord {
	return search.AStar(start,
		search.Callbacks(callback...),
		search.Cost(cost),
		search.GoalSet(goal),
		search.Grid(grid, diag),
		search.Heuristic(heuristic),
	)
}
