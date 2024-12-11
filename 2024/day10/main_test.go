package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732
`

const testInputB = testInputA

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, 36},
	}

	for _, tt := range tests {
		t.Run(`2024-10 A `+tt.name, func(t *testing.T) {
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
		{"basic B", testInputB, 81},
	}

	for _, tt := range tests {
		t.Run(`2024-10 B `+tt.name, func(t *testing.T) {
			result := solutionB([]byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
