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

	// fmt.Printf("\n插入结果是否正确: %v\n", tree.Verify())

	if !tree.Verify() {
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

	// fmt.Printf("\n删除结果是否正确: %v\n", tree.Verify())
	// tree.Dot()

	if !tree.Verify() {
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

func Test_TreapSearchRange(t *testing.T) {
	tree := NewTree(treapComparator{})

	tree.Insert(30, 30)
	tree.Insert(20, 20)
	tree.Insert(10, 10)
	tree.Insert(40, 40)
	tree.Insert(50, 50)
	tree.Insert(100, 100)
	tree.Insert(60, 60)
	tree.Insert(80, 80)
	tree.Insert(90, 90)

	// fmt.Println(tree)
	nodeList := tree.SearchRange(70, 110)
	// for i, v := range nodeList {
	// 	fmt.Printf("第%v个节点 -- key: %v\n", i+1, v.key)
	// }

	verifArr := []int{80, 90, 100}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatalf("Test_TreapSearchRange err: v.GetKey().(int) != verifArr[%d], v.GetKey().(int): %v, verifArr[%d]: %v\n", i, v.GetKey().(int), i, verifArr[i])
		}
	}

	if len(nodeList) != len(verifArr) {
		t.Fatalf("Test_TreapSearchRange err: len(nodeList) != len(verifArr), len(nodeList): %v, len(verifArr): %v\n", len(nodeList), len(verifArr))
	}
}

func Test_TreapSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(treapComparator{})

	tree.Insert(30, 30)
	tree.Insert(20, 20)
	tree.Insert(10, 10)
	tree.Insert(40, 40)
	tree.Insert(50, 50)
	tree.Insert(100, 100)
	tree.Insert(60, 60)
	tree.Insert(80, 80)
	tree.Insert(90, 90)

	nodeList := tree.SearchRangeLowerBoundKeyWithLimit(65, 2)
	// for i, node := range nodeList {
	// 	fmt.Printf("第%v个节点: \t%v\n", i+1, node.key)
	// }

	verifArr := []int{80, 90}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatalf("Test_TreapSearchRangeLowerBoundKeyWithLimit err: v.GetKey().(int) != verifArr[%d], v.GetKey().(int): %v, verifArr[%d]: %v\n", i, v.GetKey().(int), i, verifArr[i])
		}
	}

	if len(nodeList) != len(verifArr) {
		t.Fatalf("Test_TreapSearchRangeLowerBoundKeyWithLimit err: len(nodeList) != len(verifArr), len(nodeList): %v, len(verifArr): %v\n", len(nodeList), len(verifArr))
	}
}

func Test_TreapRandSearchRangeLowerBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(treapComparator{})
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
			t.Fatal("Test_TreapRandSearchRangeLowerBoundKeyWithLimit err")
		}
	}
}

func Test_TreapSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(treapComparator{})

	tree.Insert(30, 30)
	tree.Insert(20, 20)
	tree.Insert(10, 10)
	tree.Insert(40, 40)
	tree.Insert(50, 50)
	tree.Insert(100, 100)
	tree.Insert(60, 60)
	tree.Insert(80, 80)
	tree.Insert(90, 90)

	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(35, 5)
	// for i, node := range nodeList {
	// 	fmt.Printf("第%v个节点: \t%v\n", i+1, node.key)
	// }

	verifArr := []int{10, 20, 30}
	for i, v := range nodeList {
		if v.GetKey().(int) != verifArr[i] {
			t.Fatalf("Test_TreapSearchRangeUpperBoundKeyWithLimit err: v.GetKey().(int) != verifArr[%d], v.GetKey().(int): %v, verifArr[%d]: %v\n", i, v.GetKey().(int), i, verifArr[i])
		}
	}

	if len(nodeList) != len(verifArr) {
		t.Fatalf("Test_TreapSearchRangeUpperBoundKeyWithLimit err: len(nodeList) != len(verifArr), len(nodeList): %v, len(verifArr): %v\n", len(nodeList), len(verifArr))
	}
}

func Test_TreapRandSearchRangeUpperBoundKeyWithLimit(t *testing.T) {
	tree := NewTree(treapComparator{})
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
	// if !tree.Verify() {
	// 	t.Fatal("SearchRangeUpperBoundKeyWithLimit err")
	// }

	sKey := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	nodeList := tree.SearchRangeUpperBoundKeyWithLimit(sKey, 1000)
	for i, node := range nodeList {
		if sKey-(len(nodeList)-i-1) != node.GetKey().(int) {
			t.Fatal("Test_TreapRandSearchRangeUpperBoundKeyWithLimit err: 查找结果错误")
		}
	}
}
