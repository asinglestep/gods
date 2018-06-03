package skiplist

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_randomLevel(t *testing.T) {
	for i := 0; i < 10000; i++ {
		if level := randomLevel(); level >= 10 {
			t.Log(level)
		}
	}
}

func Test_SkipListRandInsert(t *testing.T) {
	num := 10
	arr := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	skipList := NewSkipList()
	for _, v := range arr {
		skipList.Insert(Key(v+1), Value(v+1))
	}

	skipList.PrintSkipList()
}

func Test_SkipListRandDelete(t *testing.T) {
	num := 10
	skipList := NewSkipList()

	arr := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// fmt.Println("插入数组: ", arr)
	// arr = []int{2, 9, 4, 5, 8, 1, 0, 6, 7, 3}

	for _, v := range arr {
		skipList.Insert(Key(v+1), Value(v+1))
	}

	fmt.Println("插入结果: ")
	skipList.PrintSkipList()

	dArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	// fmt.Println("删除数组: ", dArray)
	// dArray = []int{4, 2, 3, 0, 1}

	for _, v := range dArray {
		skipList.Delete(Key(v + 1))
		fmt.Printf("\n删除%v: \n", v+1)
		skipList.PrintSkipList()
	}

	fmt.Println("删除结果: ")
	skipList.PrintSkipList()
}

func Test_SkipListRandSearch(t *testing.T) {
	num := 10
	skipList := NewSkipList()

	arr := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)

	for _, v := range arr {
		skipList.Insert(Key(v+1), Value(v+1))
	}

	// fmt.Println("插入结果: ")
	// skipList.PrintSkipList()

	key := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
	fmt.Printf("查找的key: %v\n", key)

	entry := skipList.Search(Key(key))
	if entry == nil {
		t.Fatal("SkipListRandSearch err")
	}

	t.Logf("entry: %#v\n", entry)
}
