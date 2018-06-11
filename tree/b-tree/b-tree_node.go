package btree

import (
	"fmt"
	"strings"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

// Key Key
type Key int

// less less
// k < key, 返回true
func (k Key) less(key Key) bool {
	if k < key {
		return true
	}

	return false
}

// more more
// k > key, 返回true
func (k Key) more(key Key) bool {
	if k > key {
		return true
	}

	return false
}

// Value Value
type Value int

// Entry Entry
type Entry struct {
	Key   Key
	Value Value
}

// NewEntry NewEntry
func NewEntry(key Key, value Value) *Entry {
	e := &Entry{}
	e.Key = key
	e.Value = value

	return e
}

// String String
func (e *Entry) String() string {
	return fmt.Sprintf("(%v, %v)  ", e.Key, e.Value)
}

// TreeNode TreeNode
type TreeNode struct {
	parent    *TreeNode      // 父节点
	childrens []*TreeNode    // 子节点
	entries   []*utils.Entry // 节点数据

	iOffset int  // 偏移量
	bHandle bool // true: 第一次处理这个节点
}

// NewLeafNode 叶子节点
func NewLeafNode() *TreeNode {
	n := &TreeNode{}
	n.childrens = nil

	return n
}

// findKeyPosition 在节点中查找第一个大于等于key的位置，没有比key大的节点，则返回此节点最后一个key
func (node *TreeNode) findKeyPosition(comparator utils.Comparator, key interface{}) int {
	i, j := 0, len(node.entries)

	for i < j {
		h := int(uint(i+j) >> 1)
		if comparator.Compare(node.entries[h].GetKey(), key) == utils.Lt {
			i = h + 1
		} else {
			j = h
		}
	}

	return i
}

// findPrecursor 找到pos位置的前驱节点
func (node *TreeNode) findPrecursor(pos int) *TreeNode {
	node = node.childrens[pos]
	for !node.isLeaf() {
		node = node.childrens[len(node.childrens)-1]
	}

	return node
}

// findSuccessor 找到pos位置的后继节点
func (node *TreeNode) findSuccessor(pos int) *TreeNode {
	node = node.childrens[pos+1]
	for !node.isLeaf() {
		node = node.childrens[0]
	}

	return node
}

// isLeaf 是否是叶子节点
func (node *TreeNode) isLeaf() bool {
	return node.childrens == nil
}

// isFull 是否是满节点
func (node *TreeNode) isFull(maxEntry int) bool {
	return len(node.entries) == maxEntry
}

// PrintBTreeNode PrintBTreeNode
func (node *TreeNode) PrintBTreeNode() {
	if node.parent != nil {
		offset := node.parent.iOffset

		if offset == len(node.parent.entries) {
			fmt.Printf("节点key: %v, \t此节点为父节点key[%v]的右节点\n", node.entries, node.parent.entries[offset-1])
		} else {
			fmt.Printf("节点key: %v, \t此节点为父节点key[%v]的左节点\n", node.entries, node.parent.entries[offset])
		}
	} else {
		if len(node.entries) == 0 {
			fmt.Printf("此b树为一个空树\n")
		} else {
			fmt.Printf("节点key: %v, \t此节点为根节点\n", node.entries)
		}
	}
}

// dot dot
func (node *TreeNode) dot(nName string, pName string) (dNode *dot.Node, dEdge *dot.Edge) {
	// 添加一个node
	attrValues := make([]string, 0, len(node.entries))

	for i, entry := range node.entries {
		attrValues = append(attrValues, fmt.Sprintf("<f%d> | %d ", i, entry.Key))
	}

	attr := "\"" + strings.Join(attrValues, "|") + fmt.Sprintf("| <f%d>", len(node.entries)) + "\""

	dNode = &dot.Node{}
	dNode.Name = nName
	dNode.Attr = map[string]string{
		"label": attr,
	}

	// 添加一个edge
	if node.parent != nil {
		pos := node.parent.findKeyPosition(node.entries[0].Key)
		dEdge = &dot.Edge{}
		dEdge.Src = pName
		dEdge.SrcPort = ":f" + fmt.Sprintf("%d", pos)
		dEdge.Dst = nName
	}

	return dNode, dEdge
}
