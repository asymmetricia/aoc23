package aoc

import (
	"github.com/asymmetricia/aoc23/search"
)

// Dijkstra implements a generic Dijkstra's Algorithm, which is guaranteed to
// find the shortest path from start to end, with edges given by repeated calls
// to neighbors().
//
// length should return the length of any given edge. callbacks are optional,
// used for status reporting or visualization.
//
// returns the path including the start Cell; or nil if there is no path.
//
// If Cell type implements an `(Cell) Equal(Cell) bool` method, then this is used
// to compare reachable cells to the end state. Otherwise, simple equality is
// used.
func Dijkstra[Cell comparable](
	start Cell,
	end Cell,
	neighbors func(a Cell) []Cell,
	length func(a, b Cell) int,
	callback ...search.DCallbackFn[Cell]) []Cell {

	return search.Dijkstra[Cell](start,
		search.Goal(end),
		search.Neighbors(neighbors),
		search.Cost(length),
		search.DCallbacks(callback...),
	)
}
