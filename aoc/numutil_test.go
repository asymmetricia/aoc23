package aoc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSolveLinearSystem(t *testing.T) {
	type test struct {
		name    string
		coeff   [][]int
		answers []int
		want    []int
		wantErr error
	}
	var tests = []test{
		{"no integer solution", [][]int{{1, 2, 1}, {3, 2, 1}, {2, -3, 2}}, []int{1, 2, 3}, nil, NoIntegerSolution},

		// 94 x + 22 y = 8400
		// 34 x + 67 y = 5400
		{"2x2", [][]int{{94, 22}, {34, 67}}, []int{8400, 5400}, []int{80, 40}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SolveLinearSystem(tt.coeff, tt.answers)
			if tt.wantErr == nil {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestDeterminant(t *testing.T) {
	type test struct {
		name   string
		matrix [][]int
		want   int
	}

	tests := []test{
		{"2x2", [][]int{{2, 4}, {-1, 6}}, 16},
		{"3x3", [][]int{{1, 2, 3}, {4, 3, 2}, {3, 2, 1}}, 0},
		{"4x4", [][]int{{1, 2, 3, 4}, {5, 4, 2, 3}, {1, 3, 4, 5}, {3, 4, 3, 2}}, -12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Determinant(tt.matrix)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestMinor(t *testing.T) {
	type test struct {
		name   string
		matrix [][]int
		i, j   int
		want   int
	}
	tests := []test{
		{"3x3", [][]int{{1, 2, 3}, {4, 3, 2}, {3, 2, 1}}, 0, 0, -1},
		{"3x3", [][]int{{1, 2, 3}, {4, 3, 2}, {3, 2, 1}}, 1, 0, -4},
		{"3x3", [][]int{{1, 2, 3}, {4, 3, 2}, {3, 2, 1}}, 2, 0, -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Minor(tt.matrix, tt.i, tt.j)
			require.Equal(t, tt.want, got)
		})
	}
}
