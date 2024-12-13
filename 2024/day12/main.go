package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/coord"
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

	var plots []map[coord.Coord]bool
	grid.Each(func(c coord.Coord) (stop bool) {
		cv := grid.At(c)
		if cv == 0 {
			return
		}

		plot := map[coord.Coord]bool{c: true}
		visit := []coord.Coord{c}
		grid.Set(c, 0)
		for len(visit) > 0 {
			cursor := visit[0]
			for _, n := range cursor.Neighbors(false) {
				if grid.At(n) == cv {
					plot[n] = true
					grid.Set(n, 0)
					visit = append(visit, n)
				}
			}
			visit = visit[1:]
		}
		plots = append(plots, plot)
		return
	})

	log.Printf("%d plots", len(plots))

	var total int
	for _, plot := range plots {
		var perimeter int
		for c := range plot {
			for _, n := range c.Neighbors(false) {
				if !plot[n] {
					perimeter++
				}
			}
		}
		total += perimeter * len(plot)
	}

	return total
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	grid := coord.Load(lines, coord.LoadConfig{Dense: true})
	plotted := grid.Copy()

	var plots []map[coord.Coord]bool
	grid.Each(func(c coord.Coord) (stop bool) {
		if plotted.At(c) == 0 {
			return
		}

		cv := grid.At(c)
		plot := map[coord.Coord]bool{c: true}
		visit := []coord.Coord{c}
		plotted.Set(c, 0)
		for len(visit) > 0 {
			cursor := visit[0]
			for _, n := range cursor.Neighbors(false) {
				if plotted.At(n) != 0 && grid.At(n) == cv {
					plot[n] = true
					plotted.Set(n, 0)
					visit = append(visit, n)
				}
			}
			visit = visit[1:]
		}
		plots = append(plots, plot)
		return
	})

	log.Printf("%d plots", len(plots))

	var total int
	for _, plot := range plots {
		fences := map[coord.Coord]map[coord.Direction]bool{}
		for plant := range plot {
			fences[plant] = map[coord.Direction]bool{}
			for _, ord := range []coord.Direction{coord.North, coord.East, coord.South, coord.West} {
				n := plant.Move(ord)
				if !plot[n] {
					fences[plant][ord] = true
				}
			}
		}

	fences:
		for len(fences) > 0 {
			for loc, sides := range fences {
				if len(sides) == 0 {
					delete(fences, loc)
					continue fences
				}

				for side := range sides {
					var neighborDirs []coord.Direction

					switch side {
					case coord.North, coord.South:
						neighborDirs = []coord.Direction{coord.West, coord.East}
					case coord.West, coord.East:
						neighborDirs = []coord.Direction{coord.North, coord.South}
					default:
						panic(side)
					}

					delete(fences[loc], side)
					total += len(plot)

				neighborDirs:
					for _, dir := range neighborDirs {
						cursor := loc.Move(dir)
						_, ok := fences[cursor]
						for ok {
							if fences[cursor][side] {
								delete(fences[cursor], side)
							} else {
								continue neighborDirs
							}
							cursor = cursor.Move(dir)
							_, ok = fences[cursor]
						}
					}

				}
			}
		}
	}

	return total
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 12)
	aStart := time.Now()
	aSoln := solutionA(input)
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
