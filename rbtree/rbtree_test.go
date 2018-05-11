package rbtree

import (
	"fmt"
	"testing"
)

const _N = 2000

var arr = []int{34, 5, 12, 15, 71, 45, 2, 33, 3, 83, 61, 43, 11, 21, 22, 9}

func intComp(a0, b0 interface{}) int {
	a := a0.(int)
	b := b0.(int)
	if a < b {
		return -1
	} else if a == b {
		return 0
	} else {
		return 1
	}
}

func TestRbtree(t *testing.T) {
	tree := new(RBTree)
	for _, i := range arr {
		n := &Node{Item: i}
		InsertWithoutBalance(tree, n, intComp, COEXIST_IF_EXIST)
		InsertCase(tree, n)
	}
	it := NewIterator(tree, true)
	for {
		if n, end := it.Next(); end {
			break
		} else {
			fmt.Printf("%v ", n.Item)
		}
	}
	fmt.Printf("\n")
	fmt.Println(tree)

	for _, i := range arr {
		n := Find(tree, i, intComp)
		DeleteNode(tree, n)
		fmt.Println(tree)
	}
}

func TestRBMap(t *testing.T) {
	m := NewMap(intComp)
	for i, j := range arr {
		m.Set(i, j)
	}
	for i := 0; i < len(arr); i++ {
		fmt.Println(m.Get(i))
	}
}

func BenchmarkWriteMap(b *testing.B) {
	a := make([]int, _N)
	for i := 0; i < len(a); i++ {
		a[i] = i
	}

	for i := 0; i < b.N; i++ {
		m := make(map[int]int)
		for k, v := range a {
			m[k] = v
		}
		for k, _ := range a {
			delete(m, k)
		}
	}
}

func BenchmarkReadMap(b *testing.B) {
	a := make([]int, _N)
	for i := 0; i < len(a); i++ {
		a[i] = i
	}
	m := make(map[int]int)
	for k, v := range a {
		m[k] = v
	}

	for i := 0; i < b.N; i++ {
		for k, _ := range a {
			if _, exist := m[k]; exist {
			}
		}
	}
}

func BenchmarkWriteRBMap(b *testing.B) {
	a := make([]int, _N)
	for i := 0; i < len(a); i++ {
		a[i] = i
	}

	for i := 0; i < b.N; i++ {
		m := NewMap(intComp)
		for k, v := range a {
			m.Set(k, v)
		}
		for k, _ := range a {
			m.Delete(k)
		}
	}
}

func BenchmarkReadRBMap(b *testing.B) {
	a := make([]int, _N)
	for i := 0; i < len(a); i++ {
		a[i] = i
	}
	m := NewMap(intComp)
	for k, v := range a {
		m.Set(k, v)
	}

	for i := 0; i < b.N; i++ {
		for k, _ := range a {
			m.Get(k)
		}
	}
}
