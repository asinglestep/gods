package bptree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type bptreeComparator struct {
}

// Compare Compare
func (tc bptreeComparator) Compare(k1, k2 interface{}) int {
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

func Test_BpTreeInsert(t *testing.T) {
	tree := NewTree(2, bptreeComparator{})

	fmt.Println("插入10")
	tree.Insert(10, 10)
	fmt.Println(tree)

	fmt.Println("插入20")
	tree.Insert(20, 20)
	fmt.Println(tree)

	fmt.Println("插入30")
	tree.Insert(30, 30)
	fmt.Println(tree)

	fmt.Println("插入40")
	tree.Insert(40, 40)
	fmt.Println(tree)

	fmt.Println("插入50")
	tree.Insert(50, 50) // split
	fmt.Println(tree)

	fmt.Println("插入25")
	tree.Insert(25, 25)
	fmt.Println(tree)

	fmt.Println("插入5")
	tree.Insert(5, 5)
	fmt.Println(tree)

	fmt.Println("插入1")
	tree.Insert(1, 1)
	fmt.Println(tree)

	fmt.Println("插入7")
	tree.Insert(7, 7)
	fmt.Println(tree)

	fmt.Println("插入60")
	tree.Insert(60, 60)
	fmt.Println(tree)

	fmt.Println("插入27")
	tree.Insert(27, 27)
	fmt.Println(tree)

	fmt.Println("插入28")
	tree.Insert(28, 28)
	fmt.Println(tree)

	fmt.Println("插入70")
	tree.Insert(70, 70)
	fmt.Println(tree)

	fmt.Println("插入80")
	tree.Insert(80, 80)
	fmt.Println(tree)

	fmt.Println("插入90")
	tree.Insert(90, 90)
	fmt.Println(tree)

	if !tree.Verify() {
		t.Fatal("BpTree Insert")
	}

	// fmt.Println(tree)
	// if err := tree.Dot(); err != nil {
	// 	fmt.Printf("Dot error %v\n", err)
	// }

	idx := 0
	verifyArr := []int{1, 5, 7, 10, 20, 25, 27, 28, 30, 40, 50, 60, 70, 80, 90}

	iter := NewIterator(tree)
	for iter.Next() {
		if iter.entry.GetKey().(int) != verifyArr[idx] {
			t.Fatalf("Test_BpTreeInsert Iterator Error\n")
		}

		idx++
	}
}

func Test_BpTreeRandInsert(t *testing.T) {
	tree := NewTree(2, bptreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 15, 7, 18, 0, 4, 2, 3, 16, 5, 13, 10, 19, 12, 11, 14, 17, 6, 8, 9}

	for _, v := range array {
		tree.Insert(v, v)
	}

	// fmt.Println(tree)

	if !tree.Verify() {
		t.Fatalf("Bptree Insert Error\n")
	}

	idx := 0
	iter := NewIterator(tree)
	for iter.Next() {
		if iter.entry.GetKey().(int) != idx {
			t.Fatalf("Test_BpTreeInsert Iterator Error\n")
		}

		idx++
	}
}

func Test_BpTreeDelete(t *testing.T) {
	tree := NewTree(2, bptreeComparator{})

	tree.Insert(10, 10)
	tree.Insert(20, 20)
	tree.Insert(30, 30)
	tree.Insert(40, 40)

	tree.Insert(50, 50) // split

	tree.Insert(25, 25)

	tree.Insert(5, 5)

	tree.Insert(1, 1)

	tree.Insert(7, 7)

	tree.Insert(60, 60)

	tree.Insert(27, 27)
	tree.Insert(28, 28)
	tree.Insert(70, 70)
	tree.Insert(80, 80)

	tree.Insert(90, 90)

	if !tree.Verify() {
		t.Fatal("BpTree Insert")
	}

	fmt.Println(tree)

	tree.Delete(28)
	fmt.Println("\n删除28")
	fmt.Println(tree)

	tree.Delete(60) // test 修复2
	fmt.Println("\n删除60")
	fmt.Println(tree)

	tree.Delete(10)
	fmt.Println("\n删除10")
	fmt.Println(tree)

	tree.Delete(7)
	fmt.Println("\n删除7")
	fmt.Println(tree)

	tree.Delete(1)
	fmt.Println("\n删除1")
	fmt.Println(tree)

	tree.Delete(25)
	fmt.Println("\n删除25")
	fmt.Println(tree)

	tree.Delete(5) // test 修复1
	fmt.Println("\n删除5")
	fmt.Println(tree)

	tree.Delete(50) // test 修复3
	fmt.Println("\n删除50")
	fmt.Println(tree)

	tree.Delete(30)
	fmt.Println("\n删除30")
	fmt.Println(tree)

	tree.Delete(40)
	fmt.Println("\n删除40")
	fmt.Println(tree)

	tree.Delete(80)
	fmt.Println("\n删除80")
	fmt.Println(tree)

	tree.Delete(20)
	fmt.Println("\n删除20")
	fmt.Println(tree)

	tree.Delete(27)
	fmt.Println("\n删除27")
	fmt.Println(tree)

	tree.Delete(90)
	fmt.Println("\n删除90")
	fmt.Println(tree)

	tree.Delete(70)
	fmt.Println("\n删除70")
	fmt.Println(tree)
}

func Test_BpTreeRandDelete(t *testing.T) {
	tree := NewTree(2, bptreeComparator{})
	var num = 1000000

	iArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// iArray = []int{27, 28, 12, 4, 21, 26, 9, 3, 6, 25, 8, 15, 17, 24, 7, 2, 29, 23, 1, 18, 10, 0, 5, 13, 20, 16, 19, 22, 14, 11}

	// fmt.Printf("插入数据 %v\n", iArray)

	for _, v := range iArray {
		tree.Insert(v, v)
	}

	if !tree.Verify() {
		t.Fatalf("Bptree Insert Error, iArray %v\n", iArray)
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	// dArray = []int{12, 14, 1, 9, 7, 4, 2, 13, 0, 11, 6, 10, 3, 5, 8}

	// fmt.Printf("删除数据 %v\n", dArray)

	for _, v := range dArray {
		tree.Delete(v)
	}

	// tree.Dot()

	if !tree.Verify() {
		t.Fatalf("Bptree Delete Error, dArray %v\n", dArray)
	}
}

func Test_BpTreeRandSearch(t *testing.T) {
	tree := NewTree(2, bptreeComparator{})
	var num = 1000000

	iArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// iArray = []int{27, 28, 12, 4, 21, 26, 9, 3, 6, 25, 8, 15, 17, 24, 7, 2, 29, 23, 1, 18, 10, 0, 5, 13, 20, 16, 19, 22, 14, 11}

	// fmt.Printf("插入数据 %v\n", iArray)

	for _, v := range iArray {
		tree.Insert(v, v)
	}

	if !tree.Verify() {
		t.Fatalf("Bptree Insert Error, iArray %v\n", iArray)
	}

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	// fmt.Printf("查找数据 %v\n", sKey)

	entry := tree.Search(sKey)
	if entry.GetKey().(int) != sKey {
		t.Fatalf("Bptree Search Error, sKey %v\n", sKey)
	}
}

func Test_BpTreeRandSearchRange(t *testing.T) {
	tree := NewTree(2, bptreeComparator{})
	var num = 1000000

	iArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// iArray = []int{27, 28, 12, 4, 21, 26, 9, 3, 6, 25, 8, 15, 17, 24, 7, 2, 29, 23, 1, 18, 10, 0, 5, 13, 20, 16, 19, 22, 14, 11}

	// fmt.Printf("插入数据 %v\n", iArray)

	for _, v := range iArray {
		tree.Insert(v, v)
	}

	if !tree.Verify() {
		t.Fatalf("Bptree Insert Error, iArray %v\n", iArray)
	}

	minKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	maxKey := minKey + 1000

	entries := tree.SearchRange(minKey, maxKey)
	for i, v := range entries {
		if minKey+i != v.GetKey().(int) {
			t.Fatalf("Bptree SearchRange Error, minKey %v\n", minKey)
		}

		if maxKey < v.GetKey().(int) {
			t.Fatalf("Bptree SearchRange Error, maxKey %v\n", maxKey)
		}
	}
}
