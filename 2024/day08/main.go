package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/coord"
	"strings"
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

	ants := map[rune][]coord.Coord{}

	grid.Each(func(c coord.Coord) (stop bool) {
		sym := grid.At(c)
		if sym == '.' {
			return
		}

		ants[sym] = append(ants[sym], c)
		return
	})

	var antinodes []coord.Coord
	grid.Each(func(c coord.Coord) (stop bool) {
		for ant, locs := range ants {
			for _, a := range locs {
				if c == a {
					continue
				}
				diff := c.Minus(a)
				b := c.Minus(diff).Minus(diff)
				if grid.At(b) == ant {
					if grid.At(c) == '.' {
						grid.Set(c, '#')
					}
					antinodes = append(antinodes, c)
					return
				}
			}
		}
		return
	})

	return len(antinodes)
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)

	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	grid := coord.Load(lines, coord.LoadConfig{Dense: true})

	ants := map[rune][]coord.Coord{}

	grid.Each(func(c coord.Coord) (stop bool) {
		sym := grid.At(c)
		if sym == '.' {
			return
		}

		ants[sym] = append(ants[sym], c)
		return
	})

	minx, miny, maxx, maxy := grid.Rect()
	var antinodes []coord.Coord
	anm := grid.Copy()
	grid.Each(func(c coord.Coord) (stop bool) {
		if grid.At(c) != '.' {
			anm.Set(c, '#')
			antinodes = append(antinodes, c)
			return
		}
		for ant, locs := range ants {
			for _, a := range locs {
				if c == a {
					continue
				}
				diff := a.Minus(c)
				diff = diff.Unit()
				cursor := a.Plus(diff)
				for cursor.X >= minx && cursor.X <= maxx &&
					cursor.Y >= miny && cursor.Y <= maxy {
					if grid.At(cursor) == ant {
						anm.Set(c, '#')
						antinodes = append(antinodes, c)
						return
					}
					cursor = cursor.Plus(diff)
				}
			}
		}
		return
	})

	anm.Print()
	return len(antinodes)
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 8)
	log.Printf("input solution A: %d", solutionA(input))
	log.Printf("input solution B: %d", solutionB(input))
}
