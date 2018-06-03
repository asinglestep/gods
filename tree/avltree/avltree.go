package avltree

import (
	"container/list"
	"fmt"
	"os/exec"
	"runtime"

	dot "github.com/asinglestep/godot"
)

const (
	LEFT_2_HIGHER_THAN_RIGHT = 2
	RIGHT_2_HIGHER_THAN_LEFT = -2
)

// Tree Tree
type Tree struct {
	root *TreeNode
}

// NewTree 创建一个avl树
func NewTree() *Tree {
	t := &Tree{}
	t.root = NewSentinel()

	return t
}

// Insert 插入一个节点
func (t *Tree) Insert(key Key) {
	// 插入新节点
	newNode := NewTreeNode(key)
	fixNode := t.insertNode(newNode)

	// 插入修复
	t.insertFixUp(fixNode)
	return
}

// insertNode 插入一个新节点
func (t *Tree) insertNode(newNode *TreeNode) (fixNode *TreeNode) {
	if t.root.isSentinel() {
		t.root = newNode
		return newNode
	}

	tmpNode := t.root

	for {
		if newNode.key.less(tmpNode.key) {
			if tmpNode.left.isSentinel() {
				tmpNode.left = newNode
				break
			}

			tmpNode = tmpNode.left
		} else {
			if tmpNode.right.isSentinel() {
				tmpNode.right = newNode
				break
			}

			tmpNode = tmpNode.right
		}
	}

	newNode.parent = tmpNode
	return newNode
}

// insertFixUp 插入修复
func (t *Tree) insertFixUp(fixNode *TreeNode) {
	key := fixNode.key
	tmpNode := fixNode.parent

	for tmpNode != nil {
		// 插入节点后，树的高度没有变化，不需要修复
		if tmpNode.height == tmpNode.max()+1 {
			return
		}

		tmpNodeParent := tmpNode.parent
		isLeft := tmpNode.isLeft()

		switch tmpNode.balanceFactor() {
		// 左子树比右子树高2
		case LEFT_2_HIGHER_THAN_RIGHT:
			// 插入节点小于左节点
			if key.less(tmpNode.left.key) {
				tmpNode = tmpNode.rightRotate()
			} else {
				tmpNode = tmpNode.leftRightRotate()
			}

		// 右子树比左子树高2
		case RIGHT_2_HIGHER_THAN_LEFT:
			// 插入节点小于右节点
			if key.less(tmpNode.right.key) {
				tmpNode = tmpNode.rightLeftRotate()
			} else {
				tmpNode = tmpNode.leftRotate()
			}
		}

		if tmpNodeParent == nil {
			// 修复根节点
			t.root = tmpNode
		} else if isLeft {
			tmpNodeParent.left = tmpNode
		} else {
			tmpNodeParent.right = tmpNode
		}

		tmpNode.height = tmpNode.max() + 1
		tmpNode = tmpNode.parent
	}
}

// Delete 删除key指定的节点
func (t *Tree) Delete(key Key) {
	fixNode := t.deleteNode(key)
	// 没找到不处理
	if fixNode == nil {
		return
	}

	t.deleteFixUp(fixNode)
	return
}

// deleteNode 删除节点
func (t *Tree) deleteNode(key Key) (fixNode *TreeNode) {
	delNode := t.root

	for {
		// 没找到
		if delNode.isSentinel() {
			return nil
		}

		if key.equal(delNode.key) {
			break
		}

		if key.less(delNode.key) {
			delNode = delNode.left
		} else {
			delNode = delNode.right
		}
	}

	hasLeft := !delNode.left.isSentinel()
	hasRight := !delNode.right.isSentinel()

	// 删除节点有左右子节点
	if hasLeft && hasRight {
		findNode := delNode

		if delNode.balanceFactor() == 1 {
			// 前驱节点只可能有右节点
			delNode = delNode.findPrecursor()
			hasRight = true
		} else {
			// 后继节点只可能有左节点
			delNode = delNode.findSuccessor()
			hasLeft = true
		}

		findNode.key = delNode.key
	}

	parent := delNode.parent
	children := delNode.right
	if hasLeft {
		children = delNode.left
	}

	if delNode.isLeft() {
		// 删除节点的孩子节点变成父节点的左节点
		parent.left = children
	} else if delNode.isRight() {
		// 删除节点的孩子节点变成父节点的右节点
		parent.right = children
	} else {
		// 删除根节点
		t.root = children
	}

	// 改变删除节点的孩子节点的父节点
	children.parent = parent

	// 删除节点
	delNode.parent = nil
	delNode.left = nil
	delNode.right = nil
	return parent
}

