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
	tree.PrintRbTree()

	fmt.Println("\n插入2")
	tree.Insert(2, 2)
	tree.PrintRbTree()

	fmt.Println("\n插入1")
	tree.Insert(1, 1) // test case 2.4
	tree.PrintRbTree()

	fmt.Println("\n插入4")
	tree.Insert(4, 4) // test case 2.3
	tree.PrintRbTree()

	fmt.Println("\n插入5")
	tree.Insert(5, 5) // test case 2.5
	tree.PrintRbTree()

	fmt.Println("\n插入10")
	tree.Insert(10, 10)
	tree.PrintRbTree()

	fmt.Println("\n插入6")
	tree.Insert(6, 6) // test case 2.6
	tree.PrintRbTree()

	fmt.Println("\n插入8")
	tree.Insert(8, 8)
	tree.PrintRbTree()

	fmt.Println("\n插入9")
	tree.Insert(9, 9) // test case 2.7
	tree.PrintRbTree()

	// fmt.Printf("插入测试\n")
	// fmt.Printf("\n输入数组: \n")
	// fmt.Println("[3, 2, 1, 4, 5, 10, 6, 8, 9]")
	// fmt.Printf("\n中序遍历结果: \n")
	// tree.PrintRbTree()
	// fmt.Printf("\n测试结果是否正确: %v\n", tree.VerifRbTree())

	// if err := tree.Dot(); err != nil {
	// 	fmt.Printf("Dot error %v\n", err)
	// }

	if !tree.VerifRbTree() {
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
	// tree.PrintRbTree()
	// fmt.Printf("\n测试结果是否正确: %v\n", tree.VerifRbTree())
	if !tree.VerifRbTree() {
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
	tree.PrintRbTree()

	fmt.Println("\n删除9")
	tree.Delete(9) // test case 2.1
	tree.PrintRbTree()

	fmt.Println("\n删除5")
	tree.Delete(5) // test case 2.2.3
	tree.PrintRbTree()

	fmt.Println("\n删除6")
	tree.Delete(6) // test case 2.2.1
	tree.PrintRbTree()

	fmt.Println("\n删除10")
	tree.Delete(10)
	tree.PrintRbTree()

	fmt.Println("\n删除8")
	tree.Delete(8) // test case 2.3.4, case 2.3.1
	tree.PrintRbTree()

	fmt.Println("\n插入10")
	tree.Insert(10, 10)
	tree.PrintRbTree()

	fmt.Println("\n插入8")
	tree.Insert(8, 8)
	tree.PrintRbTree()

	fmt.Println("\n删除1")
	tree.Delete(1) // test case 2.2.4
	tree.PrintRbTree()

	fmt.Println("\n删除10")
	tree.Delete(10)
	tree.PrintRbTree()

	fmt.Println("\n删除8")
	tree.Delete(8) // test case 2.3.3
	tree.PrintRbTree()

	fmt.Println("\n插入1")
	tree.Insert(1, 1)
	tree.PrintRbTree()

	fmt.Println("\n删除4")
	tree.Delete(4) // test case 2.3.2
	tree.PrintRbTree()

	fmt.Println("\n删除3")
	tree.Delete(3)
	tree.PrintRbTree()

	fmt.Println("\n插入4")
	tree.Insert(4, 4)
	tree.PrintRbTree()

	fmt.Println("\n插入3")
	tree.Insert(3, 3)
	tree.PrintRbTree()

	fmt.Println("\n删除1")
	tree.Delete(1) // test case 2.2.2
	tree.PrintRbTree()

	fmt.Println("\n删除3")
	tree.Delete(3)
	tree.PrintRbTree()

	fmt.Println("\n删除4")
	tree.Delete(4)
	tree.PrintRbTree()

	fmt.Println("\n删除2")
	tree.Delete(2)
	tree.PrintRbTree()
}

func Test_RbTreeRandDelete(t *testing.T) {
	tree := NewTree(rbtreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}

	// tree.PrintRbTree()

	// if !tree.VerifRbTree() {
	// 	fmt.Printf("插入错误, array %v\n", array)
	// } else {
	// 	fmt.Println("插入正确")
	// }

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)

	for _, v := range dArray {
		tree.Delete(v)
	}

	// fmt.Println("")
	// tree.PrintRbTree()

	// if !tree.VerifRbTree() {
	// 	fmt.Printf("删除错误, array %v, dArray %v\n", array, dArray)
	// } else {
	// 	fmt.Printf("删除正确\n")
	// }

	if !tree.VerifRbTree() {
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

	// tree.PrintRbTree()
	nodeList := tree.SearchRange(7, 11)
	// for i, v := range nodeList {
	// 	fmt.Printf("第%v个节点 -- key: %v\n", i+1, v.key)
	// }

	verifArr := []int{8, 9, 10}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatal("Test_RbTreeSearchRange err")
		}
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
			t.Fatal("Test_RbTreeSearchRangeLowerBoundKeyWithLimit err")
		}
	}
}

func Test_RbTreeRandSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(rbtreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}

	// tree.PrintRbTree()

	// if !tree.VerifRbTree() {
	// 	fmt.Printf("插入错误, array %v\n", array)
	// } else {
	// 	fmt.Println("插入正确")
	// }

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey+i != node.GetKey().(int) {
			t.Fatal("SearchRangeLowerBoundKeyWithLimit err")
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
			t.Fatal("Test_RbTreeSearchRangeUpperBoundKeyWithLimit err")
		}
	}
}

func Test_RbTreeRandSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(rbtreeComparator{})
	var num = 1000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range array {
		tree.Insert(v, v)
	}

	// tree.PrintRbTree()

	if !tree.VerifRbTree() {
		fmt.Printf("插入错误, array %v\n", array)
	} else {
		fmt.Println("插入正确")
	}

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey-(len(nodeList)-i-1) != node.GetKey().(int) {
			t.Fatal("SearchRangeUpperBoundKeyWithLimit err")
		}
	}
}
