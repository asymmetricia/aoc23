package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `
`

const testInputB = testInputA

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, -1},
	}

	for _, tt := range tests {
		t.Run(`2024-03 A `+tt.name, func(t *testing.T) {
			result := solutionA(tt.name, []byte(tt.input))
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
		{"basic B", testInputB, -1},
	}

	for _, tt := range tests {
		t.Run(`2024-03 B `+tt.name, func(t *testing.T) {
			result := solutionB(tt.name, []byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
