package main

import (
	"bytes"
	"fmt"
	"github.com/asymmetricia/aoc23/coord"
	"github.com/sirupsen/logrus"
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

var numPad = coord.DenseWorld{
	{'7', '8', '9'},
	{'4', '5', '6'},
	{'1', '2', '3'},
	{0, '0', 'A'},
}

var dirPad = coord.DenseWorld{
	{0, '^', 'A'},
	{'<', 'v', '>'},
}

func path(grid coord.DenseWorld, start, end rune) []string {
	if start == end {
		return []string{"A"}
	}

	froms := grid.Find(start)
	if len(froms) == 0 {
		panic(fmt.Sprintf("could not find %c (%d) on grid", start, start))
	}
	from := froms[0]

	to := grid.Find(end)[0]

	var ret []string

	for _, opt := range []struct {
		b bool
		d coord.Direction
		r string
	}{
		{from.Y < to.Y, coord.South, "v"},
		{from.Y > to.Y, coord.North, "^"},
		{from.X < to.X, coord.East, ">"},
		{from.X > to.X, coord.West, "<"},
	} {
		if !opt.b {
			continue
		}
		nv := grid.At(from.Move(opt.d))
		if nv == 0 {
			continue
		}

		for _, p := range path(grid, nv, end) {
			ret = append(ret, opt.r+p)
		}
	}

	return ret
}

var cost func(c string, depth int) int

func init() {
	cost = aoc.Cache2(costImpl)
}

// cost computes the cost to activate button `c` at depth `depth`. Depth 0 is the
// human, higher depths (1 and 2, in part A) are robots. Does not handle door
// presses, because they don't loop to "A" but also are not free.
func costImpl(c string, depth int) int {
	if depth == 0 {
		return len(c)
	}

	var sum int
	var ret string
	cursor := 'A'
	for _, r := range c {
		nextPaths := path(dirPad, cursor, r)
		minCost := math.MaxInt
		var minPath string
		for _, nextPath := range nextPaths {
			nextPathCost := cost(nextPath, depth-1)
			if nextPathCost < minCost {
				minCost = nextPathCost
			}
		}
		sum += minCost
		ret = ret + minPath
		cursor = r
	}
	return sum
}

func solutionA(input []byte) int {
	return solve(input, 2)
}

func solve(input []byte, numPadRobotCount int) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	var value int
	cursor := 'A'
	for _, line := range lines {
		num := aoc.MustAtoi(strings.TrimPrefix(strings.TrimSuffix(line, "A"), "0"))
		var length int
		var seq string
		for _, c := range line {
			bestPathCost := math.MaxInt
			var bestPathStr string
			numPaths := path(numPad, cursor, c)
			for _, numPath := range numPaths {
				candidateCost := cost(numPath, numPadRobotCount)
				//log.Printf("%c to %c - %s - %d - %s", cursor, c, string(numPath), candidateCost, candidateStr)
				if candidateCost < bestPathCost {
					bestPathCost = candidateCost
				}
			}
			cursor = c
			length += bestPathCost
			seq = seq + bestPathStr
		}
		//log.Printf("%s - %s", line, seq)
		value += num * length
	}

	return value
}

func solutionB(input []byte) int {
	return solve(input, 25)
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 21)
	aStart := time.Now()
	aSoln := solutionA(input)
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
