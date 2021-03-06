package rbtree

import (
	"fmt"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

const (
	RED   Color = 1
	BLACK Color = 2
)

// Color Color
type Color uint32

func (c Color) String() string {
	switch c {
	case 1:
		return "Red"
	case 2:
		return "Black"
	}

	return "Unknown"
}

// TreeNode 红黑树节点
type TreeNode struct {
	color  Color
	entry  *utils.Entry
	left   *TreeNode
	right  *TreeNode
	parent *TreeNode
}

// Sentinel Sentinel
var Sentinel = &TreeNode{
	color: BLACK,
	entry: nil,
}

// NewTreeNode 新建一个节点
func NewTreeNode(entry *utils.Entry) *TreeNode {
	node := &TreeNode{}
	node.color = RED
	node.entry = entry
	node.left = Sentinel
	node.right = Sentinel

	return node
}

// GetKey 获取key
func (node *TreeNode) GetKey() interface{} {
	return node.entry.GetKey()
}

// GetValue 获取value
func (node *TreeNode) GetValue() interface{} {
	return node.entry.GetValue()
}

// leftRotate 左旋
func (node *TreeNode) leftRotate() *TreeNode {
	r := node.right
	rl := r.left

	node.right = rl
	r.left = node

	r.parent = node.parent
	node.parent = r

	if !rl.isSentinel() {
		rl.parent = node
	}

	return r
}

// rightRotate 右旋
func (node *TreeNode) rightRotate() *TreeNode {
	l := node.left
	lr := l.right

	node.left = lr
	l.right = node

	l.parent = node.parent
	node.parent = l

	if !lr.isSentinel() {
		lr.parent = node
	}

	return l
}

// isBlack 是否是黑色节点
func (node *TreeNode) isBlack() bool {
	if node.isSentinel() {
		return true
	}

	return node.color == BLACK
}

// isRed 是否是红色节点
func (node *TreeNode) isRed() bool {
	if node.isSentinel() {
		return false
	}

	return node.color == RED
}

// isLeft 是否是左节点
func (node *TreeNode) isLeft() bool {
	// node为根节点
	if node.parent == nil {
		return false
	}

	return node == node.parent.left
}

// isRight 是否是右节点
func (node *TreeNode) isRight() bool {
	// node为根节点
	if node.parent == nil {
		return false
	}

	return node == node.parent.right
}

// findSuccessor 找到node的后继节点
func (node *TreeNode) findSuccessor() *TreeNode {
	node = node.right

	for {
		if node.left.isSentinel() {
			return node
		}

		node = node.left
	}
}

// findPrecursor 找到node的前驱节点
func (node *TreeNode) findPrecursor() *TreeNode {
	node = node.left

	for {
		if node.right.isSentinel() {
			return node
		}

		node = node.right
	}
}

// isSentinel 是否是哨兵节点
func (node *TreeNode) isSentinel() bool {
	return node == Sentinel
}

// getBrother 获取兄弟节点
func (node *TreeNode) getBrother() *TreeNode {
	if node.isRight() {
		return node.parent.left
	}

	return node.parent.right
}

// free free
func (node *TreeNode) free() {
	node.parent = nil
	node.left = nil
	node.right = nil
	node.entry = nil
}

// minimum 以当前节点为根节点，中序遍历后，树的最小节点
func (node *TreeNode) minimum() *TreeNode {
	if node.isSentinel() {
		return nil
	}

	for !node.left.isSentinel() {
		node = node.left
	}

	return node
}

// maximum 以当前节点为根节点，中序遍历后，树的最大节点
func (node *TreeNode) maximum() *TreeNode {
	if node.isSentinel() {
		return nil
	}

	for !node.right.isSentinel() {
		node = node.right
	}

	return node
}

// next 中序遍历node的下一个节点
func (node *TreeNode) next() *TreeNode {
	if node.isSentinel() {
		return nil
	}

	// 在右子树中找最小的节点
	if n := node.right.minimum(); n != nil {
		return n
	}

	parent := node.parent
	for parent != nil && node.isRight() {
		node = parent
		parent = node.parent
	}

	return parent
}

// prev 中序遍历node的上一个节点
func (node *TreeNode) prev() *TreeNode {
	if node.isSentinel() {
		return nil
	}

	// 在左子树中找最大的节点
	if n := node.left.maximum(); n != nil {
		return n
	}

	parent := node.parent
	for parent != nil && node.isLeft() {
		node = parent
		parent = node.parent
	}

	return parent
}

// dot dot
func (node *TreeNode) dot() (dNode *dot.Node, dEdge *dot.Edge) {
	color := "#FF0000"
	if node.isBlack() {
		color = "#0F0F0F"
	}

	// 添加node
	dNode = &dot.Node{}
	dNode.Name = fmt.Sprintf("%d", node.GetKey())

	dNode.Attr = map[string]string{
		"label":     fmt.Sprintf("\"<f0> | %d | <f1> \"", node.GetKey()),
		"fillcolor": "\"" + color + "\"",
		"style":     "filled",
		"fontcolor": "\"#FFFFFF\"",
		"color":     "\"#FFFF00\"",
	}

	// 添加edge
	if node.parent != nil {
		dEdge = &dot.Edge{}
		dEdge.Src = fmt.Sprintf("%d", node.parent.GetKey())

		if node.isLeft() {
			dEdge.SrcPort = ":f0"
		} else {
			dEdge.SrcPort = ":f1"
		}

		dEdge.Dst = fmt.Sprintf("%d", node.GetKey())
	}

	return dNode, dEdge
}

// nodeColorStat 节点颜色统计
type nodeColorStat struct {
	blackCount      int64
	leftBlackCount  int64
	rightBlackCount int64
}

// reverse 倒序
func reverse(list []*TreeNode) []*TreeNode {
	listLen := len(list)

	for i := 0; i < listLen/2; i++ {
		list[i], list[listLen-i-1] = list[listLen-i-1], list[i]
	}

	return list
}
