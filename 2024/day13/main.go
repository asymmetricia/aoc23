package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
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

func canSolve(an, bn, pn int) (bool, int, int) {
	lcm, _ := aoc.LeastCommonMultiple(aoc.PrimeFactors(an), aoc.PrimeFactors(bn))
	bmax := pn / bn * bn
	for bi := 0; bi <= lcm/bn; bi++ {
		for ai := 0; ai <= lcm/an; ai++ {
			t := ai*an + bmax - bi*bn
			if t > pn {
				break
			}
			if t < pn {
				continue
			}
			return true, ai * an, bmax - bi*bn
		}
	}
	return false, 0, 0
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

	cost := new(int64)
	n := int64(len(machines))
	wg := &sync.WaitGroup{}
	for i, m := range machines {
		//wg.Wait()
		wg.Add(1)
		go func(i int, m Machine) {
			defer atomic.AddInt64(&n, -1)
			defer wg.Done()
			defer func() { log.Print(inc[0], i, atomic.LoadInt64(&n)) }()
			// solve a*AX+b*BX=PX and a*AY+b*BY=PY
			// a X+1 Y+2
			// b X+2 Y+3
			// a=(PX-b*BX)/AX

			// suppose PX=a*AX + b*BX
			// to reduce b by b', we have to increase a by a', such that a'*AX = b'*BX

			ok, ax, bx := canSolve(m.AX, m.BX, m.PX)
			if !ok {
				return
			}

			ok, ay, by := canSolve(m.AY, m.BY, m.PY)
			if !ok {
				return
			}

			lcmx, _ := aoc.LeastCommonMultiple(aoc.PrimeFactors(m.AX), aoc.PrimeFactors(m.BX))
			lcmy, _ := aoc.LeastCommonMultiple(aoc.PrimeFactors(m.AY), aoc.PrimeFactors(m.BY))
			// solve ax/m.AX + i * lcmx/m.AX = ay/m.AY + j * lcmy/m.AY

			// axp and bxp are starting points for a+b presses that solve X
			var axp = ax / m.AX
			var bxp = bx / m.BX

			// lcmpxa is the number of A presses that equate to lcmpxb B presses; so
			// axp+lcmpxa and bxp-lcmpxb will yield the same X value.
			var lcmpxa = lcmx / m.AX
			var lcmpxb = lcmx / m.BX

			// ayp and bxp are starting points for a+b pressed that solve Y
			var ayp = ay / m.AY
			var byp = by / m.BY
			// lcmpya is the number of Y presses that equate to lcmpyb B presses; so
			// axp+lcmpxa and bxp-lcmpxb will yield the same Y value
			var lcmpya = lcmy / m.AY
			var lcmpyb = lcmy / m.BY

			// solve the system:
			// axp+lcmpxa*n = ayp+lcmpya*m
			// bxp-lcmpxb*n = byp-lcmpyb*m

			for {
				for axp < ayp || bxp > byp {
					axp += lcmpxa
					bxp -= lcmpxb
				}
				for axp > ayp || bxp < byp {
					ayp += lcmpya
					byp -= lcmpyb
				}

				if bxp < 0 || byp < 0 {
					return
				}

				if axp == ayp && bxp == byp {
					log.Printf("%d: %d * %d + %d * %d = %d, want %d", i, axp, m.AX, (m.PX-axp*m.AX)/m.BX, m.BX, axp*m.AX+(m.PX-axp*m.AX)/m.BX*m.BX, m.PX)
					log.Printf("%d: %d * %d + %d * %d = %d, want %d", i, ayp, m.AY, (m.PY-ayp*m.AY)/m.BY, m.BY, ayp*m.AY+(m.PY-ayp*m.AY)/m.BY*m.BY, m.PY)
					atomic.AddInt64(cost, int64(3*axp+(m.PX-axp*m.AX)/m.BX))
					return
				}
			}
		}(i, m)
	}

	wg.Wait()
	return int(*cost)
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
