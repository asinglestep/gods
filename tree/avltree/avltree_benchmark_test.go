package avltree

import (
	"math/rand"
	"testing"
	"time"
)

func Benchmark_AvlTreeRandInsert(b *testing.B) {
	tree := NewTree()
	var num = 10000000
	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for i := 0; i < b.N; i++ {
		for _, v := range array {
			tree.Insert(Key(v))
		}
	}
}

func Benchmark_AvlTreeRandDelete(b *testing.B) {

	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(Key(v))
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	for _, v := range dArray {
		tree.Delete(Key(v))
	}
}

func Benchmark_AvlTreeRandRangeSearchBack(b *testing.B) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(Key(v))
	}

	sKey := Key(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num))
	tree.SearchRangeLowerBoundKeyWithLimit(sKey, 100000)
}
