package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const testInputA = `47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47
`

const testInputB = testInputA

func TestSolutionA(t *testing.T) {
	type test struct {
		name   string
		input  string
		expect int
	}

	tests := []test{
		{"basic A", testInputA, 143},
	}

	for _, tt := range tests {
		t.Run(`2024-05 A `+tt.name, func(t *testing.T) {
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
		{"basic B", testInputB, 123},
	}

	for _, tt := range tests {
		t.Run(`2024-05 B `+tt.name, func(t *testing.T) {
			result := solutionB(tt.name, []byte(tt.input))
			require.Equal(t, tt.expect, result)
		})
	}
}
