package bptree

import (
	"github.com/asinglestep/gods/utils"
)

// Iterator Iterator
type Iterator struct {
	tree     *Tree
	leaf     *TreeLeaf
	entry    *utils.Entry
	entryPos int
}

// NewIterator NewIterator
func NewIterator(tree *Tree) *Iterator {
	return NewIteratorWithLeaf(tree, tree.minimum(), 0)
}

// NewIteratorWithLeaf NewIteratorWithLeaf
func NewIteratorWithLeaf(tree *Tree, leaf *TreeLeaf, pos int) *Iterator {
	iter := &Iterator{}
	iter.tree = tree
	iter.leaf = leaf
	iter.entry = nil
	iter.entryPos = pos

	return iter
}

// Next Next
func (iter *Iterator) Next() bool {
	if iter.leaf == nil {
		return false
	}

	if iter.entryPos < len(iter.leaf.entries) {
		iter.entry = iter.leaf.entries[iter.entryPos]
		iter.entryPos++
		return true
	}

	iter.leaf = iter.leaf.next
	if iter.leaf != nil {
		iter.entry = iter.leaf.entries[0]
		iter.entryPos = 1
		return true
	}

	return false
}

// GetKey GetKey
func (iter *Iterator) GetKey() interface{} {
	return iter.entry.GetKey()
}

// GetValue GetValue
func (iter *Iterator) GetValue() interface{} {
	return iter.entry.GetValue()
}
