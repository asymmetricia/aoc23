package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/coord"
	"github.com/asymmetricia/aoc23/isovox"
	"github.com/asymmetricia/aoc23/search"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func solutionA(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	grid := coord.Load(lines, coord.LoadConfig{Dense: true})

	zeroes := grid.Find('0')
	nines := grid.Find('9')

	var scores int

	for _, zero := range zeroes {
		for _, nine := range nines {
			path := search.AStar(
				zero,
				search.Goal(nine),
				search.Neighbors(func(a coord.Coord) []coord.Coord {
					var ret []coord.Coord
					av := grid.At(a)
					for _, neigh := range a.Neighbors(false) {
						nv := grid.At(neigh)
						if nv-av == 1 {
							ret = append(ret, neigh)
						}
					}
					return ret
				}))
			if path != nil {
				scores += 1
			}
		}
	}

	return scores
}

func rate(grid coord.World, from coord.Coord) []coord.Coord {
	fv := grid.At(from)
	if fv == '9' {
		return []coord.Coord{from}
	}
	var ret []coord.Coord
	for _, neigh := range from.Neighbors(false) {
		tv := grid.At(neigh)
		if tv-fv == 1 {
			ret = append(ret, rate(grid, neigh)...)
		}
	}
	return ret
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	grid := coord.Load(lines, coord.LoadConfig{Dense: true})

	zeroes := grid.Find('0')

	var scores int

	enc, err := aoc.NewMP4Encoder("../../2024-10.mp4", 60, log)
	if err != nil {
		log.Fatalf("could not initialize encoder: %v", err)
	}
	w := &isovox.World{map[isovox.Coord]*isovox.Voxel{}}
	for _, zero := range zeroes {
		grid.Each(func(c coord.Coord) (stop bool) {
			w.Voxels[isovox.Coord{c.X, c.Y, int(grid.At(c) - '0')}] = &isovox.Voxel{Color: aoc.TolScale('0', '9', grid.At(c))}
			return
		})

		w.Voxels[isovox.Coord{zero.X, zero.Y, 0}].Color = aoc.TolVibrantRed

		points := []coord.Coord{zero}
		for i := '1'; i <= '9'; i++ {
			var np []coord.Coord
			for _, point := range points {
				for _, neigh := range point.Neighbors(false) {
					if grid.At(neigh) == i {
						np = append(np, neigh)
						col := aoc.TolVibrantTeal
						if i == '9' {
							col = aoc.TolVibrantMagenta
						}
						w.Voxels[isovox.Coord{neigh.X, neigh.Y, int(i - '0')}] = &isovox.Voxel{Color: col}
					}
				}
			}
			points = np

			if err := enc.Encode(w.Render(10)); err != nil {
				log.Fatal(err)
			}

			if len(points) == 0 {
				break
			}
		}

		peaks := rate(grid, zero)
		scores += len(peaks)
	}

	if err := enc.Close(); err != nil {
		log.Fatal(err)
	}

	return scores
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 10)
	log.Printf("input solution A: %d", solutionA(input))
	log.Printf("input solution B: %d", solutionB(input))
}
