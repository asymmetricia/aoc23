package aoc

// Count counts the number of times the function `f` returns `true` for the slice
// `s`
func Count[T any](s []T, f func(t T) bool) int {
	c := 0
	for _, v := range s {
		if f(v) {
			c++
		}
	}
	return c
}

// CountEq counts the number of times the value `b` appears in the slice `s`
func CountEq[T comparable](s []T, b T) int {
	c := 0
	for _, v := range s {
		if v == b {
			c++
		}
	}
	return c
}
