package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `125 17
`

const testInputB = testInputA

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, 55312},
	}

	for _, tt := range tests {
		t.Run(`2024-11 A `+tt.name, func(t *testing.T) {
			result := solutionA([]byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
