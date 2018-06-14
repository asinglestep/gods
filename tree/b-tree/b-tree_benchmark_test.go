package btree

import (
	"math/rand"
	"testing"
	"time"
)

func Benchmark_BTreeRandInsert(b *testing.B) {
	tree := NewTree(DEGREE, btreeComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{0, 1, 5, 2, 3, 6, 4, 7, 9, 8}
	// array = []int{8, 6, 7, 0, 9, 3, 5, 4, 2, 1}

	for _, v := range array {
		tree.Insert(v, v)
	}
}

func Benchmark_BTreeRandDelete(b *testing.B) {
	tree := NewTree(DEGREE, btreeComparator{})
	var num = 10000000

	iArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	for _, v := range iArray {
		tree.Insert(v, v)
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	for _, v := range dArray {
		tree.Delete(v)
	}
}

func Benchmark_BTreeRangeSearch(b *testing.B) {
	tree := NewTree(DEGREE, btreeComparator{})
	var num = 10000000

	insertArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	for _, v := range insertArray {
		tree.Insert(v, v)
	}

	minKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	maxKey := minKey + 1000
	tree.SearchRange(minKey, maxKey)
}
