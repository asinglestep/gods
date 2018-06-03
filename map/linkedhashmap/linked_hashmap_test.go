package linkedhashmap

import (
	"testing"
)

func Test_Put(t *testing.T) {
	lmap := NewLinkedHashMap(10)
	lmap.Put(Entry{"1", "1"})
	lmap.Put(Entry{"2", "2"})
	lmap.Put(Entry{"3", "3"})
	lmap.Put(Entry{"4", "4"})
	lmap.Put(Entry{"5", "5"})

	lmap.Put(Entry{"6", "6"})
	lmap.Put(Entry{"7", "7"})
	lmap.Put(Entry{"8", "8"})
	lmap.Put(Entry{"9", "9"})
	lmap.Put(Entry{"10", "10"})

	lmap.Put(Entry{"11", "11"})

	iter := NewIterator(lmap)
	for iter.Next() {
		t.Logf("%#v\n", iter.Entry())
	}
}

func Test_Get(t *testing.T) {
	lmap := NewLinkedHashMap(10)
	lmap.Put(Entry{"1", "1"})
	lmap.Put(Entry{"2", "2"})
	lmap.Put(Entry{"3", "3"})
	lmap.Put(Entry{"4", "4"})
	lmap.Put(Entry{"5", "5"})

	lmap.Put(Entry{"6", "6"})
	lmap.Put(Entry{"7", "7"})
	lmap.Put(Entry{"8", "8"})
	lmap.Put(Entry{"9", "9"})
	lmap.Put(Entry{"10", "10"})

	lmap.Put(Entry{"11", "11"})

	if _, err := lmap.Get("1"); err != errNotExist {
		t.Fatal("get 1")
	}

	e, err := lmap.Get("10")
	if err != nil {
		t.Fatal("get 10")
	}

	t.Log(e)
}
