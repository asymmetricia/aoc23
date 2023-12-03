package main

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

type Pull struct {
	Red, Green, Blue int
}

type Game struct {
	Id    int
	Pulls []Pull
}

var gameRe = regexp.MustCompile(`Game (\d+): (.*)`)

func ParseGame(in string) Game {
	res := gameRe.FindStringSubmatch(in)
	if res == nil {
		logrus.Fatalf("%q did not match RE %v", in, gameRe)
	}
	id, err := strconv.Atoi(res[1])
	if err != nil {
		logrus.Fatalf("game ID %q was not an integer: %v", res[1], err)
	}
	g := Game{Id: id}
	for _, pull := range strings.Split(res[2], ";") {
		g.Pulls = append(g.Pulls, ParsePull(pull))
	}
	return g
}

func ParsePull(pull string) Pull {
	var p Pull
	pull = strings.TrimSpace(pull)
	cubes := strings.Split(pull, ",")
	for _, cube := range cubes {
		cube = strings.TrimSpace(cube)
		numS, color := aoc.Split2(cube, " ")
		num := aoc.MustAtoi(numS)
		switch strings.ToLower(color) {
		case "red":
			p.Red = num
		case "green":
			p.Green = num
		case "blue":
			p.Blue = num
		}
	}
	return p
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	accum := 0
games:
	for _, line := range lines {
		g := ParseGame(line)
		for _, pull := range g.Pulls {
			if pull.Red > 12 || pull.Green > 13 || pull.Blue > 14 {
				continue games
			}
		}
		accum += g.Id
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

	input := aoc.Input(2023, 02)
	log.Printf("input solution: %d", solution("input", input))
}
