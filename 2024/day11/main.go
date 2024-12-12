package main

import (
	"bytes"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

var rule = aoc.Cache(ruleImpl)

func ruleImpl(stone int) []int {
	if stone == 0 {
		return []int{1}
	}
	if st := strconv.Itoa(stone); len(st)%2 == 0 {
		a, _ := strconv.ParseInt(st[:len(st)/2], 10, 64)
		b, _ := strconv.ParseInt(st[len(st)/2:], 10, 64)
		return []int{int(a), int(b)}
	}
	return []int{stone * 2024}
}

var ruleLen func(int, int) int

func init() {
	ruleLen = aoc.Cache2(ruleLenImpl)
}

func ruleLenImpl(stone int, i int) int {
	if i == 0 {
		return 1
	}

	var sum int
	for _, s := range rule(stone) {
		sum += ruleLen(s, i-1)
	}

	return sum
}

func solutionA(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	stones := aoc.Ints(lines[0])

	var sum int
	for _, stone := range stones {
		sum += ruleLen(stone, 25)
	}

	return sum
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	stones := aoc.Ints(lines[0])

	var sum int
	for _, stone := range stones {
		sum += ruleLen(stone, 75)
	}

	return sum
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 11)
	aStart := time.Now()
	a := solutionA(input)
	log.Printf("input solution A: %d, %dms", a, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	b := solutionB(input)
	log.Printf("input solution B: %d, %dms", b, time.Since(bStart).Milliseconds())
}
