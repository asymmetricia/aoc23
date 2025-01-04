package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

const testInputA = `1
10
100
2024
`

const testInputB = `1
2
3
2024
`

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, 37327623},
	}

	for _, tt := range tests {
		t.Run(`2024-22 A `+tt.name, func(t *testing.T) {
			result := solutionA([]byte(tt.input))
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
		{"basic B", testInputB, 23},
	}

	for _, tt := range tests {
		t.Run(`2024-22 B `+tt.name, func(t *testing.T) {
			result := solutionB([]byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}

func Test_evolve(t *testing.T) {
	secretNumber := 123
	seq := []int{
		15887950,
		16495136,
		527345,
		704524,
		1553684,
		12683156,
		11100544,
		12249484,
		7753432,
		5908254}
	for _, s := range seq {
		secretNumber = evolve(secretNumber)
		require.Equal(t, s, secretNumber)
	}
}

func Test_price(t *testing.T) {
	tests := []struct {
		secretNumber int
		deltas       [4]int8
		want         int8
	}{
		{1, [4]int8{-2, 1, -1, 3}, 7},
		{2, [4]int8{-2, 1, -1, 3}, 7},
		{3, [4]int8{-2, 1 - 1, 3}, 0},
		{2024, [4]int8{-2, 1, -1, 3}, 9},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d x %v", tt.secretNumber, tt.deltas), func(t *testing.T) {
			got := priceTree(tt.secretNumber, tt.deltas)
			require.Equal(t, tt.want, got)
		})
	}
}

func Benchmark_evolve(b *testing.B) {
	s := 1
	for i := 0; i < b.N; i++ {
		s = evolve(s)
	}
}

func Test_shifts(t *testing.T) {
	got := shifts(123)
	require.Equal(t, []int8{math.MinInt8, -3, 6, -1, -1, 0, 2, -2, 0, -2}, got[:10])
}
