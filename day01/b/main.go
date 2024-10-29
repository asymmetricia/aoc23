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

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	nums := map[string]rune{
		"zero":  '0',
		"one":   '1',
		"two":   '2',
		"three": '3',
		"four":  '4',
		"five":  '5',
		"six":   '6',
		"seven": '7',
		"eight": '8',
		"nine":  '9',
	}
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	accum := 0
	for _, line := range lines {
		first := -1
		last := -1
		for i, c := range line {
			if c < '0' || c > '9' {
				for k, v := range nums {
					if len(line[i:]) >= len(k) && line[i:i+len(k)] == k {
						c = v
					}
				}
			}
			if c >= '0' && c <= '9' {
				if first == -1 {
					first = int(c - '0')
					last = first
				} else {
					last = int(c - '0')
				}
			}

		}
		log.Printf("%s = %d", line, first*10+last)
		accum += first*10 + last
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

	input := aoc.Input(2023, 1)
	log.Printf("input solution: %d", solution("input", input))
}
