package main

import (
	"bytes"
	"math/big"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

type Node string

func NewNode(s string) Node {
	return Node(s)
}

var log = logrus.StandardLogger()

func bf(nodes map[Node][2]Node, steps []int, cursors []Node) int {
	c := 0
	stop := false
	for !stop {
		stop = true
		for i := range cursors {
			cursors[i] = nodes[cursors[i]][steps[c%len(steps)]]
			stop = stop && cursors[i][2] == 'Z'
		}
		c++
	}
	return c
}

func solve(nodes map[Node][2]Node, steps []int, cursors []Node) *big.Int {
	ret := big.NewInt(1)
	retFactors := map[int]uint{}
	for i, cursor := range cursors {
		c := 0
		for {
			cMod := c % len(steps)
			cursor = nodes[cursor][steps[cMod]]
			c++
			if cursor[2] == 'Z' {
				break
			}
		}
		next := c
		for {
			nextMod := next % len(steps)
			cursor = nodes[cursor][steps[nextMod]]
			next++
			if cursor[2] == 'Z' {
				break
			}
		}
		factors := aoc.PrimeFactors(next - c)
		_, retFactors = aoc.LeastCommonMultiple(factors, retFactors)
		log.Print(i, " ", c, " ", next, " ", next-c, " ", factors, " ", retFactors)
	}

	for f, n := range retFactors {
		f := big.NewInt(int64(f))
		for i := uint(0); i < n; i++ {
			ret.Mul(ret, f)
			log.Print(f, " ", ret)
		}
	}

	// nope: 16683196289967
	return ret
}

func solution(name string, input []byte) *big.Int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	var steps []int
	var nodes = map[Node][2]Node{}
	for _, c := range lines[0] {
		if c == 'L' {
			steps = append(steps, 0)
		} else {
			steps = append(steps, 1)
		}
	}
	var starts []Node
	for _, line := range lines[2:] {
		n := NewNode(line[0:3])
		if line[2] == 'A' {
			starts = append(starts, n)
		}
		nodes[n] = [2]Node{NewNode(line[7:10]), NewNode(line[12:15])}
	}

	log.Print(bf(nodes, steps, starts[0:2]))
	return solve(nodes, steps, starts)
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

	input := aoc.Input(2023, 8)
	log.Printf("input solution: %d", solution("input", input))
}
