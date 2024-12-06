package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSolution(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"empty", "", -1},
	}

	for _, tt := range tests {
		t.Run(`2024-01 a `+tt.name, func(t *testing.T) {
			result := solution(tt.name, []byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
