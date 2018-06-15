package rbtree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type rbtreeComparator struct {
}

// Compare Compare
func (tc rbtreeComparator) Compare(k1, k2 interface{}) int {
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

func Test_RbTreeInsert(t *testing.T) {
	tree := NewTree(rbtreeComparator{})
	fmt.Println("插入3")
	tree.Insert(3, 3)
	fmt.Println(tree)

	fmt.Println("\n插入2")
	tree.Insert(2, 2)
	fmt.Println(tree)

	fmt.Println("\n插入1")
	tree.Insert(1, 1) // test case 2.4
	fmt.Println(tree)

	fmt.Println("\n插入4")
	tree.Insert(4, 4) // test case 2.3
	fmt.Println(tree)

	fmt.Println("\n插入5")
	tree.Insert(5, 5) // test case 2.5
	fmt.Println(tree)

	fmt.Println("\n插入10")
	tree.Insert(10, 10)
	fmt.Println(tree)

	fmt.Println("\n插入6")
	tree.Insert(6, 6) // test case 2.6
	fmt.Println(tree)

	fmt.Println("\n插入8")
	tree.Insert(8, 8)
	fmt.Println(tree)

	fmt.Println("\n插入9")
	tree.Insert(9, 9) // test case 2.7
	fmt.Println(tree)

	// fmt.Printf("插入测试\n")
	// fmt.Printf("\n输入数组: \n")
	// fmt.Println("[3, 2, 1, 4, 5, 10, 6, 8, 9]")
	// fmt.Printf("\n中序遍历结果: \n")
	// fmt.Println(tree)
	// fmt.Printf("\n测试结果是否正确: %v\n", tree.Verify())

	// if err := tree.Dot(); err != nil {
	// 	fmt.Printf("Dot error %v\n", err)
	// }

	if !tree.Verify() {
		t.Fatal("Test_RbTreeInsert err")
	}
}

func Test_RbTreeRandInsert(t *testing.T) {
	tree := NewTree(rbtreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// fmt.Printf("\n随机插入%v个数\n", num)
	// fmt.Printf("\n输入数组: \n")
	// fmt.Println(array)

	for _, v := range array {
		tree.Insert(v, v)
	}

	// fmt.Printf("\n中序遍历结果: \n")
	// fmt.Println(tree)
	// fmt.Printf("\n测试结果是否正确: %v\n", tree.Verify())
	if !tree.Verify() {
		t.Fatal("Test_RbTreeRandInsert err")
	}

	idx := 0
	iter := NewIterator(tree)
	for iter.Next() {
		if iter.GetKey().(int) != idx {
			t.Fatalf("want %v, got %v\n", idx, iter.GetKey().(int))
		}

		idx++
	}
}

func Test_RbTreeDelete(t *testing.T) {
	tree := NewTree(rbtreeComparator{})
	tree.Insert(3, 3)
	tree.Insert(2, 2)
	tree.Insert(1, 1)
	tree.Insert(4, 4)
	tree.Insert(5, 5)
	tree.Insert(10, 10)
	tree.Insert(6, 6)
	tree.Insert(8, 8)
	tree.Insert(9, 9)
	fmt.Println("红黑树:")
	fmt.Println(tree)
}

func Test_RbTreeRandDelete(t *testing.T) {
	tree := NewTree(rbtreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}

	// fmt.Println(tree)

	// if !tree.Verify() {
	// 	fmt.Printf("插入错误, array %v\n", array)
	// } else {
	// 	fmt.Println("插入正确")
	// }

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)

	for _, v := range dArray {
		tree.Delete(v)
	}

	// fmt.Println("")
	// fmt.Println(tree)

	// if !tree.Verify() {
	// 	fmt.Printf("删除错误, array %v, dArray %v\n", array, dArray)
	// } else {
	// 	fmt.Printf("删除正确\n")
	// }

	if !tree.Verify() {
		t.Fatal("Test_RbTreeRandDelete err")
	}

	idx := num / 2
	iter := NewIterator(tree)
	for iter.Next() {
		if iter.GetKey().(int) != idx {
			t.Fatalf("want %v, got %v\n", idx, iter.GetKey().(int))
		}

		idx++
	}
}

func Test_RbTreeSearch(t *testing.T) {
	tree := NewTree(rbtreeComparator{})

	tree.Insert(3, 3)
	tree.Insert(2, 2)
	tree.Insert(1, 1)
	tree.Insert(4, 4)
	tree.Insert(5, 5)
	tree.Insert(10, 10)
	tree.Insert(6, 6)
	tree.Insert(8, 8)
	tree.Insert(9, 9)

	node := tree.Search(10)
	if node == nil || node.GetKey().(int) != 10 {
		t.Fatalf("查找错误: 查找节点%v\n", 10)
	}

	node = tree.Search(11)
	if node != nil {
		t.Fatalf("查找错误: 查找节点%v\n", 11)
	}
}

func Test_RbTreeSearchRange(t *testing.T) {
	tree := NewTree(rbtreeComparator{})

	tree.Insert(3, 3)
	tree.Insert(2, 2)
	tree.Insert(1, 1)
	tree.Insert(4, 4)
	tree.Insert(5, 5)
	tree.Insert(10, 10)
	tree.Insert(6, 6)
	tree.Insert(8, 8)
	tree.Insert(9, 9)

	// fmt.Println(tree)
	nodeList := tree.SearchRange(7, 11)
	// for i, v := range nodeList {
	// 	fmt.Printf("第%v个节点 -- key: %v\n", i+1, v.key)
	// }

	verifArr := []int{8, 9, 10}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatalf("Test_RbTreeSearchRange err: v.GetKey().(int) != verifArr[%d], v.GetKey().(int): %v, verifArr[%d]: %v\n", i, v.GetKey().(int), i, verifArr[i])
		}
	}

	if len(nodeList) != len(verifArr) {
		t.Fatalf("Test_RbTreeSearchRange err: len(nodeList) != len(verifArr), len(nodeList): %v, len(verifArr): %v\n", len(nodeList), len(verifArr))
	}
}

