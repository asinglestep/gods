package avltree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_AvlTreeInsert(t *testing.T) {
	tree := NewTree()

	tree.Insert(20)
	tree.Insert(40)
	tree.Insert(60) // test case 2.3
	tree.Insert(80)
	tree.Insert(70) // test case 2.4
	tree.Insert(15)
	tree.Insert(18) // test case 2.2
	tree.Insert(10)
	tree.Insert(5) // test case 2.1

	fmt.Printf("插入测试\n")
	fmt.Printf("\n输入数组: \n")
	fmt.Println("[20, 40, 60, 80, 70, 15, 18, 10, 5]")
	// fmt.Printf("\n中序遍历结果: \n")
	// tree.PrintAvlTree()
	if err := tree.Dot(); err != nil {
		fmt.Printf("Dot Error %v\n", err)
	}
	fmt.Printf("\n测试结果是否正确: %v\n", tree.VerifAvlTree())
}

func Test_AvlTreeRandInsert(t *testing.T) {
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
	// tree.PrintAvlTree()
	fmt.Printf("\n测试结果是否正确: %v\n", tree.VerifAvlTree())

}

func Test_AvlTreeDelete(t *testing.T) {
	tree := NewTree()

	tree.Insert(20)
	tree.Insert(40)
	tree.Insert(60)
	tree.Insert(80)
	tree.Insert(70)
	tree.Insert(15)
	tree.Insert(18)
	tree.Insert(10)
	tree.Insert(5)
	tree.PrintAvlTree()

	tree.Delete(15)
	fmt.Println("\n删除15")
	tree.PrintAvlTree()
	tree.Delete(20) // test case 2.1.2
	fmt.Println("\n删除20")
	tree.PrintAvlTree()

	tree.Delete(5)
	fmt.Println("\n删除5")
	tree.PrintAvlTree()
	tree.Delete(18)
	fmt.Println("\n删除18")
	tree.PrintAvlTree()
	tree.Delete(10) // test case 2.2.2
	fmt.Println("\n删除10")
	tree.PrintAvlTree()

	tree.Delete(80) // test case 2.1.1
	fmt.Println("\n删除80")
	tree.PrintAvlTree()

	tree.Insert(65)
	fmt.Println("\n插入65")
	tree.PrintAvlTree()

	tree.Delete(40) // test case 2.2.1
	fmt.Println("\n删除40")
	tree.PrintAvlTree()

	tree.Delete(65)
	fmt.Println("\n删除65")
	tree.PrintAvlTree()

	tree.Delete(70)
	fmt.Println("\n删除70")
	tree.PrintAvlTree()

	tree.Delete(60)
	fmt.Println("\n删除60")
	tree.PrintAvlTree()
}

func Test_AvlTreeRandDelete(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(Key(v))
	}

	if !tree.VerifAvlTree() {
		fmt.Printf("\n插入 - 测试结果是否正确: %v\n", array)
	} else {
		t.Log("插入正确")
	}

	// tree.PrintAvlTree()

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	// dArray = []int{0, 3, 1, 4, 2}

	for _, v := range dArray {
		tree.Delete(Key(v))
		// fmt.Printf("删除%v\n", v)
		// tree.PrintAvlTree()
	}

	if !tree.VerifAvlTree() {
		fmt.Printf("\n删除错误: iArray %v\ndArray %v\n", array, dArray)
	} else {
		t.Log("删除正确")
	}
}

func Test_AvlTreeSearch(t *testing.T) {
	tree := NewTree()

	tree.Insert(20)
	tree.Insert(40)
	tree.Insert(60)
	tree.Insert(80)
	tree.Insert(70)
	tree.Insert(15)
	tree.Insert(18)
	tree.Insert(10)
	tree.Insert(5)

	node := tree.Search(10)
	if node == nil || node.key != 10 {
		t.Fatalf("查找错误: 查找节点%v\n", 10)
	}

	node = tree.Search(11)
	if node != nil {
		t.Fatalf("查找错误: 查找节点%v\n", 11)
	}
}

func Test_AvlTreeRandSearch(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(Key(v))
	}

	if !tree.VerifAvlTree() {
		fmt.Printf("\n插入 - 测试结果是否正确: %v\n", array)
	} else {
		t.Log("插入正确")
	}

	sKey := Key(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num))
	node := tree.Search(sKey)
	if node == nil || node.key != sKey {
		t.Fatalf("查找错误: 查找节点%v", sKey)
	}
}

func Test_AvlTreeSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree()

	tree.Insert(20)
	tree.Insert(40)
	tree.Insert(60)
	tree.Insert(80)
	tree.Insert(70)
	tree.Insert(15)
	tree.Insert(18)
	tree.Insert(10)
	tree.Insert(5)
	// tree.PrintAvlTree()

	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(7, 5)
	for i, node := range nodeList {
		fmt.Printf("第%v个节点: \t%v\n", i+1, node.key)
	}
}

func Test_AvlTreeRandSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(Key(v))
	}

	if !tree.VerifAvlTree() {
		fmt.Printf("\n插入 - 测试结果是否正确: %v\n", array)
	} else {
		t.Log("插入正确")
	}

	sKey := Key(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num))
	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey+Key(i) != node.key {
			t.Fatalf("sKey+Key(i) != node.key")
		}
	}
}

func Test_AvlTreeSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree()

	tree.Insert(20)
	tree.Insert(40)
	tree.Insert(60)
	tree.Insert(80)
	tree.Insert(70)
	tree.Insert(15)
	tree.Insert(18)
	tree.Insert(10)
	tree.Insert(5)

	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(79, 5)
	for i, node := range nodeList {
		fmt.Printf("第%v个节点: \t%v\n", i+1, node.key)
	}
}

func Test_AvlTreeRandSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{1, 3, 2, 7, 5, 8, 4, 6, 9, 0}
	for _, v := range array {
		tree.Insert(Key(v))
	}

	if !tree.VerifAvlTree() {
		fmt.Printf("\n插入 - 测试结果是否正确: %v\n", array)
	} else {
		t.Log("插入正确")
	}

	sKey := Key(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num))
	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey-Key(len(nodeList)-i-1) != node.key {
			t.Fatalf("sKey+Key(len(nodeList)-i) != node.key ")
		}
	}
}

func Test_AvlTreeSearchRange(t *testing.T) {
	tree := NewTree()

	tree.Insert(20)
	tree.Insert(40)
	tree.Insert(60)
	tree.Insert(80)
	tree.Insert(70)
	tree.Insert(15)
	tree.Insert(18)
	tree.Insert(10)
	tree.Insert(5)

	nodeList := tree.SearchRange(50, 90)
	for i, node := range nodeList {
		fmt.Printf("第%v个节点: \t%v\n", i+1, node.key)
	}
}
