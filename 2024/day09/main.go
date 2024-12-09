package main

import (
	"bytes"
	"fmt"
	"github.com/asymmetricia/aoc23/canvas"
	"github.com/asymmetricia/aoc23/coord"
	"golang.org/x/exp/slices"
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
	uniq := map[string]bool{}
	for _, line := range lines {
		uniq[line] = true
	}
	log.Printf("read %d lines (%d unique)", len(lines), len(uniq))

	var disk []int
	for idx, i := range lines[0] {
		i -= '0'
		for i > 0 {
			if idx%2 == 1 {
				disk = append(disk, -1)
			} else {
				disk = append(disk, idx/2)
			}
			i--
		}
	}
	for j := len(disk) - 1; j >= 0; j-- {
		if disk[j] != -1 {
			mv := disk[j]
			disk[j] = -1
			for i, v := range disk {
				if v == -1 {
					disk[i] = mv
					break
				}
			}
		}
	}

	checksum := 0
	for i, v := range disk {
		if v > -1 {
			checksum += i * v
		}
	}

	return checksum
}

type block struct {
	id   int
	len  int
	free bool
}

func render(disk []block, reading, writing, percent int) *canvas.Canvas {
	const width = 120
	const height = 28
	const blockSize = 32
	cv := &canvas.Canvas{Timing: 1}
	canvas.TextBox{
		Top:         0,
		Left:        0,
		Title:       []rune(" Elfen™ Defragmenter "),
		TitleColor:  aoc.TolVibrantMagenta,
		Footer:      []rune(" !!! UNREGISTERED SHAREWARE !!! "),
		FooterColor: aoc.TolVibrantRed,
		FrameColor:  aoc.TolVibrantBlue,
		Width:       width - 2,
		Height:      height - 2,
	}.On(cv)

	canvas.TextBox{
		Top:        29,
		Left:       10,
		Title:      []rune(" Status "),
		TitleColor: aoc.TolVibrantMagenta,
		FrameColor: aoc.TolVibrantBlue,
		Width:      22,
		Height:     4,
	}.On(cv)
	cv.PrintAt(28, 30, fmt.Sprintf("% 3d%%", percent), aoc.TolVibrantGrey)
	cv.PrintAt(12, 31, string(canvas.ProgressBar(percent, 20)), aoc.TolVibrantGrey)
	cv.PrintAt(12, 32, "Full Defragmentation", aoc.TolVibrantGrey)

	legendLeft := width - 35
	canvas.TextBox{
		Top:        29,
		Left:       legendLeft,
		Title:      []rune(" Legend "),
		TitleColor: aoc.TolVibrantMagenta,
		FrameColor: aoc.TolVibrantBlue,
		Width:      22,
		Height:     4,
	}.On(cv)
	cv.PrintAt(legendLeft+2, 30, string([]rune{aoc.BlockFull, aoc.BlockDark, aoc.BlockMedium, aoc.BlockLight}), aoc.TolVibrantRed)
	cv.PrintAt(legendLeft+6, 30, " - Writing", aoc.TolVibrantGrey)
	cv.PrintAt(legendLeft+2, 31, string([]rune{aoc.BlockFull, aoc.BlockDark, aoc.BlockMedium, aoc.BlockLight}), aoc.TolVibrantCyan)
	cv.PrintAt(legendLeft+6, 31, " - Reading", aoc.TolVibrantGrey)
	cv.PrintAt(legendLeft+2, 33, fmt.Sprintf("   %c = %d blocks", aoc.BlockFull, blockSize), aoc.TolVibrantGrey)

	cursor := coord.C(1, 1)
	var accum []int
	var accumType int
	for segIdx, seg := range disk {
		for i := 0; i < seg.len; i++ {
			if seg.free {
				accum = append(accum, 0)
			} else {
				accum = append(accum, 1)
			}

			if segIdx == writing {
				accumType = 1
			} else if segIdx == reading {
				accumType = 2
			}

			if len(accum) == blockSize {
				sum := aoc.Sum(accum)
				v := aoc.BlockLight
				if sum > blockSize*3/4 {
					v = aoc.BlockFull
				} else if sum > blockSize*1/2 {
					v = aoc.BlockDark
				} else if sum > blockSize*1/4 {
					v = aoc.BlockMedium
				}
				col := aoc.TolVibrantGrey
				if accumType == 1 {
					col = aoc.TolVibrantRed
				} else if accumType == 2 {
					col = aoc.TolVibrantCyan
				}

				cv.Set(cursor.X, cursor.Y, canvas.Cell{Color: col, Value: v})
				cursor.X++
				if cursor.X >= width-1 {
					cursor.X = 1
					cursor.Y++
				}
				accum = accum[0:0]
				accumType = 0
			}
		}
	}
	return cv
}

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	var disk []block
	for i, v := range lines[0] {
		if i%2 == 0 {
			disk = append(disk, block{i / 2, int(v - '0'), false})
		} else {
			disk = append(disk, block{0, int(v - '0'), true})
		}
	}

	var stack []*canvas.Canvas
	for i := len(disk) - 1; i >= 0; i-- {
		movedTo := -1
		if disk[i].free {
			continue
		}
		toMove := disk[i]
		for j := 0; j < i; j++ {
			if !disk[j].free {
				continue
			}
			if disk[j].len < toMove.len {
				continue
			}

			movedTo = j

			disk = slices.Insert(disk, j, toMove)
			disk[j+1].len -= toMove.len
			if disk[i].free {
				disk[i].len += toMove.len
			} else if i+2 < len(disk) && disk[i+2].free {
				disk[i+2].len += toMove.len
			} else {
				disk[i+1].free = true
			}
			disk = slices.Delete(disk, i+1, i+2)
			break
		}
		stack = append(stack, render(disk, i, movedTo, 100-i*100/len(disk)))
	}

	checksum := 0
	blockId := 0
	for _, seg := range disk {
		for i := 0; i < seg.len; i++ {
			if !seg.free {
				fmt.Print(seg.id)
				checksum += blockId * seg.id
			} else {
				fmt.Print(".")
			}
			blockId++
		}
	}

	if err := canvas.RenderMp4(stack, "../../2024-09.mp4", log); err != nil {
		log.Fatal(err)
	}

	return checksum
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 9)
	log.Printf("input solution A: %d", solutionA(input))
	log.Printf("input solution B: %d", solutionB(input))
}
