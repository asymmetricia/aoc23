package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

var log = logrus.StandardLogger()

var aReg = regexp.MustCompile(`Button A: X([+-]\d+), Y([+-]\d+)`)
var bReg = regexp.MustCompile(`Button B: X([+-]\d+), Y([+-]\d+)`)
var pReg = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

type Machine struct {
	AX, AY, BX, BY int
	PX, PY         int
}

func solutionA(input []byte) int {
	return solutionB(input, false)
}

func solutionB(input []byte, inc ...bool) int {
	if len(inc) == 0 {
		inc = []bool{true}
	}
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")
	var machines []Machine
	for i := 0; i < len(lines); i += 4 {
		aMatch := aReg.FindStringSubmatch(lines[i])
		bMatch := bReg.FindStringSubmatch(lines[i+1])
		pMatch := pReg.FindStringSubmatch(lines[i+2])
		if aMatch == nil || bMatch == nil || pMatch == nil {
			panic(strings.Join(lines[i:i+3], "\n"))
		}
		m := Machine{
			AX: aoc.Int(aMatch[1]),
			AY: aoc.Int(aMatch[2]),
			BX: aoc.Int(bMatch[1]),
			BY: aoc.Int(bMatch[2]),
			PX: aoc.Int(pMatch[1]),
			PY: aoc.Int(pMatch[2]),
		}
		if inc[0] {
			m.PX += 10000000000000
			m.PY += 10000000000000
		}
		machines = append(machines, m)
	}

	var cost int
	for _, m := range machines {
		soln, _ := aoc.SolveLinearSystem(
			[][]int{
				{m.AX, m.BX},
				{m.AY, m.BY},
			}, []int{
				m.PX, m.PY,
			})
		if soln != nil {
			cost += soln[0]*3 + soln[1]
		}
	}

	return cost
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 13)
	aStart := time.Now()
	aSoln := solutionA(input)
	if aSoln != 25629 {
		panic(fmt.Sprintf("bad solution for part A %d", aSoln))
	}
	log.Printf("input solution A: %d (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
