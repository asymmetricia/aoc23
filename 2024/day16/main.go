package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/canvas"
	"github.com/asymmetricia/aoc23/coord"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solutionA(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	grid := coord.Load(lines, coord.LoadConfig{Dense: true})

	start := grid.Find('S')[0]

	costs := map[coord.Coord]int{}
	from := map[coord.Coord]coord.Coord{}
	facing := map[coord.Coord]coord.Direction{start: coord.East}

	toVisit := []coord.Coord{start}

	for len(toVisit) > 0 {
		slices.SortFunc(toVisit, func(a, b coord.Coord) bool {
			return costs[a] < costs[b]
		})
		visit := toVisit[0]
		toVisit = toVisit[1:]

		for _, dir := range coord.CardinalDirections {
			nbor := visit.Move(dir)
			if grid.At(nbor) == '#' {
				continue
			}
			if _, ok := costs[nbor]; ok {
				continue
			}
			costs[nbor] = costs[visit] + 1
			if dir != facing[visit] {
				costs[nbor] += 1000
			}
			facing[nbor] = dir
			from[nbor] = visit
			if grid.At(nbor) == 'E' {
				toVisit = nil
				break
			}
			toVisit = append(toVisit, nbor)
		}
	}

	return costs[grid.Find('E')[0]]
}

type Position struct {
	Facing   coord.Direction
	Position coord.Coord
}

func (p Position) Neighbors(grid coord.World) []Position {
	ret := []Position{
		{p.Facing.CW(), p.Position},
		{p.Facing.CW().CW(), p.Position},
		{p.Facing.CW().CW().CW(), p.Position},
	}
	if grid.At(p.Position.Move(p.Facing)) != '#' {
		ret = append(ret, Position{
			Facing:   p.Facing,
			Position: p.Position.Move(p.Facing),
		})
	}
	return ret
}

func computeCosts(grid coord.World, s rune, enc *aoc.MP4Encoder) map[Position]int {
	start := grid.Find(s)[0]
	costs := map[Position]int{}
	toVisit := map[Position]bool{}

	if s == 'S' {
		costs[Position{coord.East, start}] = 0
		toVisit[Position{coord.East, start}] = true
	} else {
		for _, d := range coord.CardinalDirections {
			costs[Position{d, start}] = 0
			toVisit[Position{d, start}] = true
		}
	}

	var visits int
	for len(toVisit) > 0 {
		visits++
		k := maps.Keys(toVisit)
		slices.SortFunc(k, func(a, b Position) bool {
			return costs[a] < costs[b]
		})

		visit := k[0]
		delete(toVisit, visit)

		if visits%50 == 0 {
			cv := &canvas.Canvas{Palette: append(aoc.TolVibrant, aoc.TolIncandescent...)}
			grid.Each(func(c coord.Coord) (stop bool) {
				col := aoc.TolVibrantGrey
				v := grid.At(c)
				switch v {
				case '#':
					v = aoc.BlockLight
				case 'S', 'E':
					col = aoc.TolVibrantMagenta
				}
				if visit.Position == c {
					col = aoc.TolVibrantMagenta
				} else {
					for _, dir := range coord.CardinalDirections {
						pos := Position{dir, c}
						c, ok := costs[pos]
						if ok {
							col = aoc.TolScale(0, 115500, c, aoc.TolIncandescent)
							v = aoc.BlockFull
							break
						} else if toVisit[pos] {
							col = aoc.TolVibrantRed
							v = aoc.BlockMedium
						}
					}
				}
				cv.Set(c.X, c.Y, canvas.Cell{Color: col, Value: v})
				return
			})
			if err := enc.Encode(cv.Render()); err != nil {
				log.Fatal(err)
			}
		}

		for _, neigh := range visit.Neighbors(grid) {
			c, ok := costs[neigh]

			cc := costs[visit]
			if neigh.Facing == visit.Facing {
				cc++
			} else {
				cc += 1000
			}

			if !ok || cc < c {
				costs[neigh] = cc
				toVisit[neigh] = true
			}
		}
	}

	return costs
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	grid := coord.Load(lines, coord.LoadConfig{Dense: true})

	start, _ := grid.Find('S')[0], grid.Find('E')[0]

	enc, err := aoc.NewMP4Encoder("../../2024-16.mp4", 60, log)
	if err != nil {
		log.Fatal(err)
	}

	// compute costs to reach each square starting at S
	fwdCosts := computeCosts(grid, 'S', enc)
	revCosts := computeCosts(grid, 'E', enc)
	best := revCosts[Position{coord.East, start}]

	for _, dir := range coord.CardinalDirections {
		c := start.North()
		p := Position{dir, c}
		rp := Position{dir.Opposite(), c}
		log.Print(dir, fwdCosts[p], revCosts[rp], fwdCosts[p]+revCosts[rp])
	}

	cv := &canvas.Canvas{Palette: append(aoc.TolVibrant, aoc.TolIncandescent...)}
	grid.Each(func(c coord.Coord) (stop bool) {
		col := aoc.TolVibrantGrey
		v := grid.At(c)
		switch v {
		case '#':
			v = aoc.BlockLight
		case 'S', 'E':
			col = aoc.TolVibrantMagenta
		}
		cv.Set(c.X, c.Y, canvas.Cell{Color: col, Value: v})
		return
	})
	if err := enc.Encode(cv.Render()); err != nil {
		log.Fatal(err)
	}

	grid.Each(func(c coord.Coord) (stop bool) {
		for _, dir := range coord.CardinalDirections {
			if fwdCosts[Position{dir, c}]+revCosts[Position{dir.Opposite(), c}] == best {
				grid.Set(c, aoc.BlockFull)
				cv.Set(c.X, c.Y, canvas.Cell{Color: aoc.TolVibrantTeal, Value: aoc.BlockFull})
				if err := enc.Encode(cv.Render()); err != nil {
					log.Fatal(err)
				}
			}
		}
		return
	})

	if err := enc.Close(); err != nil {
		log.Fatal(err)
	}

	log.Print(best)
	return len(grid.Find(aoc.BlockFull))
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 16)
	aStart := time.Now()
	aSoln := solutionA(input)
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
