package main

import (
	"bytes"
	"os"
	"strings"
	"unicode"

	"github.com/asymmetricia/aoc23/set"
	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	accum := 0
	for _, card := range lines {
		haveL, wantL := aoc.Split2(aoc.After(card, ": "), " | ")
		have := set.FromWords(haveL)
		want := set.FromWords(wantL)
		is := have.Intersect(want)
		if len(is) == 0 {
			continue
		}
		v := 1
		for i := 0; i < len(is)-1; i++ {
			v *= 2
		}
		accum += v
	}

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

	input := aoc.Input(2023, 4)
	log.Printf("input solution: %d", solution("input", input))
}
