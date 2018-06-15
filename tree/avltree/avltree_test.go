package avltree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type avltreeComparator struct {
}

// Compare Compare
func (ac avltreeComparator) Compare(k1, k2 interface{}) int {
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

func Test_AvlTreeInsert(t *testing.T) {
	tree := NewTree(avltreeComparator{})

	fmt.Println("插入20")
	tree.Insert(20, 20)
	fmt.Println(tree)

	fmt.Println("\n插入40")
	tree.Insert(40, 40)
	fmt.Println(tree)

	fmt.Println("\n插入60")
	tree.Insert(60, 60) // test case 2.3
	fmt.Println(tree)

	fmt.Println("\n插入80")
	tree.Insert(80, 80)
	fmt.Println(tree)

	fmt.Println("\n插入70")
	tree.Insert(70, 70) // test case 2.4
	fmt.Println(tree)

	fmt.Println("\n插入15")
	tree.Insert(15, 15)
	fmt.Println(tree)

	fmt.Println("\n插入18")
	tree.Insert(18, 18) // test case 2.2
	fmt.Println(tree)

	fmt.Println("\n插入10")
	tree.Insert(10, 10)
	fmt.Println(tree)

	fmt.Println("\n插入5")
	tree.Insert(5, 5) // test case 2.1
	fmt.Println(tree)

	// fmt.Printf("插入测试\n")
	// fmt.Printf("\n输入数组: \n")
	// fmt.Println("[20, 40, 60, 80, 70, 15, 18, 10, 5]")
	// fmt.Printf("\n中序遍历结果: \n")
	// fmt.Println(tree)
	// if err := tree.Dot(); err != nil {
	// 	fmt.Printf("Dot Error %v\n", err)
	// }
	// fmt.Printf("\n测试结果是否正确: %v\n", tree.Verify())

	if !tree.Verify() {
		t.Fatal("Test_AvlTreeInsert err")
	}

	idx := 0
	verifArr := []int{5, 10, 15, 18, 20, 40, 60, 70, 80}
	iter := NewIterator(tree)
	for iter.Next() {
		if iter.GetKey().(int) != verifArr[idx] {
			t.Fatalf("want %v, got %v\n", idx, iter.GetKey().(int))
		}

		idx++
	}
}

func Test_AvlTreeRandInsert(t *testing.T) {
	tree := NewTree(avltreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{3, 0, 1, 7, 2, 4, 9, 8, 5, 6}
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
		t.Fatal("Test_AvlTreeRandInsert err")
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

func Test_AvlTreeDelete(t *testing.T) {
	tree := NewTree(avltreeComparator{})

	tree.Insert(20, 20)
	tree.Insert(40, 40)
	tree.Insert(60, 60)
	tree.Insert(80, 80)
	tree.Insert(70, 70)
	tree.Insert(15, 15)
	tree.Insert(18, 18)
	tree.Insert(10, 10)
	tree.Insert(5, 5)
	fmt.Println("avl树:")
	fmt.Println(tree)
}

func Test_AvlTreeRandDelete(t *testing.T) {
	tree := NewTree(avltreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{5, 4, 2, 1, 6, 3, 9, 0, 8, 7}
	// fmt.Printf("插入: %v\n", array)
	for _, v := range array {
		tree.Insert(v, v)
	}

	if !tree.Verify() {
		t.Fatalf("Test_AvlTreeRandDelete insert err, array %v\n", array)
	}

	// fmt.Println(tree)

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	// dArray = []int{2, 3, 4, 1, 0}
	// fmt.Printf("删除: %v\n", dArray)
	for _, v := range dArray {
		// fmt.Printf("删除%v\n", v)
		tree.Delete(v)
		// fmt.Println(tree)
	}

	if !tree.Verify() {
		t.Fatalf("Test_AvlTreeRandDelete delete err, array %v\n", dArray)
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

func Test_AvlTreeRandSearch(t *testing.T) {
	tree := NewTree(avltreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(v, v)
	}

	// if !tree.Verify() {
	// 	fmt.Printf("\n插入 - 测试结果是否正确: %v\n", array)
	// } else {
	// 	fmt.Println("插入正确")
	// }

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	node := tree.Search(sKey)
	if node == nil || node.GetKey().(int) != sKey {
		t.Fatalf("查找错误: 查找节点%v", sKey)
	}
}

func Test_AvlTreeSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(avltreeComparator{})

	tree.Insert(20, 20)
	tree.Insert(40, 40)
	tree.Insert(60, 60)
	tree.Insert(80, 80)
	tree.Insert(70, 70)
	tree.Insert(15, 15)
	tree.Insert(18, 18)
	tree.Insert(10, 10)
	tree.Insert(5, 5)
	// fmt.Println(tree)

	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(7, 5)
	// for i, node := range nodeList {
	// 	fmt.Printf("第%v个节点: \t%v\n", i+1, node.GetKey())
	// }

	verifArr := []int{10, 15, 18, 20, 40}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatalf("Test_AvlTreeSearchRangeLowerBoundKeyWithLimit err: v.GetKey().(int) != verifArr[%d], v.GetKey().(int): %v, verifArr[%d]: %v\n", i, v.GetKey().(int), i, verifArr[i])
		}
	}

	if len(nodeList) != len(verifArr) {
		t.Fatalf("Test_AvlTreeSearchRangeLowerBoundKeyWithLimit err: len(nodeList) != len(verifArr), len(nodeList): %v, len(verifArr): %v\n", len(nodeList), len(verifArr))
	}
}

func Test_AvlTreeRandSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(avltreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(v, v)
	}

	// if !tree.Verify() {
	// 	fmt.Printf("\n插入 - 测试结果是否正确: %v\n", array)
	// } else {
	// 	fmt.Println("插入正确")
	// }

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey+i != node.GetKey() {
			t.Fatal("Test_AvlTreeRandSearchRangeLowerBoundKeyWithLimit err")
		}
	}
}

func Test_AvlTreeSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(avltreeComparator{})

	tree.Insert(20, 20)
	tree.Insert(40, 40)
	tree.Insert(60, 60)
	tree.Insert(80, 80)
	tree.Insert(70, 70)
	tree.Insert(15, 15)
	tree.Insert(18, 18)
	tree.Insert(10, 10)
	tree.Insert(5, 5)

	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(79, 5)
	// for i, node := range nodeList {
	// 	fmt.Printf("第%v个节点: \t%v\n", i+1, node.GetKey())
	// }

	verifArr := []int{18, 20, 40, 60, 70}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatalf("Test_AvlTreeSearchRangeUpperBoundKeyWithLimit err: v.GetKey().(int) != verifArr[%d], v.GetKey().(int): %v, verifArr[%d]: %v\n", i, v.GetKey().(int), i, verifArr[i])
		}
	}

	if len(nodeList) != len(verifArr) {
		t.Fatalf("Test_AvlTreeSearchRangeUpperBoundKeyWithLimit err: len(nodeList) != len(verifArr), len(nodeList): %v, len(verifArr): %v\n", len(nodeList), len(verifArr))
	}
}

