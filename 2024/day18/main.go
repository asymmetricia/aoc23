package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/coord"
	"github.com/asymmetricia/aoc23/search"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solutionA(input []byte, test bool) int {
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	dim := 71
	if test {
		dim = 7
	}
	grid := &coord.DenseWorld{}
	for y := 0; y < dim; y++ {
		row := make([]rune, dim)
		for i := range row {
			row[i] = '.'
		}
		*grid = append(*grid, row)
	}

	count := 1024
	if test {
		count = 12
	}
	for i := 0; i < count; i++ {
		c := coord.MustFromComma(lines[i])
		grid.Set(c, aoc.BlockFull)
	}

	path := search.AStar(coord.C(0, 0),
		search.Grid(grid, false),
		search.Goal(coord.C(dim-1, dim-1)),
		search.DistanceHeuristic(),
	)
	return len(path) - 1
}

func solutionB(input []byte, test bool) string {
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	dim := 71
	if test {
		dim = 7
	}
	grid := &coord.DenseWorld{}
	for y := 0; y < dim; y++ {
		row := make([]rune, dim)
		for i := range row {
			row[i] = '.'
		}
		*grid = append(*grid, row)
	}

	count := 1024
	if test {
		count = 12
	}
	for i := 0; i < count; i++ {
		c := coord.MustFromComma(lines[i])
		grid.Set(c, aoc.BlockFull)
	}

	i := count
	for {
		grid.Set(coord.MustFromComma(lines[i]), aoc.BlockFull)

		path := search.Dijkstra(coord.C(0, 0),
			search.Goal(coord.C(dim-1, dim-1)),
			search.Grid(grid, false))
		if path == nil {
			return lines[i]
		}
		i++
	}
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 18)
	aStart := time.Now()
	aSoln := solutionA(input, false)
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input, false)
	log.Printf("input solution B: %s (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
