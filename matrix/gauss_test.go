package matrix

import "testing"

func TestGauss(t *testing.T) {
	a := [][]float64{
		[]float64{0.001, 2, 3, 1},
		[]float64{-1, 3.710, 4.623, 2},
		[]float64{-2, 1.07, 5.643, 3},
	}
	t.Log(Gauss(a))
}