func Test_AvlTreeRandSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(avltreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(v, v)
	}

	// if !tree.Verify() {
	// 	fmt.Printf("\n插入 - 测试结果是否正确: %v\n", array)
	// } else {
	// 	fmt.Println("插入正确")
	// }

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey-(len(nodeList)-i-1) != node.GetKey() {
			t.Fatal("Test_AvlTreeRandSearchRangeUpperBoundKeyWithLimit err")
		}
	}
}

func Test_AvlTreeSearchRange(t *testing.T) {
	tree := NewTree(avltreeComparator{})

	tree.Insert(20, 20)
	tree.Insert(40, 40)
	tree.Insert(60, 60)
	tree.Insert(80, 80)
	tree.Insert(70, 70)
	tree.Insert(15, 15)
	tree.Insert(18, 18)
	tree.Insert(10, 10)
	tree.Insert(5, 5)

	nodeList := tree.SearchRange(50, 90)
	// for i, node := range nodeList {
	// 	fmt.Printf("第%v个节点: \t%v\n", i+1, node.GetKey())
	// }

	verifArr := []int{60, 70, 80}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatalf("Test_AvlTreeSearchRange err: v.GetKey().(int) != verifArr[%d],  v.GetKey().(int): %v, verifArr[%d]: %v\n", i, v.GetKey().(int), i, verifArr[i])
		}
	}

	if len(nodeList) != len(verifArr) {
		t.Fatalf("Test_AvlTreeSearchRange err: len(nodeList) != len(verifArr), len(nodeList): %v, len(verifArr): %v\n", len(nodeList), len(verifArr))
	}
}
