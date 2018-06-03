package rbtree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_RbTreeInsert(t *testing.T) {
	tree := NewTree()
	tree.Insert(3)
	tree.Insert(2)
	tree.Insert(1) // test case 2.4
	tree.Insert(4) // test case 2.3
	tree.Insert(5) // test case 2.5
	tree.Insert(10)
	tree.Insert(6) // test case 2.6
	tree.Insert(8)
	tree.Insert(9) // test case 2.7

	fmt.Printf("插入测试\n")
	fmt.Printf("\n输入数组: \n")
	fmt.Println("[3, 2, 1, 4, 5, 10, 6, 8, 9]")
	fmt.Printf("\n中序遍历结果: \n")
	tree.PrintRbTree()
	fmt.Printf("\n测试结果是否正确: %v\n", tree.VerifRbTree())

	if err := tree.Dot(); err != nil {
		fmt.Printf("Dot error %v\n", err)
	}
}

func Test_RbTreeRandInsert(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(Key(v))
	}

	fmt.Printf("\n随机插入%v个数\n", num)
	// fmt.Printf("\n输入数组: \n")
	// fmt.Println(array)
	// fmt.Printf("\n中序遍历结果: \n")
	// tree.PrintRbTree()
	fmt.Printf("\n测试结果是否正确: %v\n", tree.VerifRbTree())
}

func Test_RbTreeDelete(t *testing.T) {
	tree := NewTree()
	tree.Insert(3)
	tree.Insert(2)
	tree.Insert(1)
	tree.Insert(4)
	tree.Insert(5)
	tree.Insert(10)
	tree.Insert(6)
	tree.Insert(8)
	tree.Insert(9)
	tree.PrintRbTree()

	tree.Delete(9) // test case 2.1
	fmt.Println("\n删除9")
	tree.PrintRbTree()
	tree.Delete(5) // test case 2.2.3
	fmt.Println("\n删除5")
	tree.PrintRbTree()

	tree.Delete(6) // test case 2.2.1
	fmt.Println("\n删除6")
	tree.PrintRbTree()
	tree.Delete(10)
	fmt.Println("\n删除10")
	tree.PrintRbTree()
	tree.Delete(8) // test case 2.3.4, case 2.3.1
	fmt.Println("\n删除8")
	tree.PrintRbTree()

	tree.Insert(10)
	fmt.Println("\n插入10")
	tree.PrintRbTree()
	tree.Insert(8)
	fmt.Println("\n插入8")
	tree.PrintRbTree()

	tree.Delete(1) // test case 2.2.4
	fmt.Println("\n删除1")
	tree.PrintRbTree()

	tree.Delete(10)
	fmt.Println("\n删除10")
	tree.PrintRbTree()

	tree.Delete(8) // test case 2.3.3
	fmt.Println("\n删除8")
	tree.PrintRbTree()

	tree.Insert(1)
	fmt.Println("\n插入1")
	tree.PrintRbTree()

	tree.Delete(4) // test case 2.3.2
	fmt.Println("\n删除4")
	tree.PrintRbTree()

	tree.Delete(3)
	fmt.Println("\n删除3")
	tree.PrintRbTree()

	tree.Insert(4)
	fmt.Println("\n插入4")
	tree.PrintRbTree()
	tree.Insert(3)
	fmt.Println("\n插入3")
	tree.PrintRbTree()

	tree.Delete(1) // test case 2.2.2
	fmt.Println("\n删除1")
	tree.PrintRbTree()

	tree.Delete(3)
	fmt.Println("\n删除3")
	tree.PrintRbTree()

	tree.Delete(4)
	fmt.Println("\n删除4")
	tree.PrintRbTree()

	tree.Delete(2)
	fmt.Println("\n删除2")
	tree.PrintRbTree()
}

func Test_RbTreeRandDelete(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(Key(v))
	}

	// tree.PrintRbTree()

	if !tree.VerifRbTree() {
		fmt.Printf("插入错误, array %v\n", array)
	} else {
		fmt.Printf("插入正确\n")
	}

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)

	for _, v := range dArray {
		tree.Delete(Key(v))
	}

	// fmt.Println("")
	// tree.PrintRbTree()

	if !tree.VerifRbTree() {
		fmt.Printf("删除错误, array %v, dArray %v\n", array, dArray)
	} else {
		fmt.Printf("删除正确\n")
	}
}

func Test_RbTreeSearch(t *testing.T) {
	tree := NewTree()

	tree.Insert(3)
	tree.Insert(2)
	tree.Insert(1)
	tree.Insert(4)
	tree.Insert(5)
	tree.Insert(10)
	tree.Insert(6)
	tree.Insert(8)
	tree.Insert(9)

	node := tree.Search(10)
	if node == nil || node.key != 10 {
		t.Fatalf("查找错误: 查找节点%v\n", 10)
	}

	node = tree.Search(11)
	if node != nil {
		t.Fatalf("查找错误: 查找节点%v\n", 11)
	}
}

func Test_RbTreeSearchRange(t *testing.T) {
	tree := NewTree()

	tree.Insert(3)
	tree.Insert(2)
	tree.Insert(1)
	tree.Insert(4)
	tree.Insert(5)
	tree.Insert(10)
	tree.Insert(6)
	tree.Insert(8)
	tree.Insert(9)

	// tree.PrintRbTree()
	nodeList := tree.SearchRange(7, 11)
	for i, v := range nodeList {
		fmt.Printf("第%v个节点 -- key: %v\n", i+1, v.key)
	}
}

func Test_RbTreeSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree()

	tree.Insert(3)
	tree.Insert(2)
	tree.Insert(1)
	tree.Insert(4)
	tree.Insert(5)
	tree.Insert(10)
	tree.Insert(6)
	tree.Insert(8)
	tree.Insert(9)

	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(6, 2)
	for i, node := range nodeList {
		fmt.Printf("第%v个节点: \t%v\n", i+1, node.key)
	}
}

func Test_RbTreeRandSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(Key(v))
	}

	// tree.PrintRbTree()

	if !tree.VerifRbTree() {
		fmt.Printf("插入错误, array %v\n", array)
	} else {
		fmt.Printf("插入正确\n")
	}

	sKey := Key(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num))
	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey+Key(i) != node.key {
			t.Fatalf("SearchRangeLowerBoundKeyWithLimit错误\n")
		}
	}
}

func Test_RbTreeSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree()

	tree.Insert(3)
	tree.Insert(2)
	tree.Insert(1)
	tree.Insert(4)
	tree.Insert(5)
	tree.Insert(10)
	tree.Insert(6)
	tree.Insert(8)
	tree.Insert(9)

	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(3, 5)
	for i, node := range nodeList {
		fmt.Printf("第%v个节点: \t%v\n", i+1, node.key)
	}
}

func Test_RbTreeRandSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(Key(v))
	}

	// tree.PrintRbTree()

	if !tree.VerifRbTree() {
		fmt.Printf("插入错误, array %v\n", array)
	} else {
		fmt.Printf("插入正确\n")
	}

	sKey := Key(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num))
	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey-Key(len(nodeList)-i-1) != node.key {
			t.Fatalf("SearchRangeUpperBoundKeyWithLimit错误\n")
		}
	}
}
