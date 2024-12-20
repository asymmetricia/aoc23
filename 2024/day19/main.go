package main

import (
	"bytes"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solutionA(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	patterns := strings.Split(lines[0], ", ")
	re := regexp.MustCompile("^(" + strings.Join(patterns, "|") + ")+$")
	designs := lines[2:]

	var count int
	for _, design := range designs {
		if re.MatchString(design) {
			count++
		}
	}

	return count
}

var count func(design string, patterns string) int

func countImpl(design string, patterns string) int {
	var result int
	for _, pattern := range strings.Split(patterns, ",") {
		if design == pattern {
			result++
		} else if strings.HasPrefix(design, pattern) {
			result += count(strings.TrimPrefix(design, pattern), patterns)
		}
	}
	return result
}

func init() {
	count = aoc.Cache2(countImpl)
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	patterns := strings.Split(lines[0], ", ")
	designs := lines[2:]

	var result int
	for _, design := range designs {
		result += count(design, strings.Join(patterns, ","))
	}

	return result
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 19)
	aStart := time.Now()
	aSoln := solutionA(input)
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
