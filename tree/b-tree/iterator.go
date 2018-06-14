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
}

// NewIterator NewIterator
func NewIterator(tree *Tree) *Iterator {
	iter := &Iterator{}
	iter.tree = tree
	iter.node = iter.tree.minimum()
	iter.entry = nil
	iter.entryPos = 0

	return iter
}

// NewIteratorWithKey 从指定的key开始迭代
func NewIteratorWithKey(tree *Tree, key interface{}) *Iterator {
	iter := &Iterator{}
	iter.tree = tree

	node, pos := tree.lookup(tree.root, key)
	if node.isLeaf() {
		iter.node = node
		iter.entryPos = pos
		if pos != len(iter.node.entries) {
			// 从当前entry开始
			iter.entry = iter.node.entries[pos]
		} else {
			// 从最后一个entry开始
			iter.entry = iter.node.entries[pos-1]
		}
	} else {
		// 从 以当前节点为根节点，中序遍历后，树的最大节点 开始
		iter.node = node.childrens[pos].maximum()
		iter.entryPos = len(iter.node.entries)
		iter.entry = iter.node.entries[len(iter.node.entries)-1]
	}

	return iter
}

// Next Next
func (iter *Iterator) Next() bool {
	if iter.entryPos < len(iter.node.entries) {
		iter.entry = iter.node.entries[iter.entryPos]
		iter.entryPos++
		return true
	}

	parent := iter.node.parent
	for parent != nil {
		pos := parent.findKeyPosition(iter.tree.comparator, iter.entry.GetKey())
		if pos+1 < len(parent.childrens) {
			iter.entry = parent.entries[pos]
			iter.node = parent.childrens[pos+1].minimum()
			iter.entryPos = 0
			return true
		}

		parent = parent.parent
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
