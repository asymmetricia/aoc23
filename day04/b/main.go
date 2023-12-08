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

type Card struct {
	Want, Have set.Set[string]
	Count      int64
}

var log = logrus.StandardLogger()

func solution(name string, input []byte) int64 {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	var cards []*Card
	for _, card := range lines {
		haveL, wantL := aoc.Split2(aoc.After(card, ": "), " | ")
		have := set.FromWords(haveL)
		want := set.FromWords(wantL)
		cards = append(cards, &Card{want, have, 1})
	}

	var total int64
	for id, card := range cards {
		n := len(card.Have.Intersect(card.Want))
		logrus.Printf("card %d wins %d of next %d cards", id, card.Count, n)
		for i := id + 1; i < id+1+n; i++ {
			cards[i].Count += card.Count
		}
		total += card.Count
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

	input := aoc.Input(2023, 04)
	log.Printf("input solution: %d", solution("input", input))
}
