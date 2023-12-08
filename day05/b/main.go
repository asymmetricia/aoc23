package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

// dest, src, length

var log = logrus.StandardLogger()

type Range struct {
	Start, Length int
}

func (r Range) String() string {
	return fmt.Sprintf("[%d,%d]", r.Start, r.Start+r.Length-1)
}

type MapEntry struct {
	From, To, Length int
}

type Map []MapEntry

func (m MapEntry) Convert(from Range) (to []Range, remainder []Range) {
	if from.Start+from.Length < m.From || m.From+m.Length <= from.Start {
		return nil, []Range{from}
	}

	// from (1 len 10)
	// entry (5->15 len 5)
	// -> rem (1 len 4)
	//    from (5 len 6)
	if from.Start < m.From {
		rem := Range{from.Start, m.From - from.Start}
		remainder = append(remainder, rem)

		from.Length -= rem.Length
		from.Start = m.From
	}

	// rem (1 len 4)
	// from (5 len 6) -- 5 6 7 8 9 10
	// entry (5->15 len 5) -- 5 6 7 8 9
	// -> rem (1 len 4) (10 len 1)
	//    from (5 len 5)
	if from.Start+from.Length > m.From+m.Length {
		// we end after the conversion range ends
		rem := Range{from.Start + m.Length, from.Length - m.Length}
		remainder = append(remainder, rem)
		from.Length = m.Length
	}

	return []Range{{m.To + from.Start - m.From, aoc.Min(from.Length, m.Length)}}, remainder
}

func (m Map) Convert(from ...Range) (to []Range) {
	// if entries are in our input range, some parts of them become part of the output range

	// run the from range through the first conversion
	// collect the result, and run each remaining range through the second conversion
	// etc.

	for _, e := range m {
		var rems []Range
		for _, f := range from {
			result, split := e.Convert(f)
			to = append(to, result...)
			rems = append(rems, split...)
		}
		from = rems
	}

	return append(to, from...)
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
	for i := 0; i < len(seeds); i += 2 {
		seed := Range{seeds[i], seeds[i+1]}
		soil := SeedToSoil.Convert(seed)
		fert := SoilToFert.Convert(soil...)
		water := FertToWater.Convert(fert...)
		light := WaterToLight.Convert(water...)
		temp := LightToTemp.Convert(light...)
		humid := TempToHumid.Convert(temp...)
		loc := HumidToLoc.Convert(humid...)
		log.Printf("%v -> %v -> %v -> %v -> %v -> %v -> %v -> %v", seed, soil, fert, water, light, temp, humid, loc)
		for _, loc := range loc {
			if loc.Start < best {
				best = loc.Start
			}
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

	input := aoc.Input(2023, 05)
	log.Printf("input solution: %d", solution("input", input))
}
