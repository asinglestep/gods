package btree

import (
	// "fmt"

	"github.com/asinglestep/gods/utils"
)

// Iterator Iterator
type Iterator struct {
	tree     *Tree
	node     *TreeNode
	entry    *utils.Entry
	entryPos int
	index    int
}

// NewIterator NewIterator
func NewIterator(tree *Tree) *Iterator {
	iter := &Iterator{}
	iter.tree = tree
	iter.node = nil
	iter.entry = nil
	iter.entryPos = -1
	iter.index = -1

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
		iter.entryPos++
		iter.entry = iter.node.entries[iter.entryPos]
	} else if iter.entryPos+1 < len(iter.node.entries) {
		iter.entryPos++
		iter.entry = iter.node.entries[iter.entryPos]
	} else {
		iter.next()
	}

	return true
}

func (iter *Iterator) next() {
	parent := iter.node.parent

	for parent != nil {
		pos := parent.findKeyPosition(iter.tree.comparator, iter.entry.GetKey())
		if pos+1 < len(parent.childrens) {
			iter.entry = parent.entries[pos]
			iter.node = parent.childrens[pos+1].minimum()
			iter.entryPos = -1
			return
		}

		parent = parent.parent
	}
}

// GetKey GetKey
func (iter *Iterator) GetKey() interface{} {
	return iter.entry.GetKey()
}

// GetValue GetValue
func (iter *Iterator) GetValue() interface{} {
	return iter.entry.GetValue()
}
