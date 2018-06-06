package treap

import (
	"math/rand"
	"testing"
	"time"
)

func Benchmark_TreapRandInsert(b *testing.B) {
	tree := NewTree(treapComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}
}

func Benchmark_TreapRandDelete(b *testing.B) {
	tree := NewTree(treapComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range dArray {
		tree.Delete(v)
	}
}

func Benchmark_TreapRandSearch(b *testing.B) {
	tree := NewTree(treapComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}

	sArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	for _, v := range sArray {
		tree.Search(v)
	}
}
