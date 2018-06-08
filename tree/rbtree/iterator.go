package rbtree

// Iterator Iterator
type Iterator struct {
	index int
	tree  *Tree
	node  *TreeNode
}

// NewIterator NewIterator
func NewIterator(tree *Tree) *Iterator {
	iter := &Iterator{}
	iter.index = -1
	iter.tree = tree
	iter.node = nil

	return iter
}

// Next Next
func (iter *Iterator) Next() bool {
	iter.index++
	if iter.index >= iter.tree.size {
		return false
	}

	if iter.index == 0 {
		iter.node = iter.tree.minimum()
	} else {
		iter.node = iter.node.next()
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