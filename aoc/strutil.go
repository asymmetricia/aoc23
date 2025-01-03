package aoc

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strconv"
	"strings"
)

func Before(haystack, needle string) string {
	pos := strings.Index(haystack, needle)
	if pos < 0 {
		return haystack
	}
	return haystack[:pos]
}

func After(haystack, needle string) string {
	pos := strings.Index(haystack, needle)
	if pos < 0 {
		return ""
	}
	return haystack[pos+len(needle):]
}

func Split2(haystack, needle string) (string, string) {
	pos := strings.Index(haystack, needle)
	if pos < 0 {
		return haystack, ""
	}
	return haystack[:pos], haystack[pos+len(needle):]
}

func Int(in string) int {
	i, err := strconv.Atoi(in)
	if err != nil {
		panic(fmt.Sprintf("could not parse %q as string: %v", in, err))
	}
	return i
}

func Ints(in string) []int {
	return Map(strings.Fields(in), Int)
}

func Filter[X any](x []X, f func(X) bool) []X {
	var ret []X
	for _, xx := range x {
		if f(xx) {
			ret = append(ret, xx)
		}
	}
	return ret
}

func Map[X any, Y any](x []X, f func(X) Y) []Y {
	ret := make([]Y, len(x))
	for i, xx := range x {
		ret[i] = f(xx)
	}
	return ret
}

func Sum[X constraints.Integer | constraints.Float](values []X) X {
	var sum X
	for _, value := range values {
		sum += value
	}
	return sum
}
