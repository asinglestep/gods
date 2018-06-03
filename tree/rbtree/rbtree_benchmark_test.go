package rbtree

import (
	"math/rand"
	"testing"
	"time"
)

func Benchmark_RbTreeRandInsert(b *testing.B) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for i := 0; i < b.N; i++ {
		for _, v := range array {
			tree.Insert(Key(v))
		}
	}
}

func Benchmark_RbTreeRandDelete(b *testing.B) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(Key(v))
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range dArray {
		tree.Delete(Key(v))
	}
}

func Benchmark_RbTreeRandSearchRangeLowerBoundKeyWithLimit(b *testing.B) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(Key(v))
	}

	sKey := Key(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num))
	tree.SearchRangeLowerBoundKeyWithLimit(sKey, 1000)
}
