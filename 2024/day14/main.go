package main

import (
	"bytes"
	"fmt"
	"github.com/asymmetricia/aoc23/canvas"
	"github.com/asymmetricia/aoc23/coord"
	"image/png"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

type Robot struct {
	Position coord.Coord
	Velocity coord.Coord
}

func (r *Robot) AdvanceA(width, height int) {
	r.Position.X += r.Velocity.X
	// width == 7
	// 0..6
	// 0-1 == 6
	for r.Position.X < 0 {
		r.Position.X += width
	}
	// 6+1 == 0
	for r.Position.X >= width {
		r.Position.X -= width
	}

	r.Position.Y += r.Velocity.Y
	// 1-3 == 5
	// 0-2 == 5
	// 6-1 == 5
	// 5   == 5
	for r.Position.Y < 0 {
		r.Position.Y += height
	}
	for r.Position.Y >= height {
		r.Position.Y -= height
	}
}

var robotRe = regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)

func solutionA(input []byte, test bool) int {
	width := 101
	height := 103

	if test {
		width = 11
		height = 7
	}

	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	var robots []*Robot
	for _, line := range lines {
		m := robotRe.FindStringSubmatch(line)
		robots = append(robots, &Robot{
			Position: coord.Coord{aoc.Int(m[1]), aoc.Int(m[2])},
			Velocity: coord.Coord{aoc.Int(m[3]), aoc.Int(m[4])},
		})
	}

	for step := 0; step < 100; step++ {
		for _, robot := range robots {
			robot.AdvanceA(width, height)
		}
	}

	var grid coord.World = &coord.DenseWorld{}
	for _, robot := range robots {
		if c := grid.At(robot.Position); c <= 0 {
			grid.Set(robot.Position, '1')
		} else {
			grid.Set(robot.Position, c+1)
		}
	}

	var quads [4]int
	for _, robot := range robots {
		// 11/2 == 5 ... [0,5)
		if robot.Position.X < width/2 {
			// 7/2 == 3 ... [0,3)
			if robot.Position.Y < height/2 {
				quads[0]++
			}
			// (3,6]
			if robot.Position.Y > height/2 {
				quads[1]++
			}
		}
		// (5,10]
		if robot.Position.X > width/2 {
			if robot.Position.Y < height/2 {
				quads[2]++
			}
			if robot.Position.Y > height/2 {
				quads[3]++
			}
		}
	}

	return quads[0] * quads[1] * quads[2] * quads[3]
}

func countCluster(world coord.World, start coord.Coord, visited map[coord.Coord]bool) int {
	visited[start] = true
	var ret int
	for _, n := range start.Neighbors(false) {
		v := world.At(n)
		if !visited[n] && v >= '1' && v <= '9' {
			ret = ret + 1 + countCluster(world, n, visited)
		}
	}
	return ret
}

func solutionB(input []byte, test bool) int {

	width := 101
	height := 103

	if test {
		width = 11
		height = 7
	}

	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	var robots []*Robot
	for _, line := range lines {
		m := robotRe.FindStringSubmatch(line)
		robots = append(robots, &Robot{
			Position: coord.Coord{aoc.Int(m[1]), aoc.Int(m[2])},
			Velocity: coord.Coord{aoc.Int(m[3]), aoc.Int(m[4])},
		})
	}

	var grid coord.World = &coord.DenseWorld{}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			grid.Set(coord.C(x, y), ' ')
		}
	}

	last := time.Now()
	for step := 0; step < 10000; step++ {
		if time.Since(last) > time.Second {
			log.Print(step)
			last = time.Now()
		}
		for _, robot := range robots {
			v := grid.At(robot.Position)
			if v >= '2' && v <= '9' {
				grid.Set(robot.Position, v-1)
			} else {
				grid.Set(robot.Position, ' ')
			}
			robot.AdvanceA(width, height)
			v = grid.At(robot.Position)
			if v >= '1' && v <= '8' {
				grid.Set(robot.Position, v+1)
			} else if v != '9' {
				grid.Set(robot.Position, '1')
			}
		}

		render := false
		visited := map[coord.Coord]bool{}
		grid.Each(func(c coord.Coord) (stop bool) {
			if visited[c] || grid.At(c) == ' ' {
				return
			}

			if countCluster(grid, c, visited) > 30 {
				render = true
				return true
			}

			return
		})

		if render {
			cv := &canvas.Canvas{}
			grid.Each(func(c coord.Coord) (stop bool) {
				cv.Set(c.X, c.Y, canvas.Cell{Color: aoc.TolVibrantGrey, Value: grid.At(c)})
				return
			})
			f, err := os.OpenFile(fmt.Sprintf("../../2024-12-%d.png", step+1), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				log.Fatal(err)
			}
			png.Encode(f, cv.Render())
			f.Sync()
			f.Close()
		}
	}

	return -1
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 14)
	aStart := time.Now()
	aSoln := solutionA(input, false)
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input, false)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
