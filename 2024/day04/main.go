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

func check(grid coord.World, c coord.Coord, dir coord.Direction, word []rune) bool {
	if len(word) == 0 {
		return true
	}
	if grid.At(c) != word[0] {
		return false
	}
	return check(grid, c.Move(dir), dir, word[1:])
}

func solutionA(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	grid := coord.Load(lines, coord.LoadConfig{Dense: true})

	count := 0

	grid.Each(func(c coord.Coord) bool {
		for _, d := range coord.Directions {
			if check(grid, c, d, []rune("XMAS")) {
				count++
			}
		}
		return false
	})

	return count
}

func solutionB(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)

	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	grid := coord.Load(lines, coord.LoadConfig{Dense: true})

	count := 0

	grid.Each(func(c coord.Coord) (stop bool) {
		if grid.At(c) != 'A' {
			return
		}

		nw, ne, sw, se := grid.At(c.NorthWest()), grid.At(c.NorthEast()), grid.At(c.SouthWest()), grid.At(c.SouthEast())

		if !(nw == 'M' && se == 'S' || nw == 'S' && se == 'M') {
			return
		}

		if !(ne == 'M' && sw == 'S' || ne == 'S' && sw == 'M') {
			return
		}

		count++

		return
	})

	return count
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 04)
	log.Printf("input solution A: %d", solutionA("input", input))
	log.Printf("input solution B: %d", solutionB("input", input))
}
