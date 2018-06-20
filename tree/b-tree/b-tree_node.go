package btree

import (
	"fmt"
	"strings"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

// TreeNode TreeNode
type TreeNode struct {
	parent    *TreeNode      // 父节点
	childrens []*TreeNode    // 子节点
	entries   []*utils.Entry // 节点数据
}

// NewNode NewNode
func NewNode() *TreeNode {
	n := &TreeNode{}

	return n
}

// insertEntry 将新的entry插入到node的pos位置上
func (node *TreeNode) insertEntry(entry *utils.Entry, pos int) {
	newEntries := make([]*utils.Entry, len(node.entries)+1)
	newEntries[pos] = entry
	copy(newEntries[:pos], node.entries[:pos])
	copy(newEntries[pos+1:], node.entries[pos:])
	node.entries = newEntries
}

// insertChildren 将children插入到node的pos位置上
func (node *TreeNode) insertChildren(children *TreeNode, pos int) {
	newCs := make([]*TreeNode, len(node.childrens)+1)
	newCs[pos] = children
	copy(newCs[:pos], node.childrens[:pos])
	copy(newCs[pos+1:], node.childrens[pos:])
	node.childrens = newCs
}

// updateChildrensParent 更新的node的childrens的父节点
func (node *TreeNode) updateChildrensParent(parent *TreeNode) {
	for i := range node.childrens {
		node.childrens[i].parent = parent
	}
}

// findLowerBoundKeyPosition 在节点中查找第一个大于等于key的位置，没有比key大的，则返回node.entries的长度
func (node *TreeNode) findLowerBoundKeyPosition(comparator utils.Comparator, key interface{}) (pos int, bFound bool) {
	if len(node.entries) == 0 {
		return 0, false
	}

	i, j := 0, len(node.entries)

	for i < j {
		h := int(uint(i+j) >> 1)
		if comparator.Compare(node.entries[h].GetKey(), key) == utils.Lt {
			i = h + 1
		} else {
			j = h
		}
	}

	if i == len(node.entries) {
		return i, false
	}

	if comparator.Compare(node.entries[i].GetKey(), key) == utils.Et {
		return i, true
	}

	return i, false
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

// free free
func (node *TreeNode) free() {
	node.parent = nil
	node.childrens = nil
	node.entries = nil
}

// minimum 以当前节点为根节点，中序遍历后，树的最小节点
func (node *TreeNode) minimum() *TreeNode {
	for !node.isLeaf() {
		node = node.childrens[0]
	}

	return node
}

// maximum 以当前节点为根节点，中序遍历后，树的最大节点
func (node *TreeNode) maximum() *TreeNode {
	for !node.isLeaf() {
		node = node.childrens[len(node.childrens)-1]
	}

	return node
}

// printBTreeNode printBTreeNode
func (node *TreeNode) printBTreeNode(offset *int) (str string) {
	if node.parent != nil {
		if *offset == len(node.parent.entries) {
			str = fmt.Sprintf("节点key: %v, \t此节点为父节点key[%v]的右节点\n", node.printBTreeNodeKeys(), node.parent.entries[*offset-1].GetKey())
			*offset = 0
		} else {
			str = fmt.Sprintf("节点key: %v, \t此节点为父节点key[%v]的左节点\n", node.printBTreeNodeKeys(), node.parent.entries[*offset].GetKey())
			(*offset)++
		}
	} else {
		if len(node.entries) == 0 {
			str = fmt.Sprintf("此b树为一个空树\n")
		} else {
			str = fmt.Sprintf("节点key: %v, \t此节点为根节点\n", node.printBTreeNodeKeys())
		}
	}

	return str
}

// printBTreeNodeKeys printBTreeNodeKeys
func (node *TreeNode) printBTreeNodeKeys() string {
	keys := make([]string, 0, len(node.entries))

	for _, v := range node.entries {
		keys = append(keys, fmt.Sprintf("%v", v.GetKey()))
	}

	return strings.Join(keys, ",")
}

// dot dot
func (node *TreeNode) dot(comparator utils.Comparator, nName string, pName string) (dNode *dot.Node, dEdge *dot.Edge) {
	// 添加一个node
	attrValues := make([]string, 0, len(node.entries))

	for i, entry := range node.entries {
		attrValues = append(attrValues, fmt.Sprintf("<f%d> | %d ", i, entry.GetKey()))
	}

	attr := "\"" + strings.Join(attrValues, "|") + fmt.Sprintf("| <f%d>", len(node.entries)) + "\""

	dNode = &dot.Node{}
	dNode.Name = nName
	dNode.Attr = map[string]string{
		"label": attr,
	}

	// 添加一个edge
	if node.parent != nil {
		pos, _ := node.parent.findLowerBoundKeyPosition(comparator, node.entries[0].GetKey())
		dEdge = &dot.Edge{}
		dEdge.Src = pName
		dEdge.SrcPort = ":f" + fmt.Sprintf("%d", pos)
		dEdge.Dst = nName
	}

	return dNode, dEdge
}
