package main

import (
	"bytes"
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func mix[T int](a, b T) T {
	return a ^ b
}

func prune[T int](a T) T {
	return a & (1<<24 - 1)
}

func evolve(secretNumber int) int {
	secretNumber = prune(mix(secretNumber<<6, secretNumber))
	secretNumber = prune(mix(secretNumber>>5, secretNumber))
	return prune(mix(secretNumber<<11, secretNumber))
}

type node struct {
	value  int8
	leaves map[int8]*node
}

func priceTree(isn int, deltas [4]int8) int8 {
	root := tree(isn)
	if l1, ok := root.leaves[deltas[0]]; ok {
		if l2, ok := l1.leaves[deltas[1]]; ok {
			if l3, ok := l2.leaves[deltas[2]]; ok {
				if l4, ok := l3.leaves[deltas[3]]; ok {
					return l4.value
				}
			}
		}
	}
	return 0
}

var tree = aoc.Cache(func(isn int) *node {
	root := &node{
		value:  math.MinInt8,
		leaves: make(map[int8]*node),
	}

	ones := ones(isn)
	shifts := shifts(isn)
	for i := 1; i+3 < len(shifts); i++ {
		l1, ok := root.leaves[shifts[i]]
		if !ok {
			l1 = &node{
				value:  0,
				leaves: map[int8]*node{},
			}
			root.leaves[shifts[i]] = l1
		}

		l2, ok := l1.leaves[shifts[i+1]]
		if !ok {
			l2 = &node{
				value:  0,
				leaves: map[int8]*node{},
			}
			l1.leaves[shifts[i+1]] = l2
		}

		l3, ok := l2.leaves[shifts[i+2]]
		if !ok {
			l3 = &node{
				value:  0,
				leaves: map[int8]*node{},
			}
			l2.leaves[shifts[i+2]] = l3
		}

		_, ok = l3.leaves[shifts[i+3]]
		if !ok {
			l3.leaves[shifts[i+3]] = &node{
				value: ones[i+3],
			}
		}
	}

	return root
})

var ones = aoc.Cache(func(isn int) []int8 {
	ret := []int8{int8(isn % 10)}
	for i := 1; i < 2001; i++ {
		isn = evolve(isn)
		ret = append(ret, int8(isn%10))
	}
	return ret
})

var shifts = aoc.Cache(func(isn int) []int8 {
	ones := ones(isn)
	shifts := make([]int8, 2001)
	shifts[0] = math.MinInt8
	for i := 1; i < 2001; i++ {
		shifts[i] = ones[i] - ones[i-1]
	}
	return shifts
})

func solutionA(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	isns := aoc.Ints(strings.Join(lines, " "))

	var value int
	for _, isn := range isns {
		for i := 0; i < 2000; i++ {
			isn = evolve(isn)
		}
		value += isn
	}

	return value
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	isns := aoc.Ints(strings.Join(lines, " "))

	var value int

	// initialize the hypothetical previous sequence
	deltas := [4]int8{-9, -9, -9, -10}
	for {
		// compute the next sequence
		deltas[3]++
		if deltas[3] == 10 {
			deltas[2]++
			deltas[3] = -9
		}
		if deltas[2] == 10 {
			deltas[1]++
			deltas[2] = -9
		}
		if deltas[1] == 10 {
			deltas[0]++
			deltas[1] = -9
			log.Print(deltas)
		}
		if deltas[0] == 10 {
			break
		}

		// reject any sequences that can't exist
		var valid bool
	checks:
		for check := int8(-9); check <= 9; check++ {
			check := check
			for _, d := range deltas {
				check += d
				if check > 9 || check < 0 {
					continue checks
				}
			}
			valid = true
			break
		}
		if !valid {
			continue
		}

		var candidate int
		for _, isn := range isns {
			candidate += int(priceTree(isn, deltas))
		}
		if candidate > value {
			value = candidate
			log.Print(deltas, value)
		}
	}

	return value
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 22)
	aStart := time.Now()
	aSoln := solutionA(input)
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
