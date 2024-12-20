package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############
`

const testInputB = testInputA

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, 5},
	}

	for _, tt := range tests {
		t.Run(`2024-20 A `+tt.name, func(t *testing.T) {
			result := solution([]byte(tt.input), 2, 20)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestSolutionB(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic B", testInputB, 29},
	}

	for _, tt := range tests {
		t.Run(`2024-20 B `+tt.name, func(t *testing.T) {
			result := solution([]byte(tt.input), 20, 72)
			require.Equal(t, tt.expect, result)
		})
	}
}
