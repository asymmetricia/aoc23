package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
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
		{"basic A", testInputA, 7036},
	}

	for _, tt := range tests {
		t.Run(`2024-16 A `+tt.name, func(t *testing.T) {
			result := solutionA([]byte(tt.input))
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
		{"basic B", testInputB, 45},
	}

	for _, tt := range tests {
		t.Run(`2024-16 B `+tt.name, func(t *testing.T) {
			result := solutionB([]byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
