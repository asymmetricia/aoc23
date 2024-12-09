package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/canvas"
	"github.com/asymmetricia/aoc23/coord"
	"strconv"
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

func render(grid coord.World, anm coord.World, ant rune, lines []canvas.Line, mode int, node coord.Coord) *canvas.Canvas {
	cv := &canvas.Canvas{Timing: 5, Lines: lines}
	_, _, maxx, maxy := grid.Rect()
	cv.DrawRectangle(0, 0, maxx+2, maxy+2, aoc.TolVibrantBlue, aoc.Pixl)
	grid.Each(func(c coord.Coord) (stop bool) {
		col := aoc.TolVibrantGrey
		if grid.At(c) == ant {
			col = aoc.TolVibrantMagenta
		}
		v := grid.At(c)
		if v == '.' {
			v = anm.At(c)
			if v == '#' && mode == 0 {
				col = aoc.TolVibrantTeal
			}
		}
		if c == node {
			col = aoc.TolVibrantRed
		}
		cv.Set(c.X+1, c.Y+1, canvas.Cell{Color: col, Value: v})
		return
	})

	cv.DrawRectangle(maxx+3, 0, maxx+13, 4, aoc.TolVibrantBlue, aoc.Pixl)
	if mode == 0 {
		cv.PrintAt(maxx+4, 1, "Scanning", aoc.TolVibrantMagenta)
		cv.PrintAt(maxx+4, 2, "Counting", aoc.TolVibrantGrey)
		cv.PrintAt(maxx+4, 3, "Done", aoc.TolVibrantGrey)
	} else if mode > 0 {
		cv.PrintAt(maxx+4, 1, "Scanning", aoc.TolVibrantOrange)
		cv.PrintAt(maxx+4, 2, "Counting", aoc.TolVibrantMagenta)
		cv.PrintAt(maxx+4, 3, "Done", aoc.TolVibrantGrey)
		cv.DrawRectangle(maxx+3, 5, maxx+13, 8, aoc.TolVibrantBlue, aoc.Pixl)
		cv.PrintAt(maxx+4, 6, "Antinodes", aoc.TolVibrantTeal)
		cv.PrintAt(maxx+4, 7, strconv.Itoa(mode), aoc.TolVibrantTeal)
	} else if mode < 0 {
		cv.PrintAt(maxx+4, 1, "Scanning", aoc.TolVibrantOrange)
		cv.PrintAt(maxx+4, 2, "Counting", aoc.TolVibrantOrange)
		cv.PrintAt(maxx+4, 3, "Done", aoc.TolVibrantMagenta)
		cv.DrawRectangle(maxx+3, 5, maxx+13, 8, aoc.TolVibrantBlue, aoc.Pixl)
		cv.PrintAt(maxx+4, 6, "Antinodes", aoc.TolVibrantTeal)
		cv.PrintAt(maxx+4, 7, strconv.Itoa(-mode), aoc.TolVibrantTeal)
	}
	return cv
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

	var stack []*canvas.Canvas
	anm := grid.Copy()

	for ant, locs := range ants {
		var lines []canvas.Line

		minx, miny, maxx, maxy := grid.Rect()

		for _, a := range locs {
			for _, b := range locs {
				if a == b {
					continue
				}
				delta := b.Minus(a).Unit()
				minpt, maxpt := a, a
				for {
					anm.Set(minpt, '#')
					cdMin := minpt.Minus(delta)
					if cdMin.X >= minx && cdMin.X <= maxx &&
						cdMin.Y >= miny && cdMin.Y <= maxy {
						minpt = cdMin
					} else {
						break
					}
				}

				for {
					anm.Set(maxpt, '#')
					cdMax := maxpt.Plus(delta)
					if cdMax.X >= minx && cdMax.X <= maxx &&
						cdMax.Y >= miny && cdMax.Y <= maxy {
						maxpt = cdMax
					} else {
						break
					}
				}

				lines = append(lines, canvas.Line{
					A:     minpt.Plus(coord.C(1, 1)),
					B:     maxpt.Plus(coord.C(1, 1)),
					Color: aoc.TolVibrantMagenta,
				})
			}
		}

		stack = append(stack, render(grid, anm, ant, lines, 0, coord.C(-1, -1)))
	}

	antinodes := anm.Find('#')
	for i, node := range antinodes {
		if i == 0 {
			continue
		}
		if i < 50 ||
			i < 200 && i%5 == 0 ||
			i%10 == 0 {
			stack = append(stack, render(grid, anm, 0, nil, i, node))
		}
	}

	stack = append(stack, render(grid, anm, 0, nil, -len(antinodes), coord.C(-1, -1)))

	canvas.RenderGif(stack, "../../2024-08.gif", log)

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
