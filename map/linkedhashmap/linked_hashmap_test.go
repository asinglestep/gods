package linkedhashmap

import (
	"testing"

	"github.com/asinglestep/gods/utils"
)

func Test_Put(t *testing.T) {
	lmap := NewLinkedHashMap(10)
	lmap.Put(utils.NewEntry(1, 1))
	lmap.Put(utils.NewEntry(2, 2))
	lmap.Put(utils.NewEntry(3, 3))
	lmap.Put(utils.NewEntry(4, 4))
	lmap.Put(utils.NewEntry(5, 5))

	lmap.Put(utils.NewEntry(6, 6))
	lmap.Put(utils.NewEntry(7, 7))
	lmap.Put(utils.NewEntry(8, 8))
	lmap.Put(utils.NewEntry(9, 9))
	lmap.Put(utils.NewEntry(10, 10))

	lmap.Put(utils.NewEntry(11, 11))

	iter := NewIterator(lmap)

	idx := 2
	for iter.Next() {
		e, ok := iter.Entry().(*utils.Entry)
		if !ok {
			t.Fatalf("%v", ErrEntryType)
		}

		if e.GetKey().(int) != idx {
			t.Fatalf("wang %v, got %v", idx, e.GetKey().(int))
		}

		idx++
	}
}

func Test_Get(t *testing.T) {
	lmap := NewLinkedHashMap(10)
	lmap.Put(utils.NewEntry(1, 1))
	lmap.Put(utils.NewEntry(2, 2))
	lmap.Put(utils.NewEntry(3, 3))
	lmap.Put(utils.NewEntry(4, 4))
	lmap.Put(utils.NewEntry(5, 5))

	lmap.Put(utils.NewEntry(6, 6))
	lmap.Put(utils.NewEntry(7, 7))
	lmap.Put(utils.NewEntry(8, 8))
	lmap.Put(utils.NewEntry(9, 9))
	lmap.Put(utils.NewEntry(10, 10))

	lmap.Put(utils.NewEntry(11, 11))

	if _, err := lmap.Get(1); err != ErrNotExist {
		t.Fatal("1 exist")
	}

	e, err := lmap.Get(10)
	if err != nil {
		t.Fatalf("get 10 err: %v\n", err)
	}

	if e.GetKey().(int) != 10 {
		t.Fatalf("wang 10, got %v", e.GetKey().(int))
	}
}
