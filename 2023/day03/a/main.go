package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/asymmetricia/aoc23/canvas"
	"github.com/asymmetricia/aoc23/coord"
	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	var parts []int
	w := coord.Load(lines, coord.LoadConfig{Dense: true})
	var cnv canvas.Canvas
	w.Each(func(c coord.Coord) (stop bool) {
		cnv.Set(c.X, c.Y, canvas.Cell{Color: aoc.TolVibrantGrey, Value: w.At(c)})
		return false
	})
	w.Each(func(c coord.Coord) (stop bool) {
		v := w.At(c)
		if !unicode.IsDigit(v) {
			return false
		}

		// not the last digit
		if unicode.IsDigit(w.At(c.East())) {
			return false
		}

		var s string
		part := false
		cursor := c
		for unicode.IsDigit(w.At(cursor)) {
			s = string(w.At(cursor)) + s
			for _, neigh := range cursor.Neighbors(true) {
				nv := w.At(neigh)
				if nv == -1 || nv == '.' || nv == ' ' || unicode.IsDigit(nv) {
					continue
				}
				cnv.Set(neigh.X, neigh.Y, canvas.Cell{Color: aoc.TolVibrantMagenta, Value: w.At(neigh)})
				part = true
			}
			cursor = cursor.West()
		}

		if part {
			parts = append(parts, aoc.MustAtoi(s))
		}

		cursor = c
		for unicode.IsDigit(w.At(cursor)) {
			if part {
				cnv.Set(cursor.X, cursor.Y, canvas.Cell{Color: aoc.TolVibrantCyan, Value: w.At(cursor)})
			} else {
				cnv.Set(cursor.X, cursor.Y, canvas.Cell{Color: aoc.TolVibrantRed, Value: w.At(cursor)})
			}
			cursor = cursor.West()
		}

		return false
	})

	aoc.RenderPng(cnv.Render(), name+".png")
	accum := 0
	for _, part := range parts {
		accum += part
	}
	return accum
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

	input := aoc.Input(2023, 3)
	log.Printf("input solution: %d", solution("input", input))
}
