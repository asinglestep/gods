package treap

import (
	"math/rand"
	"testing"
	"time"

	"github.com/asinglestep/gods/utils"
)

func Benchmark_TreapRandInsert(b *testing.B) {
	tree := NewTree(treapComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for i := 0; i < b.N; i++ {
		for _, v := range array {
			tree.Insert(utils.NewEntry(v, v))
		}
	}
}

func Benchmark_TreapRandDelete(b *testing.B) {
	tree := NewTree(treapComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(utils.NewEntry(v, v))
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range dArray {
		tree.Delete(utils.NewEntry(v, v))
	}
}
