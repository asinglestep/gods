package treap

import (
	"math/rand"
	"testing"
	"time"
)

type treapComparator struct {
}

// Compare Compare
func (tc treapComparator) Compare(k1, k2 interface{}) int {
	i1 := k1.(int)
	i2 := k2.(int)

	if i1 > i2 {
		return 1
	}

	if i1 < i2 {
		return -1
	}

	return 0
}

func Test_TreapRandInsert(t *testing.T) {
	tree := NewTree(treapComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{8, 4, 0, 6, 3, 9, 1, 5, 2, 7}

	for _, v := range array {
		tree.Insert(v, v)
	}

	// fmt.Printf("\n随机插入%v个数\n", num)
	// fmt.Printf("随机数组: %v\n", array)

	// fmt.Printf("\n插入结果是否正确: %v\n", tree.VerifTreap())

	if !tree.VerifTreap() {
		t.Fatal("Test_TreapRandInsert err")
	}

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
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{8, 4, 0, 6, 3, 9, 1, 5, 2, 7}
	// fmt.Printf("随机插入数组: %v\n", array)

	for _, v := range array {
		tree.Insert(v, v)
	}

	delArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	// fmt.Printf("随机删除数组: %v\n", delArray)

	for _, v := range delArray {
		tree.Delete(v)
	}

	// fmt.Printf("\n删除结果是否正确: %v\n", tree.VerifTreap())
	// tree.Dot()

	if !tree.VerifTreap() {
		t.Fatal("Test_TreapRandDelete err")
	}

	idx := num / 2
	iter := NewIterator(tree)
	for iter.Next() {
		if iter.GetKey().(int) != idx {
			t.Fatalf("want %v, got %v", idx, iter.GetKey().(int))
		}

		idx++
	}
}

func Test_TreapRandSearch(t *testing.T) {
	tree := NewTree(treapComparator{})
	var num = 1000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{8, 4, 0, 6, 3, 9, 1, 5, 2, 7}
	// fmt.Printf("随机插入数组: %v\n", array)

	for _, v := range array {
		tree.Insert(v, v)
	}

	key := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	// fmt.Printf("查找key: %v\n", key)

	node := tree.Search(key)
	if node.GetKey().(int) != key {
		t.Fatalf("want %v, got %v\n", key, node.GetKey().(int))
	}
}
