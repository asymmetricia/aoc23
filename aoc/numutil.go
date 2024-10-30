package aoc

import (
	"golang.org/x/exp/constraints"
	"math"
)

func MaxFn[T any, K constraints.Ordered](a []T, fn func(T) K) (ret K) {
	if len(a) == 0 {
		return ret
	}

	best := fn(a[0])
	for _, item := range a[1:] {
		i := fn(item)
		if best < i {
			best = i
		}
	}
	return best
}

func Max[K constraints.Ordered](a ...K) K {
	first := true
	var max K
	for _, a := range a {
		if first || a > max {
			max = a
			first = false
		}
	}
	return max
}

func Min[K constraints.Ordered](a ...K) K {
	first := true
	var min K
	for _, a := range a {
		if first || a < min {
			min = a
			first = false
		}
	}
	return min
}

func Abs[K constraints.Signed | constraints.Float](a K) K {
	if a < 0 {
		return -a
	}
	return a
}

func PrimeFactors[T constraints.Integer](i T) map[T]uint {
	ret := map[T]uint{}
	var term T = T(math.Ceil(math.Sqrt(float64(i))))
	var f T = 2
	for {
		for i%f == 0 {
			i /= f
			ret[f]++
		}
		f++
		if i == 1 {
			break
		}
		if f > term {
			ret[i]++
			break
		}
	}
	return ret
}

func LeastCommonMultiple[T constraints.Integer](a map[T]uint, b map[T]uint) (T, map[T]uint) {
	factors := map[T]uint{}
	for f := range a {
		factors[f] = Max(a[f], b[f])
	}
	for f := range b {
		factors[f] = Max(a[f], b[f])
	}

	var ret T = 1
	for f, n := range factors {
		for i := uint(0); i < n; i++ {
			ret *= f
		}
	}
	return ret, factors
}
