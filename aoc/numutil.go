package aoc

import (
	"errors"
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

// PrimeFactors returns a map mapping each prime factor to its power in composing
// the given value.
func PrimeFactors[T constraints.Integer](i T) map[T]uint {
	if i == 0 {
		return nil
	}

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

var NoSolution = errors.New("no solution")
var NonSquareMatrix = errors.New("non-square matrix")
var WrongAnswerCount = errors.New("wrong answer count")
var NoIntegerSolution = errors.New("no integer solution")

// SolveLinearSystemFloat solves a system of linear equations using cramer's
// rule. coeff should be a square nxn matrix, and answers should have n entries.
// Floating point conversion may render solutions inaccurate. For guaranteed
// accurate integer solutions, see SolveLinearSystem.
//
// # Example
//
//	// 94 x + 22 y = 8400
//	// 34 x + 67 y = 5400
//	solns, err := SolveLinearSystemFloat([][]float32{{94,22}, {34,67}}, []float32{8400,5400})
//	// solns is []float32{80,40}
//	// err is nil
//	// x=80
//	// y=40
func SolveLinearSystemFloat[T constraints.Float](coeff [][]T, answers []T) ([]T, error) {
	n := len(coeff)
	if len(answers) != n {
		return nil, WrongAnswerCount
	}

	for _, coeffs := range coeff {
		if len(coeffs) != n {
			return nil, NonSquareMatrix
		}
	}

	solutions := make([]T, n)
	detA := Determinant(coeff)
	if detA == 0 {
		return nil, NoSolution
	}

	for i := 0; i < n; i++ {
		an := make([][]T, n)
		for row := range an {
			an[row] = make([]T, n)
			for col := range an[row] {
				if col == i {
					an[row][col] = answers[row]
				} else {
					an[row][col] = coeff[row][col]
				}
			}
		}
		detAn := Determinant(an)
		solutions[i] = detAn / detA
	}
	return solutions, nil
}

// SolveLinearSystem solves a system of linear equations using cramer's rule.
// coeff should be a square nxn matrix, and answers should have n entries. There
// may be floating point solutions for the system that this does not return,
// because it only returns valid integer solutions. See SolveLinearSystemFloat.
func SolveLinearSystem[T constraints.Signed](coeff [][]T, answers []T) ([]T, error) {
	n := len(coeff)
	if len(answers) != n {
		return nil, WrongAnswerCount
	}

	for _, coeffs := range coeff {
		if len(coeffs) != n {
			return nil, NonSquareMatrix
		}
	}

	solutions := make([]T, n)
	detA := Determinant(coeff)
	if detA == 0 {
		return nil, NoSolution
	}

	for i := 0; i < n; i++ {
		an := make([][]T, n)
		for row := range an {
			an[row] = make([]T, n)
			for col := range an[row] {
				if col == i {
					an[row][col] = answers[row]
				} else {
					an[row][col] = coeff[row][col]
				}
			}
		}
		detAn := Determinant(an)
		if detAn%detA != 0 {
			// no integer solution !!!
			return nil, NoIntegerSolution
		}
		solutions[i] = detAn / detA
	}
	return solutions, nil
}

func Determinant[T constraints.Signed | constraints.Float](matrix [][]T) T {
	n := len(matrix)
	if n == 2 {
		return matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
	}

	var result T
	for i := 0; i < n; i++ {
		var sign T = 1
		if (i)%2 == 1 {
			sign = -1
		}
		result += sign * Minor(matrix, i, 0) * matrix[i][0]
	}
	return result
}

func Minor[T constraints.Signed | constraints.Float](matrix [][]T, i, j int) T {
	n := len(matrix)
	minorMatrix := make([][]T, n-1)
	for row := 0; row < n-1; row++ {
		srcRow := row
		if row >= i {
			srcRow++
		}
		minorMatrix[row] = make([]T, n-1)
		for col := 0; col < n-1; col++ {
			srcCol := col
			if col >= j {
				srcCol++
			}
			minorMatrix[row][col] = matrix[srcRow][srcCol]
		}
	}
	return Determinant(minorMatrix)
}
