package linkedhashmap

import (
	"github.com/asinglestep/gods/list/adlist"
)

// Iterator Iterator
type Iterator struct {
	*adlist.Iterator
}

// NewIterator NewIterator
func NewIterator(lmap *LinkedHashMap) *Iterator {
	iter := &Iterator{}
	iter.Iterator = adlist.NewIterator(lmap.list)

	return iter
}

// Next Next
func (iter *Iterator) Next() bool {
	return iter.Iterator.Next()
}

// Entry Entry
func (iter *Iterator) Entry() interface{} {
	return iter.Iterator.Entry()
}
