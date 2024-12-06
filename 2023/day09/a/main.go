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

func Diffs(in []int) []int {
	var ret []int
	for i := 0; i < len(in)-1; i++ {
		ret = append(ret, in[i+1]-in[i])
	}
	return ret
}

func Zero(in []int) bool {
	for _, i := range in {
		if i != 0 {
			return false
		}
	}
	return true
}

func Next(in []int) int {
	diffs := Diffs(in)
	if Zero(diffs) {
		return in[0]
	}

	return in[len(in)-1] + Next(diffs)
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d %s lines (%d unique)", len(lines), name, len(uniq))

	sum := 0
	for _, line := range lines {
		fields := strings.Fields(line)
		var fieldsN []int
		for _, f := range fields {
			fieldsN = append(fieldsN, aoc.Int(f))
		}
		sum += Next(fieldsN)
	}

	return sum
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

	input := aoc.Input(2023, 9)
	log.Printf("input solution: %d", solution("input", input))
}
