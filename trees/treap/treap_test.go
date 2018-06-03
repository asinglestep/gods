package treap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_TreapRandInsert(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{8, 4, 0, 6, 3, 9, 1, 5, 2, 7}

	for _, v := range array {
		tree.Insert(Key(v))
	}

	fmt.Printf("\n随机插入%v个数\n", num)
	// fmt.Printf("随机数组: %v\n", array)

	fmt.Printf("\n插入结果是否正确: %v\n", tree.VerifTreap())

	// tree.Dot()
}

func Test_TreapRandDelete(t *testing.T) {
	tree := NewTree()
	var num = 10000000

	array := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num)
	// array = []int{8, 4, 0, 6, 3, 9, 1, 5, 2, 7}
	// fmt.Printf("随机插入数组: %v\n", array)

	for _, v := range array {
		tree.Insert(Key(v))
	}

	delArray := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(num / 2)
	// fmt.Printf("随机删除数组: %v\n", delArray)

	for _, v := range delArray {
		tree.Delete(Key(v))
	}

	fmt.Printf("\n删除结果是否正确: %v\n", tree.VerifTreap())
	// tree.Dot()
}
