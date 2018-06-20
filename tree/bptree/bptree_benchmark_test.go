package bptree

import (
	"math/rand"
	"testing"
	"time"
)

func Benchmark_RandInsert(b *testing.B) {
	tree := NewTree(2, bptreeComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{0, 1, 5, 2, 3, 6, 4, 7, 9, 8}

	for _, v := range array {
		tree.Insert(v, v)
	}
}

func Benchmark_RandDelete(b *testing.B) {
	tree := NewTree(2, bptreeComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{0, 1, 5, 2, 3, 6, 4, 7, 9, 8}

	for _, v := range array {
		tree.Insert(v, v)
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range dArray {
		tree.Delete(v)
	}
}

func Benchmark_RandSearchRange(b *testing.B) {
	// tree := NewTree(2, bptreeComparator{})
	// var num = 10000000

	// iArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	// for _, v := range iArray {
	// 	tree.Insert(v, v)
	// }

	// minKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	// maxKey := minKey + 1000

	// tree.SearchRange(minKey, maxKey)
}
