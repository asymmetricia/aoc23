package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"math/big"
	"os"
	"strings"
	"unicode"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func reconstruct(cache map[string]map[string]int, parts []string, groups []byte) int {
	ck := strings.Join(parts, "-")
	if m, ok := cache[ck]; ok {
		if v, ok := m[string(groups)]; ok {
			return v
		}
	} else {
		cache[ck] = map[string]int{}
	}

	for i, p := range parts {
		if unk := strings.IndexRune(p, '?'); unk >= 0 {
			var a []string
			if unk > 0 {
				a = []string{p[:unk]}
			}
			if p[unk+1:] != "" {
				a = append(a, p[unk+1:])
			}
			ar := reconstruct(cache, append(a, parts[i+1:]...), groups[i:])
			b := []string{p[:unk] + "#" + p[unk+1:]}
			br := reconstruct(cache, append(b, parts[i+1:]...), groups[i:])

			cache[ck][string(groups)] = ar + br
			return ar + br
		}

		if i >= len(groups) || byte(len(p)) != groups[i] {
			cache[ck][string(groups)] = 0
			return 0
		}
	}

	if len(parts) == len(groups) {
		cache[ck][string(groups)] = 1
		return 1
	}

	cache[ck][string(groups)] = 0
	return 0
}

func solution(name string, input []byte) *big.Int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d %s lines (%d unique)", len(lines), name, len(uniq))

	cache := map[string]map[string]int{}
	total := &big.Int{}
	for _, line := range lines {
		fields := strings.Fields(line)
		record := fields[0]
		var groupings []byte
		for _, f := range strings.Split(fields[1], ",") {
			groupings = append(groupings, byte(aoc.MustAtoi(f)))
		}
		base := slices.Clone(groupings)
		for i := 0; i < 4; i++ {
			record += "?" + fields[0]
			groupings = append(groupings, base...)
		}

		parts := strings.Split(record, ".")
		for i := 0; i < len(parts); i++ {
			if len(parts[i]) == 0 {
				parts = append(parts[:i], parts[i+1:]...)
				i--
			}
		}
		total.Add(total, big.NewInt(int64(reconstruct(cache, parts, groupings))))
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
