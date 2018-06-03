package treap

import (
	"fmt"

	dot "github.com/asinglestep/godot"
)

type Key int32

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

// equal equal
// k = key, 返回true
func (k Key) equal(key Key) bool {
	if k == key {
		return true
	}

	return false
}

// TreeNode treap节点
type TreeNode struct {
	key      Key
	priority uint32 // 优先级
	left     *TreeNode
	right    *TreeNode
	parent   *TreeNode
}

var SentinelNode = &TreeNode{
	key: -1,
}

// NewTreeNode NewTreeNode
func NewTreeNode(key Key, priority uint32) *TreeNode {
	node := &TreeNode{}
	node.key = key
	node.priority = priority
	node.left = SentinelNode
	node.right = SentinelNode

	return node
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
	return node.key == -1
}

// dot dot
func (node *TreeNode) dot() (dNode *dot.Node, dEdge *dot.Edge) {
	// 添加node
	dNode = &dot.Node{}
	dNode.Name = fmt.Sprintf("%d", node.key)
	dNode.Attr = map[string]string{
		"label": fmt.Sprintf("\"<f0> | k: %d、 p: %d | <f1> \"", node.key, node.priority),
	}

	// 添加edge
	if node.parent != nil {
		dEdge = &dot.Edge{}
		dEdge.Src = fmt.Sprintf("%d", node.parent.key)

		if node.isLeft() {
			dEdge.SrcPort = ":f0"
		} else {
			dEdge.SrcPort = ":f1"
		}

		dEdge.Dst = fmt.Sprintf("%d", node.key)
	}

	return dNode, dEdge
}
