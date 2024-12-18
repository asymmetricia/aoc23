package main

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/asymmetricia/aoc23/aoc"
)

type Computer struct {
	Program []int
	IP      int
	A, B, C int64
	Output  string
}

func (c Computer) String() string {
	return fmt.Sprintf("IP=%d; A:%d, B:%d, C:%d; Output:%s", c.IP, c.A, c.B, c.C, c.Output)
}

func NewComputer(program string) *Computer {
	return &Computer{
		Program: aoc.Map(strings.Split(program, ","), func(x string) int {
			return aoc.Int(x)
		}),
	}
}

func (c Computer) Combo(operand int64) int64 {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return c.A
	case 5:
		return c.B
	case 6:
		return c.C
	default:
		panic(fmt.Sprintf("invalid operand %d", operand))
	}
}

type OpCode uint8

var (
	OpAdv OpCode = 0
	OpBxl OpCode = 1
	OpBst OpCode = 2
	OpJnz OpCode = 3
	OpBxc OpCode = 4
	OpOut OpCode = 5
	OpBdv OpCode = 6
	OpCdv OpCode = 7
)

var Operations = map[OpCode]func(c *Computer, operand int64){
	OpAdv: func(c *Computer, operand int64) {
		c.A = c.A >> c.Combo(operand)
	},
	OpBxl: func(c *Computer, operand int64) {
		c.B = c.B ^ operand
	},
	OpBst: func(c *Computer, operand int64) {
		c.B = c.Combo(operand) & 0b111
	},
	OpJnz: func(c *Computer, operand int64) {
		if c.A == 0 {
			return
		}
		c.IP = int(operand) - 2
	},
	OpBxc: func(c *Computer, operand int64) {
		c.B = c.B ^ c.C
	},
	OpOut: func(c *Computer, operand int64) {
		v := c.Combo(operand) & 0b111
		if c.Output != "" {
			c.Output = fmt.Sprintf("%s,%d", c.Output, v)
		} else {
			c.Output = fmt.Sprintf("%d", v)
		}
	},
	OpBdv: func(c *Computer, operand int64) {
		c.B = c.A >> c.Combo(operand)
		for operand > 0 {
			c.B /= 2
			operand--
		}
	},
	OpCdv: func(c *Computer, operand int64) {
		c.C = c.A >> c.Combo(operand)
	},
}

func (c *Computer) Run() string {
	for c.IP < len(c.Program) && c.IP >= 0 {
		Operations[OpCode(c.Program[c.IP])](c, int64(c.Program[c.IP+1]))
		c.IP += 2
	}
	return c.Output
}

var log = logrus.StandardLogger()

func solutionA(input []byte) string {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	c := NewComputer(aoc.After(lines[4], ": "))
	c.A, _ = strconv.ParseInt(aoc.After(lines[0], ": "), 10, 64)
	c.B, _ = strconv.ParseInt(aoc.After(lines[1], ": "), 10, 64)
	c.C, _ = strconv.ParseInt(aoc.After(lines[2], ": "), 10, 64)
	log.Print(c)
	c.Run()
	log.Print(c)
	return c.Output
}

func search(c *Computer, values []int, idx int, program string, target func(int) int) []int {
	t := target(idx)
	if t < 0 {
		return values
	}
	if t >= len(program) {
		panic(fmt.Sprintf("target %d >= len(program) %d", t, len(program)))
	}
	for i := 0; i < 8; i++ {
		c.IP = 0
		c.B = 0
		c.C = 0
		c.A = 0
		c.Output = ""

		values[idx] = i
		for _, v := range values {
			c.A = c.A<<3 | int64(v)
		}
		c.Run()

		if len(c.Output) > t && c.Output[t] == program[t] {
			if idx == len(values)-1 {
				return values
			}

			if rest := search(c, slices.Clone(values), idx+1, program, target); rest != nil {
				return rest
			}
		}
	}
	return nil
}

func solutionB(input []byte, test bool) int64 {
	// trim trailing space only
	input = bytes.Replace(input, []byte("\r"), []byte(""), -1)
	input = bytes.TrimRightFunc(input, unicode.IsSpace)
	lines := strings.Split(strings.TrimRightFunc(string(input), unicode.IsSpace), "\n")

	prog := aoc.After(lines[4], ": ")

	c := NewComputer(prog)
	c.A, _ = strconv.ParseInt(aoc.After(lines[0], ": "), 10, 64)
	c.B, _ = strconv.ParseInt(aoc.After(lines[1], ": "), 10, 64)
	c.C, _ = strconv.ParseInt(aoc.After(lines[2], ": "), 10, 64)

	/*
		BST 4 // $B = $A % 0b111
		BXL 2 // $B = $B ^ 0b010
		CDV 5 // $C = $A / 2^$B
		BXC 3 // $B = $B ^ $C
		ADV 3 // $A = $A / 0b100
		BXL 7 // $B = $B ^ 0b111
		OUT 5 // output $B
		JNZ 0 // GOTO 0
	*/

	values := make([]int, len(c.Program))
	if test {
		values = search(c, values, 0, prog, func(i int) int {
			return 8 - i*2
		})
	} else {
		values = search(c, values, 0, prog, func(i int) int {
			return len(prog) - i*2 - 1
		})
	}

	var ret int64
	for _, v := range values {
		ret = ret<<3 | int64(v)
	}

	return ret
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	input := aoc.Input(2024, 17)
	aStart := time.Now()
	aSoln := solutionA(input)
	log.Printf("input solution A: %s (%dms)", aSoln, time.Since(aStart).Milliseconds())

	bStart := time.Now()
	bSoln := solutionB(input, false)
	log.Printf("input solution B: %d (%dms)", bSoln, time.Since(bStart).Milliseconds())
}
