package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/aoc"
	"github.com/asymmetricia/aoc23/coord"
	"golang.org/x/exp/slices"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func expand(world coord.World) {
	minx, miny, maxx, maxy := world.Rect()
cols:
	for x := minx; x <= maxx; x++ {
		for y := miny; y <= maxy; y++ {
			if world.At(coord.C(x, y)) == '#' {
				continue cols
			}
		}
		galaxies := world.Find('#')
		slices.SortFunc(galaxies, func(a, b coord.Coord) bool {
			return b.X < a.X
		})
		for _, c := range galaxies {
			if c.X <= x {
				continue
			}
			world.Set(c.East(), '#')
			world.Set(c, '>')
		}
		maxx++
		x++
	}

rows:
	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			if world.At(coord.C(x, y)) == '#' {
				continue rows
			}
		}
		galaxies := world.Find('#')
		slices.SortFunc(galaxies, func(a, b coord.Coord) bool {
			return b.Y < a.Y
		})
		for _, c := range galaxies {
			if c.Y <= y {
				continue
			}
			world.Set(c.South(), '#')
			world.Set(c, 'v')
		}
		maxy++
		y++
	}
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

	world := coord.Load(lines, coord.LoadConfig{Ignore: "."})

	world.Print()
	expand(world)
	world.Print()

	var sum int
	galaxies := world.Find('#')
	for i := 0; i < len(galaxies); i++ {
		iC := galaxies[i]
		for j := i + 1; j < len(galaxies); j++ {
			jC := galaxies[j]
			//log.Printf("%s -> %s == %d", iC, jC, iC.TaxiDistance(jC))
			sum += iC.TaxiDistance(jC)
		}
	}

	return sum
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

	input := aoc.Input(2023, 11)
	log.Printf("input solution: %d", solution("input", input))
}
