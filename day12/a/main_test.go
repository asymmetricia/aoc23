package main

import "testing"

func Test_check(t *testing.T) {
	tests := []struct {
		name      string
		record    string
		groupings []int
		want      bool
	}{
		{"valid", "#.#.###", []int{1, 1, 3}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := check(tt.record, tt.groupings); got != tt.want {
				t.Errorf("check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reconstruct(t *testing.T) {
	tests := []struct {
		name      string
		record    string
		groupings []int
		want      int
	}{
		{"ten", "?###????????", []int{3, 2, 1}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reconstruct(tt.record, tt.groupings); got != tt.want {
				t.Errorf("reconstruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
