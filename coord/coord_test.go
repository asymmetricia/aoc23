package coord_test

import (
	"github.com/asymmetricia/aoc23/coord"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCoord_Unit(t *testing.T) {
	tests := []struct {
		c    coord.Coord
		want coord.Coord
	}{
		{coord.Coord{-2, 4}, coord.Coord{-1, 2}},
		{coord.Coord{2, -4}, coord.Coord{1, -2}},
		{coord.Coord{4, -8}, coord.Coord{1, -2}},
	}
	for _, tt := range tests {
		t.Run(tt.c.String(), func(t *testing.T) {
			require.Equal(t, tt.want, tt.c.Unit())
		})
	}
}