func Test_RbTreeSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(rbtreeComparator{})

	tree.Insert(3, 3)
	tree.Insert(2, 2)
	tree.Insert(1, 1)
	tree.Insert(4, 4)
	tree.Insert(5, 5)
	tree.Insert(10, 10)
	tree.Insert(6, 6)
	tree.Insert(8, 8)
	tree.Insert(9, 9)

	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(6, 2)
	// for i, node := range nodeList {
	// 	fmt.Printf("第%v个节点: \t%v\n", i+1, node.key)
	// }

	verifArr := []int{6, 8}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatalf("Test_RbTreeSearchRangeLowerBoundKeyWithLimit err: v.GetKey().(int) != verifArr[%d], v.GetKey().(int): %v, verifArr[%d]: %v\n", i, v.GetKey().(int), i, verifArr[i])
		}
	}

	if len(nodeList) != len(verifArr) {
		t.Fatalf("Test_RbTreeSearchRangeLowerBoundKeyWithLimit err: len(nodeList) != len(verifArr), len(nodeList): %v, len(verifArr): %v\n", len(nodeList), len(verifArr))
	}
}

func Test_RbTreeRandSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(rbtreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}

	// fmt.Println(tree)

	// if !tree.Verify() {
	// 	fmt.Printf("插入错误, array %v\n", array)
	// } else {
	// 	fmt.Println("插入正确")
	// }

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey+i != node.GetKey().(int) {
			t.Fatal("Test_RbTreeRandSearchRangeLowerBoundKeyWithLimit err")
		}
	}
}

func Test_RbTreeSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(rbtreeComparator{})

	tree.Insert(3, 3)
	tree.Insert(2, 2)
	tree.Insert(1, 1)
	tree.Insert(4, 4)
	tree.Insert(5, 5)
	tree.Insert(10, 10)
	tree.Insert(6, 6)
	tree.Insert(8, 8)
	tree.Insert(9, 9)

	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(3, 5)
	// for i, node := range nodeList {
	// 	fmt.Printf("第%v个节点: \t%v\n", i+1, node.key)
	// }

	verifArr := []int{1, 2, 3}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatalf("Test_RbTreeSearchRangeUpperBoundKeyWithLimit err: v.GetKey().(int) != verifArr[%d], v.GetKey().(int): %v, verifArr[%d]: %v\n", i, v.GetKey().(int), i, verifArr[i])
		}
	}

	if len(nodeList) != len(verifArr) {
		t.Fatalf("Test_RbTreeSearchRangeLowerBoundKeyWithLimit err: len(nodeList) != len(verifArr), len(nodeList): %v, len(verifArr): %v\n", len(nodeList), len(verifArr))
	}
}

func Test_RbTreeRandSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(rbtreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}

	// fmt.Println(tree)

	// if !tree.Verify() {
	// 	fmt.Printf("插入错误, array %v\n", array)
	// } else {
	// 	fmt.Println("插入正确")
	// }
	if !tree.Verify() {
		t.Fatal("Test_RbTreeRandSearchRangeUpperBoundKeyWithLimit err")
	}

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey-(len(nodeList)-i-1) != node.GetKey().(int) {
			t.Fatal("Test_RbTreeRandSearchRangeUpperBoundKeyWithLimit err: 查找结果错误")
		}
	}
}
