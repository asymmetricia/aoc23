package main

import (
	"testing"
)

//func Test_check(t *testing.T) {
//	tests := []struct {
//		name      string
//		record    string
//		groupings []byte
//		valid     bool
//		partial   bool
//	}{
//		{"valid", "#.#.###", []byte{1, 1, 3}, true, false},
//		{"partial", "?.#.###", []byte{1, 1, 3}, false, true},
//		{"partial 2", ".###????????", []byte{3, 2, 1}, false, true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			valid, partial := check(tt.record, tt.groupings)
//			require.Equal(t, tt.valid, valid, "validity")
//			require.Equal(t, tt.partial, partial, "partial")
//		})
//	}
//}

func Test_reconstruct(t *testing.T) {
	tests := []struct {
		name      string
		record    []string
		groupings []byte
		want      int
	}{
		{"one.1", []string{"?"}, []byte{1}, 1},
		{"one.2", []string{"?", "?"}, []byte{1, 1}, 1},
		{"two", []string{"??", "?"}, []byte{1, 1}, 2},
		{"one.3", []string{"#???"}, []byte{1}, 1},
		{"one.4", []string{"###", "##", "#???"}, []byte{3, 2, 1}, 1},
		{"four", []string{"###", "##", "????"}, []byte{3, 2, 1}, 4},
		{"ten", []string{"###", "???????"}, []byte{3, 2, 1}, 10},
		{"ten", []string{"###????????"}, []byte{3, 2, 1}, 10},
		{"ten", []string{"?###????????"}, []byte{3, 2, 1}, 10},
		{"zero", []string{"####????????"}, []byte{3, 2, 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := map[string]map[string]int{}
			if got := reconstruct(cache, tt.record, tt.groupings); got != tt.want {
				t.Errorf("reconstruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
