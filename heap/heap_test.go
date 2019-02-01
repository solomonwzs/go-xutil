package heap

import (
	"testing"
)

func TestHeap0(t *testing.T) {
	a := []int{9, 2, 1, 5, 6, 3, 8, 7, 0, 4}
	less := func(i, j int) bool { return a[i] < a[j] }
	SliceHeapify(a, less)
	t.Log(a)

	for length := len(a) - 1; length >= 0; length-- {
		PercolateDown(a, less)
		t.Log(a[length])
		a = a[:length]
	}
}

func TestHeap1(t *testing.T) {
	a := []int{9, 2, 1, 5, 6, 3, 8, 7, 0, 4}
	heap := make([]int, 0, len(a))
	less := func(i, j int) bool {
		return heap[i] < heap[j]
	}

	for _, i := range a {
		heap = append(heap, i)
		PercolateUp(heap, less)
	}
	t.Log(heap)

	for length := len(heap) - 1; length >= 0; length-- {
		PercolateDown(heap, less)
		t.Log(heap[length])
		heap = heap[:length]
	}
}
