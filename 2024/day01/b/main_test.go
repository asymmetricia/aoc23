package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInput = `3   4
4   3
2   5
1   3
3   9
3   3
`

func TestSolution(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic", testInput, 31},
	}

	for _, tt := range tests {
		t.Run(`2024-01 a `+tt.name, func(t *testing.T) {
			result := solution(tt.name, []byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
