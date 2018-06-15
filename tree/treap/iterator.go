package treap

// Iterator Iterator
type Iterator struct {
	node   *TreeNode
	tree   *Tree
	bBegin bool
}

// NewIterator NewIterator
func NewIterator(tree *Tree) *Iterator {
	iter := &Iterator{}
	iter.tree = tree
	iter.node = iter.tree.minimum()
	iter.bBegin = true

	return iter
}

// NewIteratorWithNode 从指定的node开始迭代
func NewIteratorWithNode(tree *Tree, node *TreeNode) *Iterator {
	iter := &Iterator{}
	iter.tree = tree
	iter.node = node
	iter.bBegin = true

	return iter
}

// Next Next
func (iter *Iterator) Next() bool {
	if iter.bBegin {
		iter.bBegin = false
	} else {
		iter.node = iter.node.next()
	}

	if iter.node != nil {
		return true
	}

	return false
}

// Prev Prev
func (iter *Iterator) Prev() bool {
	if iter.bBegin {
		iter.bBegin = false
	} else {
		iter.node = iter.node.prev()
	}

	if iter.node != nil {
		return true
	}

	return false
}

// GetKey GetKey
func (iter *Iterator) GetKey() interface{} {
	return iter.node.entry.GetKey()
}

// GetValue GetValue
func (iter *Iterator) GetValue() interface{} {
	return iter.node.entry.GetValue()
}
