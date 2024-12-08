package main

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solve(ints []int, want int, modeB bool) rune {
	if len(ints) == 2 {
		if ints[0]+ints[1] == want {
			return '+'
		}
		if ints[0]*ints[1] == want {
			return '*'
		}
		if modeB && aoc.Int(strconv.Itoa(ints[0])+strconv.Itoa(ints[1])) == want {
			return '|'
		}
		return 0
	} else {
		// a b c
		// a ? b + c
		// a ? b * c
		last := ints[len(ints)-1]
		if solve(ints[:len(ints)-1], want-last, modeB) != 0 {
			return '+'
		}
		if want%last == 0 && solve(ints[:len(ints)-1], want/last, modeB) != 0 {
			return '*'
		}

		wantS := strconv.Itoa(want)
		lastS := strconv.Itoa(last)
		if modeB && strings.HasSuffix(wantS, lastS) && wantS != lastS &&
			solve(ints[:len(ints)-1], aoc.Int(strings.TrimSuffix(wantS, lastS)), modeB) != 0 {
			return '|'
		}
		return 0
	}
}

func solutionA(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d lines (%d unique)", len(lines), len(uniq))

	total := 0
	for _, line := range lines {
		want, operands := aoc.Split2(line, ": ")
		wantI, operandsI := aoc.Int(want), aoc.Ints(operands)
		if solve(operandsI, wantI, false) != 0 {
			total += wantI
		}
	}

	return total
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d lines (%d unique)", len(lines), len(uniq))

	total := 0
	for _, line := range lines {
		want, operands := aoc.Split2(line, ": ")
		wantI, operandsI := aoc.Int(want), aoc.Ints(operands)
		if solve(operandsI, wantI, true) != 0 {
			total += wantI
		}
	}

	return total
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 07)
	log.Printf("input solution A: %d", solutionA(input))
	log.Printf("input solution B: %d", solutionB(input))
}
