package main

import (
	"github.com/asymmetricia/aoc23/coord"
	"testing"
)

func Test_countMirroredRows(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  int
	}{
		{
			"simple",
			[]string{
				"##......#",
				"##......#",
			},
			1,
		},
		{
			"simple two",
			[]string{
				"..#.##.#.",
				"##......#",
				"##......#",
				"..#.##.#.",
			},
			2,
		},
		{
			"asym zero",
			[]string{
				"#.##..##.",
				"..#.##.#.",
				"##......#",
				"##......#",
				"..#.##.#.",
				"..##..##.",
				"#.#.##.#.",
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countMirroredRows(*(coord.Load(tt.lines, coord.LoadConfig{Dense: true}).(*coord.DenseWorld)))
			if got != tt.want {
				t.Errorf("countMirroredRows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_countMirroredCols(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  int
	}{
		{
			"asym three",
			[]string{
				"#.##..##.",
				"..#.##.#.",
				"##......#",
				"##......#",
				"..#.##.#.",
				"..##..##.",
				"#.#.##.#.",
			},
			5,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countMirroredCols(*(coord.Load(tt.lines, coord.LoadConfig{Dense: true}).(*coord.DenseWorld)))
			if got != tt.want {
				t.Errorf("countMirroredCols() = %v, want %v", got, tt.want)
			}
		})
	}
}
