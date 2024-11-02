package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/coord"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d %s lines (%d unique)", len(lines), name, len(uniq))

	world := coord.Load(lines, coord.LoadConfig{Dense: true})
	loop := &coord.DenseWorld{}
	src := world.Find('S')[0]
	loop.Set(src, 'S')

	cursor := src
	dir := coord.North
	if n := world.At(src.North()); cursor == src && (n == '|' || n == 'F' || n == '7') {
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

	log.Print(count)

	return count / 2
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
