package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))
`

const testInputB = `xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, 161},
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
		{"basic B", testInputB, 48},
	}

	for _, tt := range tests {
		t.Run(`2024-03 B `+tt.name, func(t *testing.T) {
			result := solutionB(tt.name, []byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
