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

	w := coord.Load(lines, coord.LoadConfig{Dense: true})
	_, _, _, maxy := w.Rect()

	changed := true
	for changed {
		changed = false
		w.Each(func(c coord.Coord) (stop bool) {
			if w.At(c) == 'O' && w.At(c.North()) == '.' {
				changed = true
				w.Set(c.North(), 'O')
				w.Set(c, '.')
			}

			return false
		})
	}

	var total int
	w.Each(func(c coord.Coord) (stop bool) {
		if w.At(c) == 'O' {
			total += maxy - c.Y + 1
		}
		return false
	})

	return total
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

	input := aoc.Input(2023, 14)
	log.Printf("input solution: %d", solution("input", input))
}
