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

func solution(name string, input []byte) int64 {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	w := coord.Load(lines, coord.LoadConfig{Dense: true})
	var cnv canvas.Canvas
	w.Each(func(c coord.Coord) (stop bool) {
		cnv.Set(c.X, c.Y, canvas.Cell{Color: aoc.TolVibrantGrey, Value: w.At(c)})
		return false
	})

	var accum int64

	gears := w.Find('*')
	for _, gear := range gears {
		first, second := coord.C(-1, -1), coord.C(-1, -1)
		cnv.Set(gear.X, gear.Y, canvas.Cell{Color: aoc.TolVibrantMagenta, Value: w.At(gear)})
		for _, neigh := range gear.Neighbors(true) {
			if unicode.IsDigit(w.At(neigh)) {
				cursor := neigh
				for unicode.IsDigit(w.At(cursor.West())) {
					cursor = cursor.West()
				}
				if first.X == -1 {
					first = cursor
				} else if first != cursor {
					second = cursor
				}
			}
		}
		if first.X != -1 && second.X != -1 {
			cursor := first
			for unicode.IsDigit(w.At(cursor.West())) {
				cursor = cursor.West()
			}
			fv := int(w.At(cursor) - '0')
			cnv.Set(cursor.X, cursor.Y, canvas.Cell{Color: aoc.TolVibrantCyan, Value: w.At(cursor)})
			for unicode.IsDigit(w.At(cursor.East())) {
				cursor = cursor.East()
				cnv.Set(cursor.X, cursor.Y, canvas.Cell{Color: aoc.TolVibrantCyan, Value: w.At(cursor)})
				fv = fv*10 + int(w.At(cursor)-'0')
			}

			cursor = second
			for unicode.IsDigit(w.At(cursor.West())) {
				cursor = cursor.West()
			}
			sv := int(w.At(cursor) - '0')
			cnv.Set(cursor.X, cursor.Y, canvas.Cell{Color: aoc.TolVibrantCyan, Value: w.At(cursor)})
			for unicode.IsDigit(w.At(cursor.East())) {
				cursor = cursor.East()
				cnv.Set(cursor.X, cursor.Y, canvas.Cell{Color: aoc.TolVibrantCyan, Value: w.At(cursor)})
				sv = sv*10 + int(w.At(cursor)-'0')
			}

			accum += int64(fv) * int64(sv)
		}
	}

	aoc.RenderPng(cnv.Render(), name+".png")

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
