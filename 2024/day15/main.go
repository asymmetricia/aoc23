package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/canvas"
	"github.com/asymmetricia/aoc23/coord"
	"github.com/asymmetricia/aoc23/isovox"
	"golang.org/x/exp/slices"
	"strings"
	"time"
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

	sep := slices.Index(lines, "")
	grid := coord.Load(lines[:sep], coord.LoadConfig{Dense: true})
	steps := strings.Join(lines[sep+1:], "")

	robot := grid.Find('@')[0]

steps:
	for _, step := range steps {
		dir := map[rune]coord.Direction{
			'^': coord.North,
			'>': coord.East,
			'v': coord.South,
			'<': coord.West,
		}[step]
		cursor := robot
		var push []coord.Coord

		for {
			cursor = cursor.Move(dir)
			v := grid.At(cursor)
			if v == '#' {
				continue steps
			}
			if v == 'O' {
				push = append(push, cursor)
				continue
			}
			break
		}

		for _, p := range push {
			grid.Set(p.Move(dir), 'O')
		}
		grid.Set(robot, '.')
		robot = robot.Move(dir)
		grid.Set(robot, '@')
	}

	var ret int
	for _, box := range grid.Find('O') {
		ret += 100*box.Y + box.X
	}

	return ret
}

func canPush(grid coord.World, obj coord.Coord, dir coord.Direction) bool {
	v := grid.At(obj)
	if v == '#' {
		return false
	}
	if v == '.' || v == 0 {
		return true
	}

	switch dir {
	case coord.East, coord.West:
		if v == '@' || v == '[' || v == ']' {
			return canPush(grid, obj.Move(dir), dir)
		}
	case coord.North, coord.South:
		if v == '@' {
			return canPush(grid, obj.Move(dir), dir)
		}
		if v == '[' {
			return canPush(grid, obj.Move(dir), dir) && canPush(grid, obj.Move(dir).East(), dir)
		}
		if v == ']' {
			return canPush(grid, obj.Move(dir), dir) && canPush(grid, obj.Move(dir).West(), dir)
		}
	}

	panic(obj)
}

func push(grid coord.World, obj coord.Coord, dir coord.Direction) {
	if grid.At(obj) == '.' || grid.At(obj) == 0 {
		return
	}
	switch dir {
	case coord.North, coord.South:
		if grid.At(obj) == '[' {
			push(grid, obj.Move(dir), dir)
			push(grid, obj.Move(dir).East(), dir)
			grid.Set(obj.Move(dir), '[')
			grid.Set(obj.Move(dir).East(), ']')
			grid.Set(obj, '.')
			grid.Set(obj.East(), '.')
			return
		}
		if grid.At(obj) == ']' {
			push(grid, obj.Move(dir).West(), dir)
			push(grid, obj.Move(dir), dir)
			grid.Set(obj.Move(dir).West(), '[')
			grid.Set(obj.Move(dir), ']')
			grid.Set(obj.West(), '.')
			grid.Set(obj, '.')
			return
		}
		fallthrough
	case coord.East, coord.West:
		push(grid, obj.Move(dir), dir)
		grid.Set(obj.Move(dir), grid.At(obj))
		grid.Set(obj, '.')
		return
	}
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	sep := slices.Index(lines, "")
	sgrid := coord.Load(lines[:sep], coord.LoadConfig{Dense: true})

	grid := &coord.DenseWorld{}
	sgrid.Each(func(c coord.Coord) (stop bool) {
		switch sgrid.At(c) {
		case '#':
			grid.Set(coord.C(c.X*2, c.Y), '#')
			grid.Set(coord.C(c.X*2+1, c.Y), '#')
		case 'O':
			grid.Set(coord.C(c.X*2, c.Y), '[')
			grid.Set(coord.C(c.X*2+1, c.Y), ']')
		case '@':
			grid.Set(coord.C(c.X*2, c.Y), '@')
		}
		return
	})
	grid.Each(func(c coord.Coord) (stop bool) {
		if grid.At(c) == 0 {
			grid.Set(c, '.')
		}
		return
	})
	steps := strings.Join(lines[sep+1:], "")

	robot := grid.Find('@')[0]

	enc, err := aoc.NewMP4Encoder("../../2024-15.mp4", 60, log)
	if err != nil {
		log.Fatal(err)
	}

	var stack []*canvas.Canvas
	last := time.Now()
	for stepIdx, step := range steps {
		if time.Since(last) >= time.Second {
			log.Printf("%d/%d", stepIdx, len(steps))
			last = time.Now()
		}
		dir := map[rune]coord.Direction{
			'^': coord.North,
			'>': coord.East,
			'v': coord.South,
			'<': coord.West,
		}[step]
		if !canPush(grid, robot, dir) {
			continue
		}
		push(grid, robot, dir)
		robot = robot.Move(dir)

		if stepIdx%5 > 0 {
			continue
		}

		cv := &canvas.Canvas{}
		iw := isovox.World{Voxels: map[isovox.Coord]*isovox.Voxel{}}
		grid.Each(func(c coord.Coord) (stop bool) {
			if grid.At(c) == '@' {
				cv.Set(c.X, c.Y, canvas.Cell{Color: aoc.TolVibrantMagenta, Value: '@'})
				iw.Voxels[isovox.Coord{c.X, c.Y, 0}] = &isovox.Voxel{Color: aoc.TolVibrantMagenta, Size: 4}
			}
			if grid.At(c) == ']' || grid.At(c) == '[' {
				cv.Set(c.X, c.Y, canvas.Cell{Color: aoc.TolVibrantCyan, Value: grid.At(c)})
				iw.Voxels[isovox.Coord{c.X, c.Y, 0}] = &isovox.Voxel{Color: aoc.TolVibrantCyan, Size: 6}
			}
			if grid.At(c) == '#' {
				cv.Set(c.X, c.Y, canvas.Cell{Color: aoc.TolVibrantGrey, Value: '#'})
				iw.Voxels[isovox.Coord{c.X, c.Y, 0}] = &isovox.Voxel{Color: aoc.TolVibrantGrey}
			}
			return
		})
		stack = append(stack, cv)
		if err := enc.Encode(iw.Render(8)); err != nil {
			log.Fatal(err)
		}
	}
	if err := enc.Close(); err != nil {
		log.Fatal(err)
	}

	canvas.RenderGif(stack, "../../2024-15.gif", log)

	var ret int
	for _, box := range grid.Find('[') {
		ret += 100*box.Y + box.X
	}

	return ret
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 15)
	aStart := time.Now()
	aSoln := solutionA(input)
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
