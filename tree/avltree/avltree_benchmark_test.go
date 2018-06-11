package avltree

import (
	"math/rand"
	"testing"
	"time"
)

func Benchmark_AvlTreeRandInsert(b *testing.B) {
	tree := NewTree(avltreeComparator{})
	var num = 10000000
	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}
}

func Benchmark_AvlTreeRandDelete(b *testing.B) {

	tree := NewTree(avltreeComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(v, v)
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	for _, v := range dArray {
		tree.Delete(v)
	}
}

func Benchmark_AvlTreeRandRangeSearchBack(b *testing.B) {
	tree := NewTree(avltreeComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(v, v)
	}

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	tree.SearchRangeLowerBoundKeyWithLimit(sKey, 100000)
}
