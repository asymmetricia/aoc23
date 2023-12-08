package aoc

import "strconv"

func MustAtoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return i
}

func MustAtoiSlice(a []string) []int {
	var ret []int
	for _, s := range a {
		ret = append(ret, MustAtoi(s))
	}
	return ret
}
