package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInput = `7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9
`

func TestSolution(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic", testInput, 4},
	}

	for _, tt := range tests {
		t.Run(`2024-02 b `+tt.name, func(t *testing.T) {
			result := solution(tt.name, []byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
