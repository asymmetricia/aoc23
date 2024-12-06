package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var nodes = map[string][2]string{}

var log = logrus.StandardLogger()

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	dirs := lines[0]
	for _, line := range lines[2:] {
		nodes[line[0:3]] = [2]string{line[7:10], line[12:15]}
	}

	cursor := "AAA"
	c := 0
	for cursor != "ZZZ" {
		fmt.Println(c)
		if dirs[c%len(dirs)] == 'L' {
			cursor = nodes[cursor][0]
		} else {
			cursor = nodes[cursor][1]
		}
		c++
	}

	return c
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
