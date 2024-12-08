package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............
`

const testInputB = testInputA

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, 14},
	}

	for _, tt := range tests {
		t.Run(`2024-08 A `+tt.name, func(t *testing.T) {
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
		{"basic B", testInputB, 34},
	}

	for _, tt := range tests {
		t.Run(`2024-08 B `+tt.name, func(t *testing.T) {
			result := solutionB([]byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
