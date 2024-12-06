package main

import (
	"bytes"
	"math"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

type MapEntry struct {
	From, To, Length int
}

func (m MapEntry) Match(from int) bool {
	return from >= m.From && from < m.From+m.Length
}

func (m MapEntry) Convert(from int) (to int) {
	return m.To + from - m.From
}

func (m MapEntry) Reverse(to int) (from int) {
	return m.From + to - m.To
}

type Map []MapEntry

func (m Map) Convert(from int) (to int) {
	for _, e := range m {
		if e.Match(from) {
			return e.Convert(from)
		}
	}
	return from
}

func solution(name string, input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	var SeedToSoil, SoilToFert, FertToWater, WaterToLight, LightToTemp, TempToHumid, HumidToLoc Map

	var target *Map
	var seeds []int

	for _, line := range lines {
		if strings.HasPrefix(line, "seeds: ") {
			for _, s := range strings.Fields(aoc.After(line, ": ")) {
				seeds = append(seeds, aoc.MustAtoi(s))
			}
			continue
		}
		switch line {
		case "seed-to-soil map:":
			target = &SeedToSoil
		case "soil-to-fertilizer map:":
			target = &SoilToFert
		case "fertilizer-to-water map:":
			target = &FertToWater
		case "water-to-light map:":
			target = &WaterToLight
		case "light-to-temperature map:":
			target = &LightToTemp
		case "temperature-to-humidity map:":
			target = &TempToHumid
		case "humidity-to-location map:":
			target = &HumidToLoc
		default:
			conv := strings.Fields(line)
			if len(conv) == 3 {
				*target = append(*target, MapEntry{
					To:     aoc.MustAtoi(conv[0]),
					From:   aoc.MustAtoi(conv[1]),
					Length: aoc.MustAtoi(conv[2]),
				})
			}
		}
	}

	best := math.MaxInt
	for _, seed := range seeds {
		soil := SeedToSoil.Convert(seed)
		fert := SoilToFert.Convert(soil)
		water := FertToWater.Convert(fert)
		light := WaterToLight.Convert(water)
		temp := LightToTemp.Convert(light)
		humid := TempToHumid.Convert(temp)
		loc := HumidToLoc.Convert(humid)
		log.Printf("%d -> %d -> %d -> %d -> %d -> %d -> %d -> %d", seed, soil, fert, water, light, temp, humid, loc)
		if loc < best {
			best = loc
		}
	}

	return best
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

	input := aoc.Input(2023, 5)
	log.Printf("input solution: %d", solution("input", input))
}
