package avltree

import (
	"fmt"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

// TreeNode avl节点
type TreeNode struct {
	height int32        // 树得高度
	entry  *utils.Entry // 数据
	right  *TreeNode    // 右子节点
	left   *TreeNode    // 左子节点
	parent *TreeNode    // 父节点
}

var Sentinel = &TreeNode{
	entry: nil,
}

// NewTreeNode 新建一个节点
func NewTreeNode(entry *utils.Entry) *TreeNode {
	node := &TreeNode{
		entry:  entry,
		height: 1,
		right:  Sentinel,
		left:   Sentinel,
	}

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

// rightRotate 右旋
func (node *TreeNode) rightRotate() *TreeNode {
	l := node.left
	lr := l.right

	l.right = node
	l.parent = node.parent

	node.left = lr
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

	r.left = node
	r.parent = node.parent

	node.right = rl
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
	return node.entry == nil
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
	// 添加node
	dNode = &dot.Node{}
	dNode.Name = fmt.Sprintf("%d", node.GetKey())
	dNode.Attr = map[string]string{
		"label": fmt.Sprintf("\"<f0> | %d | <f1> \"", node.GetKey()),
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

// reverse 倒序
func reverse(list []*TreeNode) []*TreeNode {
	listLen := len(list)

	for i := 0; i < listLen/2; i++ {
		list[i], list[listLen-i-1] = list[listLen-i-1], list[i]
	}

	return list
}
