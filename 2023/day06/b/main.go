package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solution(name string, input []byte) uint64 {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	time := aoc.MustAtoi(strings.ReplaceAll(aoc.After(lines[0], ":"), " ", ""))
	dist := aoc.MustAtoi(strings.ReplaceAll(aoc.After(lines[1], ":"), " ", ""))
	var count uint64
	for i := 1; i < time; i++ {
		if i*(time-i) > dist {
			count++
		}
	}

	return count
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

	input := aoc.Input(2023, 6)
	log.Printf("input solution: %d", solution("input", input))
}
