package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const TestInputA = `Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0
`

func TestSolutionA(t *testing.T) {
	require.Equal(t, "4,6,3,5,6,3,5,2,1,0", solutionA([]byte(TestInputA)))
}

func TestComputer(t *testing.T) {
	type test struct {
		name                string
		a, b, c             int64
		program             string
		wantA, wantB, wantC int64
		wantOutput          string
	}

	tests := []test{
		{"bst", 0, 0, 9, "2,6", -1, 1, -1, ""},
		{"basic 2", 10, 0, 0, "5,0,5,1,5,4", -1, -1, -1, "0,1,2"},
		{"basic 3", 2024, 0, 0, "0,1,5,4,3,0", 0, -1, -1, "4,2,5,6,7,7,7,7,3,1,0"},
		{"basic 4", 0, 29, 0, "1,7", 0, 26, 0, ""},
		{"basic 5", 0, 2024, 43690, "4,0", 0, 44354, -1, ""}}

	for _, tt := range tests {
		t.Run(`2024-17 A `+tt.name, func(t *testing.T) {
			c := NewComputer(tt.program)
			c.A, c.B, c.C = tt.a, tt.b, tt.c
			c.Run()
			if tt.wantA != -1 {
				require.Equal(t, tt.wantA, c.A)
			}
			if tt.wantB != -1 {
				require.Equal(t, tt.wantB, c.B)
			}
			if tt.wantC != -1 {
				require.Equal(t, tt.wantC, c.C)
			}
			require.Equal(t, tt.wantOutput, c.Output)
		})
	}
}

const testInputB = `Register A: 2024
Register B: 0
Register C: 0

Program: 0,3,5,4,3,0`

func TestSolutionB(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int64
	}

	tests := []test{
		{"basic B", testInputB, 117440},
	}

	for _, tt := range tests {
		t.Run(`2024-17 B `+tt.name, func(t *testing.T) {
			result := solutionB([]byte(tt.input), true)
			require.Equal(t, tt.expect, result)
		})
	}
}
