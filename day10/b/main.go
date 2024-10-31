package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/canvas"
	"github.com/asymmetricia/aoc23/coord"
	"image/color"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func render(world coord.World, col func(c coord.Coord) color.Color) *canvas.Canvas {
	tiles := map[rune]rune{
		'I': 'I',
		'.': ' ',
		'S': 'S',
		'|': aoc.LineV,
		'-': aoc.LineH,
		'F': aoc.LineTL,
		'7': aoc.LineTR,
		'J': aoc.LineBR,
		'L': aoc.LineBL,
	}

	cv := &canvas.Canvas{}
	_, _, maxx, maxy := world.Rect()
	cv.DrawRectangle(0, 0, maxx+2, maxy+2, aoc.TolVibrantRed, 0)
	world.Each(func(c coord.Coord) (stop bool) {
		cv.Set(c.X+1, c.Y+1, canvas.Cell{col(c), tiles[world.At(c)], 0})
		return false
	})
	return cv
}

func solution(name string, input []byte) int {
	// anim stuff
	var stack []*canvas.Canvas
	f := 0

	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d %s lines (%d unique)", len(lines), name, len(uniq))

	world := coord.Load(lines, true)
	loop := &coord.DenseWorld{}
	src := world.Find('S')[0]
	loop.Set(src, 'S')

	cursor := src
	dir := coord.North
	if n := world.At(src.North()); n == '|' || n == 'F' || n == '7' {
		dir = coord.North
		cursor = src.North()
	}
	if e := world.At(src.East()); cursor == src && (e == '-' || e == 'J' || e == '7') {
		dir = coord.East
		cursor = src.East()
	}
	if w := world.At(src.West()); cursor == src && (w == '-' || w == 'F' || w == 'L') {
		dir = coord.West
		cursor = src.West()
	}
	if s := world.At(src.South()); cursor == src && (s == '|' || s == 'J' || s == 'L') {
		dir = coord.South
		cursor = src.South()
	}
	loop.Set(cursor, world.At(cursor))

	count := 1
	for cursor != src {
		if f%25 == 0 {
			stack = append(stack, render(world, func(c coord.Coord) color.Color {
				if loop.At(c) > 0 {
					return aoc.TolVibrantMagenta
				}
				return aoc.TolVibrantGrey
			}))
		}
		f++
		switch world.At(cursor) {
		case '|':
		case '-':
		case 'F':
			if dir == coord.North {
				dir = coord.East
			} else {
				dir = coord.South
			}
		case 'L':
			if dir == coord.South {
				dir = coord.East
			} else {
				dir = coord.North
			}
		case 'J':
			if dir == coord.East {
				dir = coord.North
			} else {
				dir = coord.West
			}
		case '7':
			if dir == coord.East {
				dir = coord.South
			} else {
				dir = coord.West
			}
		}
		cursor = cursor.Move(dir)
		loop.Set(cursor, world.At(cursor))
		count++
	}

	minx, miny, maxx, maxy := loop.Rect()

	cursor = coord.C(minx, miny)
	inside := false
	for cursor.X <= maxx && cursor.Y <= maxy {
		if f%50 == 0 {
			stack = append(stack, render(loop, func(c coord.Coord) color.Color {
				if loop.At(c) == 'I' {
					return aoc.TolVibrantCyan
				}
				return aoc.TolVibrantMagenta
			}))
		}
		f++

		switch loop.At(cursor) {
		case '|':
			inside = !inside
		case 'F':
			cursor.X++
			for loop.At(cursor) == '-' {
				cursor.X++
			}
			if loop.At(cursor) == 'J' {
				inside = !inside
			}
		case 'L':
			cursor.X++
			for loop.At(cursor) == '-' {
				cursor.X++
			}
			if loop.At(cursor) == '7' {
				inside = !inside
			}
		case 0:
			if inside {
				loop.Set(cursor, 'I')
			}
		}
		cursor.X++
		if cursor.X > maxx {
			inside = false
			cursor.X = minx
			cursor.Y++
		}
	}

	loop.Print()

	canvas.RenderGif(stack, "day10_"+name+".gif", log)
	return len(loop.Find('I'))
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	test, err := os.ReadFile("test")
	if err == nil {
		log.Printf("test solution: %d", solution("test", test))
	} else {
		log.Warningf("no test data present")
	}

	input := aoc.Input(2023, 10)
	log.Printf("input solution: %d", solution("input", input))
}
