package main

import (
	"bytes"
	"slices"
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
	var a, b []int
	for _, line := range lines {
		fields := strings.Fields(line)
		a = append(a, aoc.Int(fields[0]))
		b = append(b, aoc.Int(fields[1]))
	}

	slices.Sort(a)
	slices.Sort(b)

	if len(a) != len(b) {
		panic("a and b lengths differ")
	}

	total := 0
	for i, av := range a {
		bv := b[i]
		total += aoc.Abs(av - bv)
	}

	return total
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 01)
	log.Printf("input solution: %d", solution("input", input))
}
