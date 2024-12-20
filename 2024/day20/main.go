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

func solution(input []byte, cheatLen, target int) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	grid := coord.Loadv2(lines, coord.Dense())

	start := grid.Find('S')[0]
	costs := map[coord.Coord]int{}

	openSet := map[coord.Coord]bool{
		start: true,
	}
	for len(openSet) > 0 {
		var visit coord.Coord
		for k := range openSet {
			visit = k
			break
		}
		delete(openSet, visit)
		for _, neigh := range visit.Neighbors(false) {
			v := grid.At(neigh)
			if v == '.' || v == 'E' {
				c, ok := costs[neigh]
				tentative := costs[visit] + 1
				if !ok || tentative < c {
					costs[neigh] = tentative
					openSet[neigh] = true
				}
			}
		}
	}

	cheats := map[coord.Coord]map[coord.Coord]bool{}
	paths := append(grid.Find('.'), append(grid.Find('S'), grid.Find('E')...)...)
	for i, a := range paths {
		if _, ok := cheats[a]; !ok {
			cheats[a] = map[coord.Coord]bool{}
		}

		for _, b := range paths[i+1:] {
			dist := a.TaxiDistance(b)
			if dist > cheatLen {
				continue
			}

			if aoc.Abs(costs[a]-costs[b])-dist >= target {
				cheats[a][b] = true
			}
		}
	}

	var c int
	for _, cc := range cheats {
		c += len(cc)
	}
	return c
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 20)
	aStart := time.Now()
	aSoln := solution(input, 2, 100)
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solution(input, 20, 100)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
