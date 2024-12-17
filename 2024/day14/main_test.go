package main

import (
	"github.com/asymmetricia/aoc23/coord"
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3
`

const testInputB = testInputA

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, 12},
	}

	for _, tt := range tests {
		t.Run(`2024-14 A `+tt.name, func(t *testing.T) {
			result := solutionA([]byte(tt.input), true)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestRobot_AdvanceA(t *testing.T) {
	tests := []struct {
		name                     string
		position, velocity, want coord.Coord
	}{
		{"basic", coord.C(0, 0), coord.C(1, 1), coord.C(1, 1)},
		{"neg x", coord.C(0, 0), coord.C(-1, 1), coord.C(10, 1)},
		{"neg 2x", coord.C(0, 0), coord.C(-2, 1), coord.C(9, 1)},
		{"neg y", coord.C(0, 0), coord.C(-1, -1), coord.C(10, 6)},
		{"neg 2y", coord.C(0, 0), coord.C(-2, -2), coord.C(9, 5)},
		{"example 1", coord.C(2, 4), coord.C(2, -3), coord.C(4, 1)},
		{"example 2", coord.C(4, 1), coord.C(2, -3), coord.C(6, 5)},
		{"example 3", coord.C(6, 5), coord.C(2, -3), coord.C(8, 2)},
		{"example 4", coord.C(8, 2), coord.C(2, -3), coord.C(10, 6)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Robot{
				Position: tt.position,
				Velocity: tt.velocity,
			}
			r.AdvanceA(11, 7)
		})
	}
}
