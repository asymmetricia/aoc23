package main

import (
	"bytes"
	"encoding/json"
	"regexp"
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

	total := 0

	mulRe := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	matches := mulRe.FindAllStringSubmatch(string(input), -1)
	for _, mul := range matches {
		a := aoc.Int(mul[1])
		b := aoc.Int(mul[2])
		total += a * b
	}

	return total
}

func solutionB(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)

	mulRe := regexp.MustCompile(`(mul)\((\d+),(\d+)\)|(do)\(\)|(don't)\(\)`)
	matches := mulRe.FindAllStringSubmatch(string(input), -1)
	j, _ := json.Marshal(matches)
	println(string(j))

	total := 0
	enabled := true
	for _, mul := range matches {
		if strings.HasPrefix(mul[0], "do(") {
			enabled = true
		} else if strings.HasPrefix(mul[0], "don't(") {
			enabled = false
		} else if strings.HasPrefix(mul[0], "mul(") {
			if enabled {
				a := aoc.Int(mul[2])
				b := aoc.Int(mul[3])
				total += a * b
			}
		} else {
			panic("bad op:" + mul[0])
		}
	}

	return total
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 03)
	log.Printf("input solution A: %d", solutionA("input", input))
	log.Printf("input solution B: %d", solutionB("input", input))
}
