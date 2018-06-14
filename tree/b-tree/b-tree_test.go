package btree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 0, 1, 5, 2, 3, 6, 4, 7, 9, 8
// 					3
//			1				5,		7
//		0		2		4		6		8,	9

type btreeComparator struct {
}

// Compare Compare
func (tc btreeComparator) Compare(k1, k2 interface{}) int {
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

const DEGREE = 10

// Insert Test
func Test_BTreeRandInsert(t *testing.T) {
	tree := NewTree(DEGREE, btreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{8, 6, 7, 0, 9, 3, 5, 4, 2, 1}

	for _, v := range array {
		tree.Insert(v, v)
	}

	if !tree.Verify() {
		t.Fatalf("插入 - 验证b树错误: 数组 %v\n", array)
	}

	// fmt.Println(tree)
	// if err := tree.Dot(); err != nil {
	// 	fmt.Printf("Dot error %v\n", err)
	// }
	// fmt.Println(tree.size)

	idx := 0
	iter := NewIterator(tree)
	for iter.Next() {
		if iter.GetKey().(int) != idx {
			t.Fatalf("Test_BTreeRandInsert err\n")
		}

		idx++
	}
}

func Test_BatchBTreeRandInsert(t *testing.T) {
	count := 1000

	for i := 0; i < count; i++ {
		tree := NewTree(DEGREE, btreeComparator{})
		var num = 100

		array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

		for _, v := range array {
			tree.Insert(v, v)
		}

		if !tree.Verify() {
			t.Fatal("Test_BatchBTreeRandInsert err")
		}
	}
}

// Delete Test
func Test_BTreeDelete(t *testing.T) {
	tree := NewTree(DEGREE, btreeComparator{})
	array := []int{0, 1, 5, 2, 3, 6, 4, 7, 9, 8}

	for _, v := range array {
		tree.Insert(v, v)
	}

	fmt.Println(tree)

	fmt.Println("\n删除3")
	tree.Delete(3) // test MergeNode - bBig = true
	fmt.Println(tree)

	fmt.Println("\n删除6")
	tree.Delete(6)
	fmt.Println(tree)

	fmt.Println("\n删除5")
	tree.Delete(5) // test MoveKey - bBig = true
	fmt.Println(tree)

	fmt.Println("\n删除9")
	tree.Delete(9) // test MergeNode - bBig = false
	fmt.Println(tree)

	fmt.Println("\n删除2")
	tree.Delete(2)
	fmt.Println(tree)

	fmt.Println("\n删除1")
	tree.Delete(1)
	fmt.Println(tree)

	fmt.Println("\n删除8")
	tree.Delete(8) // test MoveKey - bBig = false
	fmt.Println(tree)

	fmt.Println("\n删除0")
	tree.Delete(0)
	fmt.Println(tree)

	fmt.Println("\n删除7")
	tree.Delete(7)
	fmt.Println(tree)

	fmt.Println("\n删除4")
	tree.Delete(4)
	fmt.Println(tree)
}

func Test_BTreeRandDelete(t *testing.T) {
	tree := NewTree(DEGREE, btreeComparator{})
	var num = 1000000

	insertArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	for _, v := range insertArray {
		tree.Insert(v, v)
	}

	if !tree.Verify() {
		t.Fatalf("插入 - 验证b树错误: 数组 %v\n", insertArray)
	}

	deleteArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	for _, v := range deleteArray {
		tree.Delete(v)
	}

	if !tree.Verify() {
		t.Fatalf("删除 - 验证b树错误: 数组 %v\n", deleteArray)
	}

	// fmt.Println(tree)

	idx := num / 2
	iter := NewIterator(tree)
	for iter.Next() {
		if iter.GetKey().(int) != idx {
			t.Fatalf("Test_BTreeRandDelete err")
		}

		idx++
	}
}

func Test_BatchBTreeRandDelete(t *testing.T) {
	count := 1000

	for i := 0; i < count; i++ {
		tree := NewTree(DEGREE, btreeComparator{})
		var num = 100

		insertArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
		for _, v := range insertArray {
			tree.Insert(v, v)
		}

		if !tree.Verify() {
			t.Fatalf("插入 - 验证b树错误: 数组 %v\n", insertArray)
		}

		deleteArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
		for _, v := range deleteArray {
			tree.Delete(v)
		}

		if !tree.Verify() {
			t.Fatalf("删除 - 验证b树错误: 数组 %v\n", deleteArray)
		}
	}
}

func Test_BTreeSearch(t *testing.T) {
	tree := NewTree(DEGREE, btreeComparator{})
	var num = 1000000

	insertArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	for _, v := range insertArray {
		tree.Insert(v, v)
	}

	if !tree.Verify() {
		t.Fatalf("插入 - 验证b树错误: 数组 %v\n", insertArray)
	}

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	entry := tree.Search(sKey)
	if entry.GetKey().(int) != sKey {
		t.Fatalf("查找错误\n")
	}
}

func Test_BTreeeRangeSearch(t *testing.T) {
	var num = 1000000
	var diff = 1000
	tree := NewTree(DEGREE, btreeComparator{})

	insertArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	for _, v := range insertArray {
		tree.Insert(v, v)
	}

	if !tree.Verify() {
		t.Fatalf("插入 - 验证b树错误: 数组 %v\n", insertArray)
	}

	// fmt.Println(tree)

	minKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	maxKey := minKey + diff

	entryList := tree.SearchRange(minKey, maxKey)
	if len(entryList) != diff+1 {
		t.Fatalf("len(entryList) != diff\n")
	}

	for i, v := range entryList {
		if minKey+i != v.GetKey().(int) {
			t.Fatalf("SearchRange错误, minKey+Key(i) != v.Key\n")
		}

		if v.GetKey().(int) > maxKey {
			t.Fatalf("SearchRange错误, v.Key > maxKey\n")
		}
	}
}
