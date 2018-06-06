package treap

import (
	"fmt"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

// TreeNode treap节点
type TreeNode struct {
	left     *TreeNode
	right    *TreeNode
	parent   *TreeNode
	priority uint32       // 优先级
	entry    *utils.Entry // 数据
}

var SentinelNode = &TreeNode{
	entry: nil,
}

// NewTreeNode NewTreeNode
func NewTreeNode(entry *utils.Entry, priority uint32) *TreeNode {
	node := &TreeNode{}
	node.entry = entry
	node.priority = priority
	node.left = SentinelNode
	node.right = SentinelNode

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

	r.left = node
	r.parent = node.parent

	node.right = rl
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

	l.right = node
	l.parent = node.parent

	node.left = lr
	node.parent = l

	if !lr.isSentinel() {
		lr.parent = node
	}

	return l
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

// isSentinel 是否是哨兵节点
func (node *TreeNode) isSentinel() bool {
	return node.entry == nil
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

// free free
func (node *TreeNode) free() {
	node.parent = nil
	node.left = nil
	node.right = nil
	node.entry = nil
}

// dot dot
func (node *TreeNode) dot() (dNode *dot.Node, dEdge *dot.Edge) {
	// 添加node
	dNode = &dot.Node{}
	dNode.Name = fmt.Sprintf("%d", node.entry.GetKey())
	dNode.Attr = map[string]string{
		"label": fmt.Sprintf("\"<f0> | k: %d、 p: %d | <f1> \"", node.entry.GetKey(), node.priority),
	}

	// 添加edge
	if node.parent != nil {
		dEdge = &dot.Edge{}
		dEdge.Src = fmt.Sprintf("%d", node.parent.entry.GetKey())

		if node.isLeft() {
			dEdge.SrcPort = ":f0"
		} else {
			dEdge.SrcPort = ":f1"
		}

		dEdge.Dst = fmt.Sprintf("%d", node.entry.GetKey())
	}

	return dNode, dEdge
}
