package main

import (
	"bytes"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

type Type int

const (
	HighCard Type = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

func HandType(a string) Type {
	hand := []rune(a)
	for i, card := range hand {
		if card == 'J' {
			newHand := slices.Clone(hand)
			best := HighCard
			for k := range values {
				if k == 'J' {
					continue
				}
				newHand[i] = rune(k)
				newHandType := HandType(string(newHand))
				if newHandType > best {
					best = newHandType
				}
			}
			return best
		}
	}
	sort.Slice(hand, func(i, j int) bool {
		return hand[i] < hand[j]
	})
	if hand[0] == hand[1] && hand[1] == hand[2] && hand[2] == hand[3] && hand[3] == hand[4] {
		return FiveOfAKind
	}

	if hand[0] == hand[1] && hand[1] == hand[2] && hand[2] == hand[3] ||
		hand[1] == hand[2] && hand[2] == hand[3] && hand[3] == hand[4] {
		return FourOfAKind
	}

	if hand[0] == hand[1] && hand[2] == hand[3] && hand[3] == hand[4] ||
		hand[0] == hand[1] && hand[1] == hand[2] && hand[3] == hand[4] {
		return FullHouse
	}

	if hand[0] == hand[1] && hand[1] == hand[2] ||
		hand[1] == hand[2] && hand[2] == hand[3] ||
		hand[2] == hand[3] && hand[3] == hand[4] {
		return ThreeOfAKind
	}
	if hand[0] == hand[1] && hand[2] == hand[3] ||
		hand[0] == hand[1] && hand[3] == hand[4] ||
		hand[1] == hand[2] && hand[3] == hand[4] {
		return TwoPair
	}
	if hand[0] == hand[1] || hand[1] == hand[2] || hand[2] == hand[3] || hand[3] == hand[4] {
		return OnePair
	}

	return HighCard
}

var values = map[uint8]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 1,
	'Q': 12,
	'K': 13,
	'A': 14}

func LessHand(a, b Hand) bool {
	ha, hb := HandType(a.Hand), HandType(b.Hand)
	if ha != hb {
		return ha < hb
	}

	for i := 0; i < 5; i++ {
		if a.Hand[i] != b.Hand[i] {
			return values[a.Hand[i]] < values[b.Hand[i]]
		}
	}

	return false
}

type Hand struct {
	Hand string
	Bid  int
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	var hands []Hand
	for _, line := range lines {
		hb := strings.Fields(line)
		hands = append(hands, Hand{hb[0], aoc.MustAtoi(hb[1])})
	}

	slices.SortFunc(hands, LessHand)

	var accum int = 0
	for i, h := range hands {
		accum += (i + 1) * h.Bid
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

	input := aoc.Input(2023, 7)
	log.Printf("input solution: %d", solution("input", input))
}
