package adlist

// Iterator Iterator
type Iterator struct {
	index int
	list  *List
	node  *Node
}

// NewIterator NewIterator
func NewIterator(list *List) *Iterator {
	iter := &Iterator{}
	iter.index = -1
	iter.list = list
	iter.node = nil

	return iter
}

// Next Next
func (iter *Iterator) Next() bool {
	iter.index++
	if iter.index >= iter.list.length {
		return false
	}

	if iter.index == 0 {
		iter.node = iter.list.head
	} else {
		iter.node = iter.node.next
	}

	return true
}

// Entry Entry
func (iter *Iterator) Entry() interface{} {
	return iter.node.entry
}
