package skiplist

import (
	"math/rand"
	"testing"
	"time"
)

type intComparator struct {
}

func (c intComparator) Compare(k1, k2 interface{}) int {
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

func Test_SkipListRandInsert(t *testing.T) {
	num := 10
	arr := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// fmt.Println("插入数组: ", arr)
	// arr = []int{2, 9, 4, 5, 8, 1, 0, 6, 7, 3}

	skipList := NewList(intComparator{})
	for _, v := range arr {
		skipList.Insert(v, v)
	}

	// fmt.Println(skipList)

	idx := 0
	iter := NewIterator(skipList)
	for iter.Next() {
		if iter.node.entry.key.(int) != idx {
			t.Fatalf("want %v, got %v\n", idx, iter.GetKey().(int))
		}

		idx++
	}
}

func Test_SkipListRandDelete(t *testing.T) {
	num := 10
	arr := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// fmt.Println("插入数组: ", arr)
	// arr = []int{2, 9, 4, 5, 8, 1, 0, 6, 7, 3}

	skipList := NewList(intComparator{})
	for _, v := range arr {
		skipList.Insert(v, v)
	}

	// fmt.Println("插入结果: ")
	// fmt.Println(skipList)

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	// fmt.Println("删除数组: ", dArray)
	// dArray = []int{4, 2, 3, 0, 1}

	for _, v := range dArray {
		skipList.Delete(v)
		// fmt.Printf("\n删除%v: \n", v+1)
		// fmt.Println(skipList)
	}

	// fmt.Println("删除结果: ")
	// fmt.Println(skipList)

	idx := 5
	iter := NewIterator(skipList)
	for iter.Next() {
		if iter.node.entry.key.(int) != idx {
			t.Fatalf("want %v, got %v\n", idx, iter.GetKey().(int))
		}

		idx++
	}
}

func Test_SkipListRandSearch(t *testing.T) {
	num := 10
	arr := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// fmt.Println("插入数组: ", arr)
	// arr = []int{2, 9, 4, 5, 8, 1, 0, 6, 7, 3}

	skipList := NewList(intComparator{})
	for _, v := range arr {
		skipList.Insert(v, v)
	}

	// fmt.Println("插入结果: ")
	// fmt.Println(skipList)

	key := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	// fmt.Printf("查找的key: %v\n", key)

	entry := skipList.Search(key)
	if entry == nil {
		t.Fatalf("want %v, got nil\n", key)
	}

	if entry.GetKey().(int) != key {
		t.Fatalf("want %v, got %v\n", key, entry.GetKey().(int))
	}
}
