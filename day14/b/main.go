package main

import (
	"bytes"
	"github.com/asymmetricia/aoc23/coord"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

func tilt(w coord.World, d coord.Direction) (changed bool) {
	minx, miny, maxx, maxy := w.Rect()

	if d == coord.North {
		for y := miny + 1; y <= maxy; y++ {
			for x := minx; x <= maxx; x++ {
				from := coord.C(x, y)
				if w.At(from) != 'O' {
					continue
				}
				to := from
				for w.At(to.North()) == '.' {
					to = to.North()
				}
				if to != from {
					w.Set(to, 'O')
					w.Set(from, '.')
					changed = true
				}
			}
		}
	}
	if d == coord.East {
		for x := maxx - 1; x >= minx; x-- {
			for y := maxy; y >= miny; y-- {
				from := coord.C(x, y)
				if w.At(from) != 'O' {
					continue
				}
				to := from
				for w.At(to.East()) == '.' {
					to = to.East()
				}
				if to != from {
					w.Set(to, 'O')
					w.Set(from, '.')
					changed = true
				}
			}
		}
	}
	if d == coord.South {
		for y := maxy - 1; y >= miny; y-- {
			for x := minx; x <= maxx; x++ {
				from := coord.C(x, y)
				if w.At(from) != 'O' {
					continue
				}
				to := from
				for w.At(to.South()) == '.' {
					to = to.South()
				}
				if to != from {
					w.Set(to, 'O')
					w.Set(from, '.')
					changed = true
				}
			}
		}
	}
	if d == coord.West {
		for x := minx - 1; x <= maxx; x++ {
			for y := maxy; y >= miny; y-- {
				from := coord.C(x, y)
				if w.At(from) != 'O' {
					continue
				}
				to := from
				for w.At(to.West()) == '.' {
					to = to.West()
				}
				if to != from {
					w.Set(to, 'O')
					w.Set(from, '.')
					changed = true
				}
			}
		}
	}
	return changed
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

	w := coord.Load(lines, coord.LoadConfig{Dense: true})
	_, _, _, maxy := w.Rect()

	states := map[string]int{}
	var last time.Time
	for n := 0; n < 1000000000; n++ {
		if time.Since(last) > 5*time.Second {
			log.Printf("%d/%d (%.2f%%)", n, 1000000000, float32(n)/1000000000*100)
			last = time.Now()
		}
		changed := false
		changed = tilt(w, coord.North) || changed
		changed = tilt(w, coord.West) || changed
		changed = tilt(w, coord.South) || changed
		changed = tilt(w, coord.East) || changed
		if !changed {
			log.Print(n)
			break
		}
		if states != nil {
			s := w.(*coord.DenseWorld).String()
			if _, ok := states[s]; ok {
				log.Printf("%d .. %d", states[s], n)
				mod := n - states[s]
				for n+mod < 1000000000 {
					n += mod
				}
				states = nil
			} else {
				states[s] = n
			}
		}
	}

	var total int
	w.Each(func(c coord.Coord) (stop bool) {
		if w.At(c) == 'O' {
			total += maxy - c.Y + 1
		}
		return false
	})

	return total
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

	input := aoc.Input(2023, 14)
	log.Printf("input solution: %d", solution("input", input))
}
