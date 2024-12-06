package main

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func isSafe(levels []int) bool {
	lt := levels[0] < levels[1]
	for i := 0; i < len(levels)-1; i++ {
		plt := levels[i] < levels[i+1]
		if plt != lt {
			return false
		}
		diff := aoc.Abs(levels[i] - levels[i+1])
		if diff < 1 || diff > 3 {
			return false
		}
	}
	return true
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	var safe int

	for _, line := range lines {
		levels := aoc.Ints(line)
		if isSafe(levels) {
			log.Printf("%v is safe", levels)
			safe++
			continue
		}
		for skip := 0; skip < len(levels); skip++ {
			dampened := make([]int, 0, len(levels))
			for i, l := range levels {
				if skip == i {
					continue
				}
				dampened = append(dampened, l)
			}
			if isSafe(dampened) {
				log.Printf("%v is safe with %d removed to yield %v", levels, skip, dampened)
				safe++
				break
			}
		}
	}

	return safe
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 02)
	log.Printf("input solution: %d", solution("input", input))
}