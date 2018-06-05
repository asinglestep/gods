package treap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/asinglestep/gods/utils"
)

type treapComparator struct {
}

// Compare Compare
func (tc treapComparator) Compare(k1, k2 interface{}) int {
	e1 := k1.(*utils.Entry)
	e2 := k2.(*utils.Entry)

	if e1.GetKey().(int) > e2.GetKey().(int) {
		return 1
	}

	if e1.GetKey().(int) < e2.GetKey().(int) {
		return -1
	}

	return 0
}

func Test_TreapRandInsert(t *testing.T) {
	tree := NewTree(treapComparator{})
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{8, 4, 0, 6, 3, 9, 1, 5, 2, 7}

	for _, v := range array {
		tree.Insert(utils.NewEntry(v, v))
	}

	// fmt.Printf("\n随机插入%v个数\n", num)
	// fmt.Printf("随机数组: %v\n", array)

	fmt.Printf("\n插入结果是否正确: %v\n", tree.VerifTreap())

	// tree.Dot()
	idx := 0
	iter := NewIterator(tree)
	for iter.Next() {
		if iter.GetKey().(int) != idx {
			t.Fatalf("want %v, got %v", idx, iter.GetKey().(int))
		}

		idx++
	}
}

func Test_TreapRandDelete(t *testing.T) {
	tree := NewTree(treapComparator{})
	var num = 10

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{8, 4, 0, 6, 3, 9, 1, 5, 2, 7}
	// fmt.Printf("随机插入数组: %v\n", array)

	for _, v := range array {
		tree.Insert(utils.NewEntry(v, v))
	}

	delArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	// fmt.Printf("随机删除数组: %v\n", delArray)

	for _, v := range delArray {
		tree.Delete(utils.NewEntry(v, v))
	}

	fmt.Printf("\n删除结果是否正确: %v\n", tree.VerifTreap())
	// tree.Dot()
	idx := num / 2
	iter := NewIterator(tree)
	for iter.Next() {
		if iter.GetKey().(int) != idx {
			t.Fatalf("want %v, got %v", idx, iter.GetKey().(int))
		}

		idx++
	}
}
