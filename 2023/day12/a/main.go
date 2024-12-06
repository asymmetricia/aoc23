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

func check(record string, groupings []int) bool {
	for _, g := range groupings {
		// consume leading `.`
		for len(record) > 0 && record[0] == '.' {
			record = record[1:]
			if len(record) == 0 {
				return false
			}
		}

		// consume the exact number of `#`
		for g > 0 {
			if len(record) == 0 || record[0] == '.' {
				return false
			}
			g--
			record = record[1:]
		}

		if len(record) > 0 && record[0] == '#' {
			return false
		}
	}

	for len(record) > 0 && record[0] == '.' {
		record = record[1:]
	}

	if len(record) > 0 {
		return false
	}

	return true
}

func reconstruct(record string, groupings []int) int {
	if strings.ContainsRune(record, '?') {
		a := strings.Replace(record, "?", "#", 1)
		b := strings.Replace(record, "?", ".", 1)
		ar := reconstruct(a, groupings)
		br := reconstruct(b, groupings)
		return ar + br
	}

	if check(record, groupings) {
		return 1
	}

	return 0
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

	total := 0
	for _, line := range lines {
		fields := strings.Fields(line)
		total += reconstruct(fields[0], aoc.MustAtoiSlice(strings.Split(fields[1], ",")))
	}

	return total
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

	input := aoc.Input(2023, 12)
	log.Printf("input solution: %d", solution("input", input))
}
