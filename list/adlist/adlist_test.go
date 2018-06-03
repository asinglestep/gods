package adlist

import (
	"testing"
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

func Test_Insert(t *testing.T) {
	list := NewList(intComparator{})
	list.AddNodeToHead(4)
	list.AddNodeToTail(5)

	list.AddNodeToHead(3)
	list.AddNodeToTail(6)

	list.AddNodeToHead(2)
	list.AddNodeToTail(7)

	list.AddNodeToHead(1)
	list.AddNodeToTail(8)

	var idx = 1
	iter := NewIterator(list)
	for iter.Next() {
		if iter.Entry().(int) != idx {
			t.Fatalf("want %v, got %v\n", idx, iter.Entry().(int))
		}

		idx++
	}
}

func Test_Delete(t *testing.T) {
	list := NewList(intComparator{})

	list.AddNodeToHead(4)
	list.AddNodeToTail(5)

	list.AddNodeToHead(3)
	list.AddNodeToTail(6)

	list.AddNodeToHead(2)
	list.AddNodeToTail(7)

	list.AddNodeToHead(1)
	list.AddNodeToTail(8)

	list.DeleteNode(list.Head())
	list.DeleteNode(list.Tail())

	var idx = 2
	iter := NewIterator(list)
	for iter.Next() {
		if iter.Entry().(int) != idx {
			t.Fatalf("want %v, got %v\n", idx, iter.Entry().(int))
		}

		idx++
	}
}

func Test_Search(t *testing.T) {
	list := NewList(intComparator{})

	list.AddNodeToHead(4)
	list.AddNodeToTail(5)

	list.AddNodeToHead(3)
	list.AddNodeToTail(6)

	list.AddNodeToHead(2)
	list.AddNodeToTail(7)

	list.AddNodeToHead(1)
	list.AddNodeToTail(8)

	if list.SearchNode(7).GetEntry().(int) != 7 {
		t.Fatalf("want 7, got %v\n", list.SearchNode(7).GetEntry().(int))
	}
}
