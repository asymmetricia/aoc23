package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/coord"
	"golang.org/x/exp/slices"
	"os"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func countMirroredRows(w coord.DenseWorld, skip ...int) int {
	_, miny, _, maxy := w.Rect()
rows:
	for y := miny; y < maxy; y++ {
		if len(skip) > 0 && y == skip[0] {
			continue
		}
		if slices.Equal(w[y], w[y+1]) {
			for i := 1; i <= y; i++ {
				if y+i < maxy && !slices.Equal(w[y-i], w[y+i+1]) {
					continue rows
				}
			}
			return y
		}
	}
	return -1
}

func countMirroredCols(w coord.DenseWorld, skip ...int) int {
	transposed := coord.DenseWorld{}
	w.Each(func(c coord.Coord) bool {
		transposed.Set(coord.C(c.Y, c.X), w.At(c))
		return false
	})
	w = transposed
	return countMirroredRows(w, skip...)
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	blocks := aoc.Blocks(input)

	var total int
blocks:
	for i, block := range blocks {
		w := coord.Load(block, coord.LoadConfig{Dense: true}).(*coord.DenseWorld)
		c := countMirroredCols(*w)
		r := countMirroredRows(*w)

		minx, miny, maxx, maxy := w.Rect()
		for y := miny; y <= maxy; y++ {
			for x := minx; x <= maxx; x++ {
				w2 := w.Copy().(*coord.DenseWorld)
				if w2.At(coord.C(x, y)) == '#' {
					w2.Set(coord.C(x, y), '.')
				} else {
					w2.Set(coord.C(x, y), '#')
				}

				c2 := countMirroredCols(*w2, c)
				r2 := countMirroredRows(*w2, r)
				if c2 != -1 && c2 != c {
					log.Printf("%d: c %d -> c2 %d", i, c, c2)
					total += c2 + 1
					continue blocks
				}

				if r2 != -1 && r2 != r {
					log.Printf("%d: r %d -> r2 %d", i, r, r2)
					total += 100 * (r2 + 1)
					continue blocks
				}
			}
		}
		log.Fatalf("%d: no match!!!", i)
	}

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
