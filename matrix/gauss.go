package matrix

import (
	"errors"
	"math"
)

func Gauss(a [][]float64) (ans []float64, err error) {
	n := len(a)
	if n == 0 {
		err = errors.New("matrix size error")
		return
	}
	if len(a[0]) != n+1 {
		err = errors.New("matrix size error")
		return
	}

	selectColE(a)
	ans = make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		ans[i] = a[i][n]
		for j := i + 1; j < n; j++ {
			ans[i] -= a[i][j] * ans[j]
		}
		ans[i] /= a[i][i]
	}

	return
}

func selectColE(a [][]float64) (err error) {
	n := len(a)
	for j := 0; j < n; j++ {
		maxRowE := j
		for i := j + 1; i < n; i++ {
			if math.Abs(a[i][j]) > math.Abs(a[maxRowE][j]) {
				maxRowE = i
			}
		}
		if maxRowE != j {
			swapRow(a, j, maxRowE)
		}

		for i := j + 1; i < n; i++ {
			if a[j][j] == 0 {
				return errors.New("not augmented matrix")
			}
			tmp := a[i][j] / a[j][j]
			for k := j; k <= n; k++ {
				a[i][k] -= a[j][k] * tmp
			}
		}
	}
	return
}

func swapRow(a [][]float64, m, n int) {
	tmp := make([]float64, len(a[0]))
	copy(tmp, a[m])
	copy(a[m], a[n])
	copy(a[n], tmp)
}