// deleteFixUp 删除修复
func (t *Tree) deleteFixUp(fixNode *TreeNode) {
	for fixNode != nil {
		isHeightChange := false
		fixNodeParent := fixNode.parent
		isLeft := fixNode.isLeft()

		// 删除节点后，树的高度没有变化，不需要修复
		if fixNode.height == fixNode.max()+1 {
			isHeightChange = true
		}

		switch fixNode.balanceFactor() {
		// 删除节点之后，左子树比右子树高2
		case LEFT_2_HIGHER_THAN_RIGHT:
			if fixNode.left.right.high() > fixNode.left.left.high() {
				// 左节点的右子树比左子树高
				fixNode = fixNode.leftRightRotate()
			} else {
				fixNode = fixNode.rightRotate()
			}

		// 删除节点之后，右子树比左子树高2
		case RIGHT_2_HIGHER_THAN_LEFT:
			if fixNode.right.left.high() > fixNode.right.right.high() {
				// 右节点的左子树比右子树高
				fixNode = fixNode.rightLeftRotate()
			} else {
				fixNode = fixNode.leftRotate()
			}

		// 删除节点之后，左右子树的高度差不为2，且高度没有发生变化，不需要修复
		default:
			if isHeightChange {
				return
			}
		}

		if fixNodeParent == nil {
			// 根节点
			t.root = fixNode
		} else if isLeft {
			fixNodeParent.left = fixNode
		} else {
			fixNodeParent.right = fixNode
		}

		fixNode.height = fixNode.max() + 1
		fixNode = fixNode.parent
	}
}

// Search 查找key指定的节点
func (t *Tree) Search(key Key) *TreeNode {
	node := t.root

	for {
		if node.isSentinel() {
			return nil
		}

		if key.equal(node.key) {
			break
		}

		if key.less(node.key) {
			node = node.left
		} else {
			node = node.right
		}
	}

	return node
}

// SearchRange 查找key在[min, max]之间的节点
func (t *Tree) SearchRange(min Key, max Key) []*TreeNode {
	stack := list.New()
	list := []*TreeNode{}
	node := t.root

	for {
		if node.isSentinel() {
			break
		}

		// 将key大于等于min的节点加入到stack中
		if !node.key.less(min) {
			stack.PushBack(node)
			node = node.left
		} else {
			node = node.right
		}
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		// node.key大于max，退出
		if node.key.more(max) {
			break
		}

		list = append(list, node)
		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}

	return list
}

// SearchRangeLowerBoundKeyWithLimit 查找大于等于key的limit个节点
func (t *Tree) SearchRangeLowerBoundKeyWithLimit(key Key, limit int64) []*TreeNode {
	var count int64
	stack := list.New()
	list := make([]*TreeNode, 0, limit)
	node := t.root

	for {
		if node.isSentinel() {
			break
		}

		// 将key大于等于key的节点加入到stack中
		if !node.key.less(key) {
			stack.PushBack(node)
			node = node.left
		} else {
			node = node.right
		}
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		list = append(list, node)
		count++
		if count == limit {
			break
		}

		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}

	return list
}

// SearchRangeUpperBoundKeyWithLimit 找到小于等于key的limit个节点
func (t *Tree) SearchRangeUpperBoundKeyWithLimit(key Key, limit int64) []*TreeNode {
	var count int64
	stack := list.New()
	list := make([]*TreeNode, 0, limit)
	node := t.root

	for {
		if node.isSentinel() {
			break
		}

		// 将key小于等于key的节点加入到stack中
		if !node.key.more(key) {
			stack.PushBack(node)
			node = node.right
		} else {
			node = node.left
		}
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		list = append(list, node)
		count++
		if count == limit {
			break
		}

		node = node.left
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.right
		}
	}

	return reverse(list)
}

// PrintAvlTree PrintAvlTree
func (t *Tree) PrintAvlTree() {
	node := t.root
	stack := list.New()

	for !node.isSentinel() {
		stack.PushBack(node)
		node = node.left
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		if node.parent != nil {
			fmt.Printf("节点为: %v, \t节点的高度为: %v, \t父节点为: %v\n", node.key, node.height, node.parent.key)
		} else {
			fmt.Printf("节点为: %v, \t节点的高度为: %v, \t此节点为根节点\n", node.key, node.height)
		}

		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}
}

// VerifAvlTree 验证是否是一个avl树
func (t *Tree) VerifAvlTree() bool {
	node := t.root
	stack := list.New()
	keys := make([]Key, 0)

	for !node.isSentinel() {
		stack.PushBack(node)
		node = node.left
	}

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		// 如果node的平衡因子大于2或者小于-2
		bf := node.balanceFactor()
		if bf >= 2 || bf <= -2 {
			fmt.Printf("bf >= 2 || bf <= -2\n")
			return false
		}

		keys = append(keys, node.key)

		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}

	// 验证顺序
	for i := 0; i < len(keys)-1; i++ {
		if keys[i] > keys[i+1] {
			fmt.Printf("顺序错误\n")
			return false
		}
	}

	return true
}

// Dot Dot
func (t *Tree) Dot() error {
	node := t.root
	if node.isSentinel() {
		return nil
	}

	stack := list.New()

	dGraph := dot.NewGraph()
	dGraph.SetNodeGlobalAttr(map[string]string{
		"height": ".1",
		"shape":  "record",
		"width":  ".1",
	})

	stack.PushBack(node)

	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		dNode, dEdge := node.dot()
		dGraph.AddNode(dNode)
		if dEdge != nil {
			dGraph.AddEdge(dEdge)
		}

		// 将左右子节点压入stack
		if !node.left.isSentinel() {
			stack.PushBack(node.left)
		}

		if !node.right.isSentinel() {
			stack.PushBack(node.right)
		}
	}

	if err := dot.GenerateDotFile("avltree.dot", dGraph); err != nil {
		return err
	}

	if err := exec.Command("dot", "-Tpng", "avltree.dot", "-o", "avltree.png").Run(); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "avltree.png")
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
