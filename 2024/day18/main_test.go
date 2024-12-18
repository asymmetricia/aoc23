package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0
`

const testInputB = testInputA

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, 22},
	}

	for _, tt := range tests {
		t.Run(`2024-18 A `+tt.name, func(t *testing.T) {
			result := solutionA([]byte(tt.input), true)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestSolutionB(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect string
	}

	tests := []test{
		{"basic B", testInputB, "6,1"},
	}

	for _, tt := range tests {
		t.Run(`2024-18 B `+tt.name, func(t *testing.T) {
			result := solutionB([]byte(tt.input), true)
			require.Equal(t, tt.expect, result)
		})
	}
}
