package avltree

import (
	"container/list"
	"fmt"
	"os/exec"
	"runtime"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

const (
	LEFT_2_HIGHER_THAN_RIGHT = 2
	RIGHT_2_HIGHER_THAN_LEFT = -2

	LEFT_1_HIGHER_THAN_RIGHT = 1
	RIGHT_1_HIGHER_THAN_LEFT = -1
)

// Tree Tree
type Tree struct {
	root       *TreeNode
	size       int // 节点数
	comparator utils.Comparator
}

// NewTree 创建一个avl树
func NewTree(comparator utils.Comparator) *Tree {
	t := &Tree{}
	t.root = NewSentinel()
	t.comparator = comparator

	return t
}

// Insert 插入一个节点
func (t *Tree) Insert(key, val interface{}) {
	// 插入新节点
	newNode := NewTreeNode(utils.NewEntry(key, val))
	t.insertNode(newNode)
	return
}

// insertNode 插入一个新节点
func (t *Tree) insertNode(node *TreeNode) {
	next := &t.root
	var parent *TreeNode

	for cur := *next; !cur.isSentinel(); cur = *next {
		res := t.comparator.Compare(node.GetKey(), cur.GetKey())
		if res == utils.Et {
			return
		}

		parent = cur
		if res == utils.Lt {
			next = &cur.left
		} else {
			next = &cur.right
		}
	}

	*next = node
	node.parent = parent

	// 插入修复
	t.insertFixUp(node)
	t.size++
}

// insertFixUp 插入修复
func (t *Tree) insertFixUp(node *TreeNode) {
	var bRotate bool
	fixNode := node.parent

	for fixNode != nil {
		bRotate = false

		// 插入节点后，树的高度没有变化，不需要修复
		if fixNode.height == fixNode.max()+1 {
			return
		}

		parent := fixNode.parent
		isLeft := fixNode.isLeft()

		switch fixNode.balanceFactor() {
		case LEFT_2_HIGHER_THAN_RIGHT:
			// 左子树比右子树高2
			bRotate = true
			// 插入节点 在 修复节点的左子树的左子树上，只需要旋转一次
			// 插入节点 在 修复节点的左子树的右子树上，需要旋转两次
			bSingleRorate := t.comparator.Compare(node.GetKey(), fixNode.left.GetKey()) == utils.Lt
			fixNode = t.caseLeft2HigherThanRight(fixNode, bSingleRorate)

		case RIGHT_2_HIGHER_THAN_LEFT:
			// 右子树比左子树高2
			bRotate = true
			// 插入节点 在 修复节点的右子树的右子树上，只需要旋转一次
			// 插入节点 在 修复节点的右子树的左子树上，需要旋转两次
			bSingleRorate := t.comparator.Compare(node.GetKey(), fixNode.right.GetKey()) == utils.Gt
			fixNode = t.caseRight2HigherThanLeft(fixNode, bSingleRorate)
		}

		if bRotate {
			t.updateChildren(parent, fixNode, isLeft)
		}

		fixNode.height = fixNode.max() + 1
		fixNode = fixNode.parent
	}
}

// Delete 删除指定的节点
func (t *Tree) Delete(key interface{}) {
	node := t.root

	for !node.isSentinel() {
		res := t.comparator.Compare(key, node.GetKey())
		if res == utils.Et {
			t.deleteNode(node)
			return
		}

		if res == utils.Lt {
			node = node.left
		} else {
			node = node.right
		}
	}
}

// deleteNode 删除节点
func (t *Tree) deleteNode(node *TreeNode) {
	hasLeft := !node.left.isSentinel()
	hasRight := !node.right.isSentinel()

	if hasLeft && hasRight {
		// 删除节点有左右子节点
		last := node

		switch node.balanceFactor() {
		case LEFT_1_HIGHER_THAN_RIGHT:
			// 前驱节点只可能有右节点
			node = node.findPrecursor()
			hasRight = true
			// hasLeft = false

		default:
			// 后继节点只可能有左节点
			node = node.findSuccessor()
			hasLeft = true
			// hasRight = false
		}

		last.entry = node.entry
	}

	parent := node.parent
	children := node.right
	if hasLeft {
		children = node.left
	}

	// 改变删除节点的孩子节点的父节点
	children.parent = parent
	t.updateChildren(parent, children, node.isLeft())

	// 删除节点
	node.free()
	t.deleteFixUp(parent)
	t.size--
}

// deleteFixUp 删除修复
func (t *Tree) deleteFixUp(node *TreeNode) {
	for node != nil {
		bHeightChange := false
		parent := node.parent
		isLeft := node.isLeft()

		// 删除节点后，虽然树的高度没有变化
		// 但是可能会导致左右子树的高度差为2
		if node.height == node.max()+1 {
			bHeightChange = true
		}

		switch node.balanceFactor() {
		case LEFT_2_HIGHER_THAN_RIGHT:
			// 删除节点之后，左子树比右子树高2
			// 修复节点的左子树的左子树比右子树高 或者 修复节点的左子树的左子树和右子树一样高，只需要旋转一次
			// 修复节点的左子树的右子树比左子树高，需要旋转两次
			bSingleRorate := node.left.right.high() <= node.left.left.high()
			node = t.caseLeft2HigherThanRight(node, bSingleRorate)

		case RIGHT_2_HIGHER_THAN_LEFT:
			// 删除节点之后，右子树比左子树高2
			// 修复节点的右子树的右子树比左子树高 或者 修复节点的右子树的右子树和左子树一样高，只需要旋转一次
			// 修复节点的右子树的左子树比右子树高，需要旋转两次
			bSingleRorate := node.right.left.high() <= node.right.right.high()
			node = t.caseRight2HigherThanLeft(node, bSingleRorate)

		default:
			if bHeightChange {
				// 删除节点之后，左右子树的高度差不为2，且高度没有发生变化，不需要修复
				return
			}
		}

		t.updateChildren(parent, node, isLeft)
		node.height = node.max() + 1
		node = node.parent
	}
}

// Search 查找key指定的节点
func (t *Tree) Search(key interface{}) *TreeNode {
	node := t.root

	for !node.isSentinel() {
		res := t.comparator.Compare(node.GetKey(), key)
		if res == utils.Et {
			return node
		}

		if res == utils.Lt {
			node = node.right
		} else {
			node = node.left
		}
	}

	return nil
}

// SearchRange 查找key在[min, max]之间的节点
func (t *Tree) SearchRange(min, max interface{}) []*TreeNode {
	stack := list.New()
	list := []*TreeNode{}
	node := t.root

	for !node.isSentinel() {
		// 将key大于等于min的节点加入到stack中
		if t.comparator.Compare(node.GetKey(), min) != utils.Lt {
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
		if t.comparator.Compare(node.GetKey(), max) == utils.Gt {
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
func (t *Tree) SearchRangeLowerBoundKeyWithLimit(key interface{}, limit int64) []*TreeNode {
	var count int64
	stack := list.New()
	list := make([]*TreeNode, 0, limit)
	node := t.root

	for !node.isSentinel() {
		// 将key大于等于key的节点加入到stack中
		if t.comparator.Compare(node.GetKey(), key) != utils.Lt {
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
func (t *Tree) SearchRangeUpperBoundKeyWithLimit(key interface{}, limit int64) []*TreeNode {
	var count int64
	stack := list.New()
	list := make([]*TreeNode, 0, limit)
	node := t.root

	for !node.isSentinel() {
		// 将key小于等于key的节点加入到stack中
		if t.comparator.Compare(node.GetKey(), key) != utils.Gt {
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
			fmt.Printf("节点为: %v, \t节点的高度为: %v, \t父节点为: %v\n", node.GetKey(), node.height, node.parent.GetKey())
		} else {
			fmt.Printf("节点为: %v, \t节点的高度为: %v, \t此节点为根节点\n", node.GetKey(), node.height)
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
	keys := make([]interface{}, 0)

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

		keys = append(keys, node.GetKey())

		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}

	// 验证顺序
	for i := 0; i < len(keys)-1; i++ {
		if t.comparator.Compare(keys[i], keys[i+1]) == utils.Gt {
			fmt.Printf("Key顺序错误\n")
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

// updateChild 更新父节点的子节点
func (t *Tree) updateChildren(parent, children *TreeNode, isLeft bool) {
	if parent == nil {
		t.root = children
	} else if isLeft {
		parent.left = children
	} else {
		parent.right = children
	}
}

// caseLeft2HigherThanRight 左子树比右子树高2
func (t *Tree) caseLeft2HigherThanRight(node *TreeNode, bSingleRorate bool) *TreeNode {
	if bSingleRorate {
		return node.rightRotate()
	}

	return node.leftRightRotate()
}

// caseRight2HigherThanLeft 右子树比左子树高2
func (t *Tree) caseRight2HigherThanLeft(node *TreeNode, bSingleRorate bool) *TreeNode {
	if bSingleRorate {
		return node.leftRotate()
	}

	return node.rightLeftRotate()
}
