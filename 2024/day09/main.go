package main

import (
	"bytes"
	"fmt"
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

func solutionB(input []byte) int {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	type block struct {
		id   int
		len  int
		free bool
	}

	var disk []block
	for i, v := range lines[0] {
		if i%2 == 0 {
			disk = append(disk, block{i / 2, int(v - '0'), false})
		} else {
			disk = append(disk, block{0, int(v - '0'), true})
		}
	}

	for i := len(disk) - 1; i >= 0; i-- {
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
