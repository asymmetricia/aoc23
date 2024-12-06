package main

import (
	"bytes"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solutionA(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	split := slices.Index(lines, "")
	rules := aoc.Map(lines[:split], func(x string) [2]int {
		a, b := aoc.Split2(x, "|")
		return [2]int{aoc.Int(a), aoc.Int(b)}
	})
	updates := aoc.Map(lines[split+1:], func(x string) []int {
		return aoc.Map(strings.Split(x, ","), aoc.Int)
	})

	sum := 0
	for _, update := range updates {
		if !slices.IsSortedFunc(update, func(a, b int) bool {
			for _, rule := range rules {
				if rule[0] == a && rule[1] == b {
					return true
				}
			}
			return false
		}) {
			continue
		}

		sum += update[len(update)/2]
	}
	return sum
}

func solutionB(name string, input []byte) int {
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	split := slices.Index(lines, "")
	rules := aoc.Map(lines[:split], func(x string) [2]int {
		a, b := aoc.Split2(x, "|")
		return [2]int{aoc.Int(a), aoc.Int(b)}
	})
	updates := aoc.Map(lines[split+1:], func(x string) []int {
		return aoc.Map(strings.Split(x, ","), aoc.Int)
	})

	for _, update := range updates {
		for _, i := range update {
			if slices.IndexFunc(rules, func(ints [2]int) bool {
				return ints[0] == i || ints[1] == i
			}) == -1 {
				panic("update missing from rules: " + strconv.Itoa(i))
			}
		}
	}

	fn := func(a, b int) bool {
		for _, rule := range rules {
			if rule[0] == a && rule[1] == b {
				return true
			}
		}
		return false
	}

	sum := 0
	for _, update := range updates {
		if slices.IsSortedFunc(update, fn) {
			continue
		}
		slices.SortStableFunc(update, fn)
		sum += update[len(update)/2]
	}

	return sum
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 05)
	log.Printf("input solution A: %d", solutionA("input", input))
	log.Printf("input solution B: %d", solutionB("input", input))
}
