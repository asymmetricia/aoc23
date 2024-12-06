package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInput = `
`

func TestSolution(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic", testInput, -1},
	}

	for _, tt := range tests {
		t.Run(`2023-11 a ` + tt.name, func(t *testing.T) {
			result := solution(tt.name, []byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
