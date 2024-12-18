package search

import (
	"github.com/asymmetricia/aoc23/coord"
	"github.com/asymmetricia/aoc23/pqueue"
	"github.com/asymmetricia/aoc23/set"
	"iter"
	"maps"
	"math"
)

type CallbackFn[Cell comparable] func(
	openSet map[Cell]bool,
	cameFrom map[Cell]Cell,
	gScore map[Cell]int,
	fScore map[Cell]int,
	current Cell)

type Config[Cell comparable] struct {
	IsGoal       func(Cell) bool
	Goals        func() iter.Seq[Cell]
	Neighbors    func(a Cell) []Cell
	Cost         func(a, b Cell) int
	NegativeCost bool
	Heuristic    func(a Cell, goals iter.Seq[Cell]) int

	DCallbacks []DCallbackFn[Cell]
	Callbacks  []CallbackFn[Cell]
}

func NewConfig[Cell comparable]() *Config[Cell] {
	return &Config[Cell]{
		Cost: func(_, _ Cell) int {
			return 1
		},
		Heuristic: func(_ Cell, _ iter.Seq[Cell]) int {
			return 0
		},
	}
}

type Option[Cell comparable] func(config *Config[Cell])

func Callback[Cell comparable](c CallbackFn[Cell]) Option[Cell] {
	return func(config *Config[Cell]) {
		config.Callbacks = append(config.Callbacks, c)
	}
}

func Callbacks[Cell comparable](c ...CallbackFn[Cell]) Option[Cell] {
	return func(config *Config[Cell]) {
		config.Callbacks = c
	}
}

func Cost[Cell comparable](f func(a Cell, b Cell) int) Option[Cell] {
	return func(config *Config[Cell]) {
		if f != nil {
			config.Cost = f
		}
	}
}

func Goal[Cell comparable](g Cell) Option[Cell] {
	return func(config *Config[Cell]) {
		config.IsGoal = func(cell Cell) bool {
			return cell == g
		}
		config.Goals = func() iter.Seq[Cell] {
			return func(yield func(Cell) bool) {
				yield(g)
			}
		}
	}
}

func GoalsFn[Cell comparable](f func() []Cell) Option[Cell] {
	return func(config *Config[Cell]) {
		config.Goals = func() iter.Seq[Cell] {
			return func(yield func(Cell) bool) {
				for _, c := range f() {
					if !yield(c) {
						return
					}
				}
			}
		}
	}
}

func Goals[Cell comparable](g ...Cell) Option[Cell] {
	s := set.FromItems(g)
	return func(config *Config[Cell]) {
		config.IsGoal = func(cell Cell) bool {
			return s[cell]
		}
		config.Goals = func() iter.Seq[Cell] {
			return maps.Keys(s)
		}
	}
}

func GoalSet[Cell comparable](gs set.Set[Cell]) Option[Cell] {
	return func(config *Config[Cell]) {
		config.IsGoal = func(cell Cell) bool {
			return gs[cell]
		}
		config.Goals = func() iter.Seq[Cell] {
			return maps.Keys(gs)
		}
	}
}

// Grid configures A* to operate on a grid. The Neighbors function is set to identify
// grid neighbors with value `.` by default
func Grid(grid coord.World, diag bool, walkable ...rune) Option[coord.Coord] {
	var s set.Set[rune]
	if len(walkable) > 0 {
		s = set.FromItems(walkable)
	} else {
		s = set.FromItem('.')
	}

	return func(config *Config[coord.Coord]) {
		config.Neighbors = func(a coord.Coord) (ret []coord.Coord) {
			for _, n := range a.Neighbors(diag) {
				if s[grid.At(n)] {
					ret = append(ret, n)
				}
			}
			return
		}
	}
}

func Heuristic[Cell comparable](n func(Cell) int) Option[Cell] {
	return func(config *Config[Cell]) {
		config.Heuristic = func(a Cell, _ iter.Seq[Cell]) int {
			return n(a)
		}
	}
}

func IsGoal[Cell comparable](f func(Cell) bool) Option[Cell] {
	return func(config *Config[Cell]) {
		config.IsGoal = f
	}
}

func DistanceHeuristic() Option[coord.Coord] {
	return func(config *Config[coord.Coord]) {
		config.Heuristic = func(a coord.Coord, goals iter.Seq[coord.Coord]) int {
			best := math.MaxFloat32
			for goal := range goals {
				if d := a.Distance(goal); d < best {
					best = d
				}
			}
			return int(best)
		}
	}
}

func Neighbors[Cell comparable](n func(Cell) []Cell) Option[Cell] {
	return func(config *Config[Cell]) {
		config.Neighbors = n
	}
}

type DCallbackFn[Cell comparable] func(
	q *pqueue.PQueue[Cell],
	dist map[Cell]int,
	prev map[Cell]Cell,
	current Cell)

func DCallbacks[Cell comparable](f ...DCallbackFn[Cell]) Option[Cell] {
	return func(config *Config[Cell]) {
		config.DCallbacks = f
	}
}