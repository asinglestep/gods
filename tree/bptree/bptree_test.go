package bptree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_BpTreeInsert(t *testing.T) {
	tree := NewTree(2)

	tree.Insert(NewEntry(10, 10))
	tree.Insert(NewEntry(20, 20))
	tree.Insert(NewEntry(30, 30))
	tree.Insert(NewEntry(40, 40))

	tree.Insert(NewEntry(50, 50)) // split

	tree.Insert(NewEntry(25, 25))

	tree.Insert(NewEntry(5, 5))

	tree.Insert(NewEntry(1, 1))

	tree.Insert(NewEntry(7, 7))

	tree.Insert(NewEntry(60, 60))

	tree.Insert(NewEntry(27, 27))
	tree.Insert(NewEntry(28, 28))
	tree.Insert(NewEntry(70, 70))
	tree.Insert(NewEntry(80, 80))

	tree.Insert(NewEntry(90, 90))

	if !tree.VerifBpTree() {
		t.Fatal("BpTree Insert")
	}

	// tree.PrintBpTree()
	if err := tree.Dot(); err != nil {
		fmt.Printf("Dot error %v\n", err)
	}
}

func Test_BpTreeRandInsert(t *testing.T) {
	tree := NewTree(2)
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 15, 7, 18, 0, 4, 2, 3, 16, 5, 13, 10, 19, 12, 11, 14, 17, 6, 8, 9}

	for _, v := range array {
		tree.Insert(NewEntry(Key(v), Value(v)))
	}

	// tree.PrintBpTree()

	if !tree.VerifBpTree() {
		t.Fatalf("Bptree Insert Error\n")
	}
}

func Test_BpTreeDelete(t *testing.T) {
	tree := NewTree(2)

	tree.Insert(NewEntry(10, 10))
	tree.Insert(NewEntry(20, 20))
	tree.Insert(NewEntry(30, 30))
	tree.Insert(NewEntry(40, 40))

	tree.Insert(NewEntry(50, 50)) // split

	tree.Insert(NewEntry(25, 25))

	tree.Insert(NewEntry(5, 5))

	tree.Insert(NewEntry(1, 1))

	tree.Insert(NewEntry(7, 7))

	tree.Insert(NewEntry(60, 60))

	tree.Insert(NewEntry(27, 27))
	tree.Insert(NewEntry(28, 28))
	tree.Insert(NewEntry(70, 70))
	tree.Insert(NewEntry(80, 80))

	tree.Insert(NewEntry(90, 90))

	if !tree.VerifBpTree() {
		t.Fatal("BpTree Insert")
	}

	tree.PrintBpTree()

	tree.Delete(28)
	fmt.Println("\n删除28")
	tree.PrintBpTree()

	tree.Delete(60) // test 修复2
	fmt.Println("\n删除60")
	tree.PrintBpTree()

	tree.Delete(10)
	fmt.Println("\n删除10")
	tree.PrintBpTree()

	tree.Delete(7)
	fmt.Println("\n删除7")
	tree.PrintBpTree()

	tree.Delete(1)
	fmt.Println("\n删除1")
	tree.PrintBpTree()

	tree.Delete(25)
	fmt.Println("\n删除25")
	tree.PrintBpTree()

	tree.Delete(5) // test 修复1
	fmt.Println("\n删除5")
	tree.PrintBpTree()

	tree.Delete(50) // test 修复3
	fmt.Println("\n删除50")
	tree.PrintBpTree()

	tree.Delete(30)
	fmt.Println("\n删除30")
	tree.PrintBpTree()

	tree.Delete(40)
	fmt.Println("\n删除40")
	tree.PrintBpTree()

	tree.Delete(80)
	fmt.Println("\n删除80")
	tree.PrintBpTree()

	tree.Delete(20)
	fmt.Println("\n删除20")
	tree.PrintBpTree()

	tree.Delete(27)
	fmt.Println("\n删除27")
	tree.PrintBpTree()

	tree.Delete(90)
	fmt.Println("\n删除90")
	tree.PrintBpTree()

	tree.Delete(70)
	fmt.Println("\n删除70")
	tree.PrintBpTree()
}

func Test_BpTreeRandDelete(t *testing.T) {
	tree := NewTree(2)
	var num = 10000000

	iArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// iArray = []int{27, 28, 12, 4, 21, 26, 9, 3, 6, 25, 8, 15, 17, 24, 7, 2, 29, 23, 1, 18, 10, 0, 5, 13, 20, 16, 19, 22, 14, 11}

	// fmt.Printf("插入数据 %v\n", iArray)

	for _, v := range iArray {
		tree.Insert(NewEntry(Key(v), Value(v)))
	}

	if !tree.VerifBpTree() {
		t.Fatalf("Bptree Insert Error, iArray %v\n", iArray)
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	// dArray = []int{12, 14, 1, 9, 7, 4, 2, 13, 0, 11, 6, 10, 3, 5, 8}

	// fmt.Printf("删除数据 %v\n", dArray)

	for _, v := range dArray {
		tree.Delete(Key(v))
	}

	// tree.Dot()

	if !tree.VerifBpTree() {
		t.Fatalf("Bptree Delete Error, dArray %v\n", dArray)
	}
}

func Test_BpTreeRandSearch(t *testing.T) {
	tree := NewTree(2)
	var num = 10000000

	iArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// iArray = []int{27, 28, 12, 4, 21, 26, 9, 3, 6, 25, 8, 15, 17, 24, 7, 2, 29, 23, 1, 18, 10, 0, 5, 13, 20, 16, 19, 22, 14, 11}

	// fmt.Printf("插入数据 %v\n", iArray)

	for _, v := range iArray {
		tree.Insert(NewEntry(Key(v), Value(v)))
	}

	if !tree.VerifBpTree() {
		t.Fatalf("Bptree Insert Error, iArray %v\n", iArray)
	}

	sKey := Key(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num))
	// fmt.Printf("查找数据 %v\n", sKey)

	entry := tree.Search(sKey)
	if entry.Key != sKey {
		t.Fatalf("Bptree Search Error, sKey %v\n", sKey)
	}
}

func Test_BpTreeRandSearchRange(t *testing.T) {
	tree := NewTree(2)
	var num = 10000000

	iArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// iArray = []int{27, 28, 12, 4, 21, 26, 9, 3, 6, 25, 8, 15, 17, 24, 7, 2, 29, 23, 1, 18, 10, 0, 5, 13, 20, 16, 19, 22, 14, 11}

	// fmt.Printf("插入数据 %v\n", iArray)

	for _, v := range iArray {
		tree.Insert(NewEntry(Key(v), Value(v)))
	}

	if !tree.VerifBpTree() {
		t.Fatalf("Bptree Insert Error, iArray %v\n", iArray)
	}

	minKey := Key(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num))
	maxKey := minKey + 1000

	entries := tree.SearchRange(minKey, maxKey)
	for i, v := range entries {
		if minKey+Key(i) != v.Key {
			t.Fatalf("Bptree SearchRange Error, minKey %v\n", minKey)
		}

		if maxKey < v.Key {
			t.Fatalf("Bptree SearchRange Error, maxKey %v\n", maxKey)
		}
	}
}
