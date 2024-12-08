package coord

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Coord struct {
	X, Y int
}

// Distance returns the pythagorean distance. This is relatively slow to compute.
func (c Coord) Distance(d Coord) float64 {
	return math.Sqrt(math.Pow(float64(c.X-d.X), 2) + math.Pow(float64(c.Y-d.Y), 2))
}

// TaxiDistance returns the taxi / manhattan distance (i.e., absolute difference
// in X values plus absolute difference in Y values). It's quite fast to compute.
func (c Coord) TaxiDistance(d Coord) int {
	dx := c.X - d.X
	dy := c.Y - d.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func C(x, y int) Coord {
	return Coord{x, y}
}

func (c Coord) Neighbors(diag bool) []Coord {
	if diag {
		return []Coord{
			c.North(), c.NorthEast(),
			c.East(), c.SouthEast(),
			c.South(), c.SouthWest(),
			c.West(), c.NorthWest(),
		}
	}
	return []Coord{
		c.North(), c.East(), c.South(), c.West(),
	}
}

func MustFromComma(xy string) Coord {
	c, e := FromComma(xy)
	if e != nil {
		panic(e)
	}
	return c
}

func FromComma(xy string) (Coord, error) {
	parts := strings.Split(strings.TrimSpace(xy), ",")
	if len(parts) != 2 {
		return Coord{}, fmt.Errorf("expected two ,-separated parts, got %d", len(parts))
	}
	var ret Coord
	var err error
	ret.X, err = strconv.Atoi(parts[0])
	if err != nil {
		return Coord{}, fmt.Errorf("bad X coordinate %q: %w", parts[0], err)
	}
	ret.Y, err = strconv.Atoi(parts[1])
	if err != nil {
		return Coord{}, fmt.Errorf("bad Y coordinate %q: %w", parts[1], err)
	}
	return ret, nil
}

func (c Coord) Move(d Direction) Coord {
	switch d {
	case North:
		return c.North()
	case NorthEast:
		return c.NorthEast()
	case East:
		return c.East()
	case SouthEast:
		return c.SouthEast()
	case South:
		return c.South()
	case SouthWest:
		return c.SouthWest()
	case West:
		return c.West()
	case NorthWest:
		return c.NorthWest()
	}
	panic("bad direction " + strconv.Itoa(int(d)))
}

func (c Coord) North() Coord {
	return Coord{c.X, c.Y - 1}
}
func (c Coord) South() Coord {
	return Coord{c.X, c.Y + 1}
}
func (c Coord) East() Coord {
	return Coord{c.X + 1, c.Y}
}
func (c Coord) West() Coord {
	return Coord{c.X - 1, c.Y}
}
func (c Coord) NorthEast() Coord {
	return Coord{c.X + 1, c.Y - 1}
}
func (c Coord) SouthEast() Coord {
	return Coord{c.X + 1, c.Y + 1}
}
func (c Coord) NorthWest() Coord {
	return Coord{c.X - 1, c.Y - 1}
}
func (c Coord) SouthWest() Coord {
	return Coord{c.X - 1, c.Y + 1}
}

func (c Coord) Execute(steps []string) Coord {
	for _, step := range steps {
		if dir, ok := DirectionStrings[step]; ok {
			c = c.Move(dir)
		} else {
			panic(step)
		}
	}
	return c
}

func (c Coord) Plus(a Coord) Coord {
	return Coord{c.X + a.X, c.Y + a.Y}
}
func (c Coord) Equal(a Coord) bool {
	return c.X == a.X && c.Y == a.Y
}

// Minus returns c-a, the coordinate {c.X-a.X, c.Y-a.Y}
func (c Coord) Minus(a Coord) Coord {
	return Coord{c.X - a.X, c.Y - a.Y}
}

// Unit returns the unit vector in integer space of the given coordinate, i.e.,
// for (6,3) it returns (2,1), accounting properly for signs, etc.
func (c Coord) Unit() Coord {
	if c.X == 0 && c.Y == 0 {
		return Coord{0, 0}
	}
	if c.X == 0 && c.Y > 0 {
		return Coord{0, 1}
	}
	if c.X == 0 && c.Y < 0 {
		return Coord{0, -1}
	}
	if c.X > 0 && c.Y == 0 {
		return Coord{1, 0}
	}
	if c.X < 0 && c.Y == 0 {
		return Coord{-1, 0}
	}

	i := 2
	for {
		if c.X > 0 && c.X < i ||
			c.X < 0 && c.X > -i ||
			c.Y > 0 && c.Y < i ||
			c.Y < 0 && c.Y > -i {
			break
		}

		for c.X%i == 0 && c.Y%i == 0 {
			c.X /= i
			c.Y /= i
		}
		i++
	}

	return c
}

// TaxiPerimeter returns the set of points that are the given TaxiDistance away
// from the coordinate. This is a circle in taxi space, which looks like a
// diamond in Cartesian space.
func (c Coord) TaxiPerimeter(dist int) []Coord {
	if dist == 0 {
		return []Coord{c}
	}

	if dist < 0 {
		panic("negative distance")
	}

	var ret []Coord
	var cursor Coord = c
	cursor.Y -= dist
	for cursor.Y < c.Y {
		cursor = cursor.SouthEast()
		ret = append(ret, cursor)
	}
	for cursor.X > c.X {
		cursor = cursor.SouthWest()
		ret = append(ret, cursor)
	}
	for cursor.Y > c.Y {
		cursor = cursor.NorthWest()
		ret = append(ret, cursor)
	}
	for cursor.X < c.X {
		cursor = cursor.NorthEast()
		ret = append(ret, cursor)
	}
	return ret
}

// Collinear returns true if all the given coordinates are on the same line. If
// there are only one or two coordinates specified, this is by definition true.
func Collinear(cs ...Coord) bool {
	if len(cs) < 3 {
		return true
	}

	dir := cs[1].Minus(cs[0]).Unit()
	for _, c := range cs[2:] {
		if c.Minus(cs[0]).Unit() != dir {
			return false
		}
	}
	return true
}
