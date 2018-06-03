package avltree

import (
	"fmt"

	dot "github.com/asinglestep/godot"
)

// Key Key
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

// TreeNode avl节点
type TreeNode struct {
	height int32     // 树得高度
	key    Key       // key
	right  *TreeNode // 右子节点
	left   *TreeNode // 左子节点
	parent *TreeNode // 父节点
}

// NewTreeNode 新建一个节点
func NewTreeNode(key Key) *TreeNode {
	node := &TreeNode{
		key:    key,
		height: 1,
		right:  NewSentinel(),
		left:   NewSentinel(),
	}

	return node
}

// NewSentinel NewSentinel
func NewSentinel() *TreeNode {
	n := &TreeNode{
		key: -1,
	}

	return n
}

// rightRotate 右旋
func (node *TreeNode) rightRotate() *TreeNode {
	l := node.left
	lr := l.right

	node.left = lr
	l.right = node

	l.parent = node.parent
	node.parent = l

	node.height = node.max() + 1

	if !lr.isSentinel() {
		lr.parent = node
	}

	return l
}

// leftRotate 左旋
func (node *TreeNode) leftRotate() *TreeNode {
	r := node.right
	rl := r.left

	node.right = rl
	r.left = node

	r.parent = node.parent
	node.parent = r

	node.height = node.max() + 1

	if !rl.isSentinel() {
		rl.parent = node
	}

	return r
}

// leftRightRotate 先左旋再右旋
func (node *TreeNode) leftRightRotate() *TreeNode {
	node.left = node.left.leftRotate()
	return node.rightRotate()
}

// rightLeftRotate 先右旋再左旋
func (node *TreeNode) rightLeftRotate() *TreeNode {
	node.right = node.right.rightRotate()
	return node.leftRotate()
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

// findPrecursor 找到前驱节点
func (node *TreeNode) findPrecursor() *TreeNode {
	node = node.left

	for {
		if node.right.isSentinel() {
			return node
		}

		node = node.right
	}
}

// balanceFactor node的平衡因子
func (node *TreeNode) balanceFactor() int32 {
	var lh, rh int32

	if !node.left.isSentinel() {
		lh = node.left.high()
	}

	if !node.right.isSentinel() {
		rh = node.right.high()
	}

	return lh - rh
}

// Height 返回树的高度
func (node *TreeNode) high() int32 {
	if node.isSentinel() {
		return 0
	}

	return node.height
}

// max 取左右子树的最大高度
func (node *TreeNode) max() int32 {
	var lh, rh int32

	if !node.left.isSentinel() {
		lh = node.left.high()
	}

	if !node.right.isSentinel() {
		rh = node.right.high()
	}

	if lh >= rh {
		return lh
	}

	return rh
}

// isSentinel 是否是哨兵节点
func (node *TreeNode) isSentinel() bool {
	return node.key == -1
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

// reverse 倒序
func reverse(list []*TreeNode) []*TreeNode {
	listLen := len(list)

	for i := 0; i < listLen/2; i++ {
		list[i], list[listLen-i-1] = list[listLen-i-1], list[i]
	}

	return list
}

// dot dot
func (node *TreeNode) dot() (dNode *dot.Node, dEdge *dot.Edge) {
	// 添加node
	dNode = &dot.Node{}
	dNode.Name = fmt.Sprintf("%d", node.key)
	dNode.Attr = map[string]string{
		"label": fmt.Sprintf("\"<f0> | %d | <f1> \"", node.key),
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
