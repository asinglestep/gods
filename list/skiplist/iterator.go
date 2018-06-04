package skiplist

// Iterator Iterator
type Iterator struct {
	index int
	node  *Node
	list  *List
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
		iter.node = iter.list.head.level[0].forward
	} else {
		iter.node = iter.node.level[0].forward
	}

	return true
}

// GetKey GetKey
func (iter *Iterator) GetKey() interface{} {
	return iter.node.entry.GetKey()
}

// GetValue GetValue
func (iter *Iterator) GetValue() interface{} {
	return iter.node.entry.GetValue()
}
