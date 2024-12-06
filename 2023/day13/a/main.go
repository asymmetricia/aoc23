package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/coord"
	"golang.org/x/exp/slices"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func countMirroredRows(w coord.DenseWorld) int {
	_, miny, _, maxy := w.Rect()
rows:
	for y := miny; y < maxy; y++ {
		if slices.Equal(w[y], w[y+1]) {
			for i := 1; i <= y; i++ {
				if y+i < maxy && !slices.Equal(w[y-i], w[y+i+1]) {
					continue rows
				}
			}
			return y + 1
		}
	}
	return 0
}

func countMirroredCols(w coord.DenseWorld) int {
	transposed := coord.DenseWorld{}
	w.Each(func(c coord.Coord) bool {
		transposed.Set(coord.C(c.Y, c.X), w.At(c))
		return false
	})
	w = transposed
	return countMirroredRows(w)
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	var blocks [][]string
	var accum []string
	for _, line := range lines {
		if line == "" {
			if len(accum) > 0 {
				blocks = append(blocks, accum)
				accum = nil
			}
		} else {
			accum = append(accum, line)
		}
	}
	if len(accum) > 0 {
		blocks = append(blocks, accum)
	}

	var total int
	for _, block := range blocks {
		w := coord.Load(block, coord.LoadConfig{Dense: true}).(*coord.DenseWorld)
		total += countMirroredCols(*w) + 100*countMirroredRows(*w)
	}

	//for y := miny; y < maxy-1; y++ {
	//	if slices.Equal(world[y], world[y+1]) {
	//	}
	//}

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

	input := aoc.Input(2023, 13)
	log.Printf("input solution: %d", solution("input", input))
}
