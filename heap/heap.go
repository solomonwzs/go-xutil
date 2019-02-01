package heap

import (
	"reflect"
)

func heapify(less func(i, j int) bool, swap func(i, j int), start, end int) {
	dad := start
	son := dad*2 + 1
	for son <= end {
		if son+1 <= end && less(son+1, son) {
			son += 1
		}
		if less(dad, son) {
			return
		} else {
			swap(dad, son)
			dad = son
			son = dad*2 + 1
		}
	}
}

func PercolateDown(slice interface{}, less func(i, j int) bool) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()

	swap(length-1, 0)
	heapify(less, swap, 0, length-2)
}

func PercolateUp(slice interface{}, less func(i, j int) bool) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()

	for son := length - 1; son > 0; {
		dad := (son - 1) / 2
		if less(dad, son) {
			return
		} else {
			swap(dad, son)
			son = dad
		}
	}
}

func SliceHeapify(slice interface{}, less func(i, j int) bool) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()

	for i := length/2 - 1; i >= 0; i-- {
		heapify(less, swap, i, length-1)
	}
}
