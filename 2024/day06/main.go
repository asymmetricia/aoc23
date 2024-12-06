package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/coord"
	"github.com/asymmetricia/aoc23/term"
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

	var guard coord.Coord
	var dir coord.Direction
	grid.Each(func(c coord.Coord) (stop bool) {
		switch grid.At(c) {
		case '^':
			dir = coord.North
			guard = c
		case '>':
			dir = coord.East
			guard = c
		case 'v':
			dir = coord.South
			guard = c
		case '<':
			dir = coord.West
			guard = c
		}
		return
	})

	minx, miny, maxx, maxy := grid.Rect()
	for minx <= guard.X && maxx >= guard.X &&
		miny <= guard.Y && maxy >= guard.Y {
		if nv := grid.At(guard.Move(dir)); nv != '#' {
			grid.Set(guard, 'X')
			guard = guard.Move(dir)
		} else {
			dir = dir.CW(false)
		}
		//term.MoveCursor(1, 1)
		//grid.Print()
		//term.ColorC(aoc.TolVibrantMagenta)
		//term.MoveCursor(guard.X+1, guard.Y+1)
		//print("@")
		//term.ColorReset()
		//time.Sleep(10 * time.Millisecond)
	}

	return len(grid.Find('X'))
}

func loop(grid coord.World, guard coord.Coord, dir coord.Direction) bool {
	type state struct {
		c coord.Coord
		d coord.Direction
	}
	states := map[state]bool{
		state{guard, dir}: true,
	}

	start := guard

	minx, miny, maxx, maxy := grid.Rect()
	for minx <= guard.X && maxx >= guard.X &&
		miny <= guard.Y && maxy >= guard.Y {
		next := guard.Move(dir)

		if nv := grid.At(next); nv != '#' {
			guard = guard.Move(dir)
		} else {
			dir = dir.CW()
		}
		sk := state{guard, dir}
		if states[sk] {
			grid.Set(start, '@')
			grid.Set(guard, '!')
			return true
		}
		states[sk] = true
	}
	return false
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	grid := coord.Load(lines, coord.LoadConfig{Dense: true})

	var guard coord.Coord
	var dir coord.Direction
	grid.Each(func(c coord.Coord) (stop bool) {
		switch grid.At(c) {
		case '^':
			dir = coord.North
			guard = c
		case '>':
			dir = coord.East
			guard = c
		case 'v':
			dir = coord.South
			guard = c
		case '<':
			dir = coord.West
			guard = c
		}
		return
	})
	start := guard
	startDir := dir

	count := 0
	tried := map[coord.Coord]bool{}
	minx, miny, maxx, maxy := grid.Rect()
	for minx <= guard.X && maxx >= guard.X &&
		miny <= guard.Y && maxy >= guard.Y {

		if guard != start && !tried[guard] {
			tried[guard] = true
			maybe := grid.Copy()
			maybe.Set(guard, '#')
			if loop(maybe, start, startDir) {
				count++
			}
		}

		next := guard.Move(dir)

		if grid.At(next) == '#' {
			dir = dir.CW()
			continue
		}

		guard = guard.Move(dir)
	}

	return count
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	term.HideCursor()

	input := aoc.Input(2024, 06)
	log.Printf("input solution A: %d", solutionA(input))

	solB := solutionB(input)
	if solB >= 2141 {
		panic("too big")
	}
	log.Printf("input solution B: %d", solB)

	term.ShowCursor()
}
