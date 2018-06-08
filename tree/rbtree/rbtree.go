package rbtree

import (
	"container/list"
	"fmt"
	"os/exec"
	"runtime"

	dot "github.com/asinglestep/godot"
	"github.com/asinglestep/gods/utils"
)

// Tree Tree
type Tree struct {
	root       *TreeNode
	size       int // 节点数
	comparator utils.Comparator
}

// NewTree NewTree
func NewTree(comparator utils.Comparator) *Tree {
	t := &Tree{}
	t.root = Sentinel
	t.comparator = comparator

	return t
}

// Insert 插入
func (t *Tree) Insert(key, val interface{}) {
	node := NewTreeNode(utils.NewEntry(key, val))
	t.insertNode(node)
	return
}

// insertNode 插入新节点
func (t *Tree) insertNode(node *TreeNode) {
	next := &t.root
	var parent *TreeNode

	for cur := *next; !cur.isSentinel(); cur = *next {
		res := t.comparator.Compare(cur.GetKey(), node.GetKey())
		if res == utils.Et {
			cur.entry.SetValue(node.GetValue())
			return
		}

		parent = cur
		if res == utils.Lt {
			// 在右子树查找
			next = &cur.right
		} else {
			// 在左子树查找
			next = &cur.left
		}
	}

	*next = node
	node.parent = parent
	t.insertFixUp(node)
	t.size++
}

// insertFixUp 插入节点后进行修复
func (t *Tree) insertFixUp(node *TreeNode) {
	for {
		parent := node.parent
		if parent == nil {
			// 修复节点为根节点，变成黑色节点
			t.caseRoot()
			return
		}

		if parent.isBlack() {
			// 修复节点的父节点为黑色，不处理
			return
		}

		// 获取叔叔节点
		uncle := parent.getBrother()
		grandfather := parent.parent
		nIsLeft := node.isLeft()
		nIsRight := node.isRight()
		pIsLeft := parent.isLeft()
		pIsRight := parent.isRight()

		switch {
		case uncle.isRed():
			// 叔叔节点也是红色
			node = t.iCaseUncleIsRed(parent, uncle)

		case nIsLeft && pIsLeft:
			// 当前节点为左节点，父节点为左节点
			t.iCaseNodeAndParentIsLeft(parent)
			return

		case nIsRight && pIsRight:
			// 当前节点为右节点，父节点是右节点
			t.iCaseNodeAndParentIsRight(parent)
			return

		case nIsLeft && pIsRight:
			// 当前节点为左节点，父节点为右节点
			grandfather.right = parent.rightRotate()
			node = parent

		case nIsRight && pIsLeft:
			// 当前节点为右节点，父节点为左节点
			grandfather.left = parent.leftRotate()
			node = parent
		}
	}
}

// Delete 删除
func (t *Tree) Delete(key interface{}) {
	node := t.root

	for !node.isSentinel() {
		res := t.comparator.Compare(node.GetKey(), key)
		if res == utils.Et {
			t.deleteNode(node)
			return
		}

		if res == utils.Lt {
			node = node.right
		} else {
			node = node.left
		}
	}
}

// deleteNode 删除节点
func (t *Tree) deleteNode(node *TreeNode) {
	var children *TreeNode

	if !node.left.isSentinel() {
		// 有左子节点，找到前驱节点
		tmpNode := node
		node = node.findPrecursor()
		tmpNode.entry = node.entry
		children = node.left
	} else if !node.right.isSentinel() {
		// 有右子节点，找到后继节点
		tmpNode := node
		node = node.findSuccessor()
		tmpNode.entry = node.entry
		children = node.right
	} else {
		children = node.left
	}

	// 删除节点的父节点
	parent := node.parent
	children.parent = parent
	t.updateChildren(parent, children, node.isLeft())

	if !node.isRed() {
		// 删除节点不为红色
		t.deleteFixUp(children)
	}

	node.free()
	t.size--
}

// deleteFixUp 删除修复
func (t *Tree) deleteFixUp(node *TreeNode) {
	for {
		parent := node.parent
		if parent == nil {
			// 修复节点为根节点
			t.caseRoot()
			return
		}

		brother := node.getBrother()
		bIsRed := brother.isRed()
		blIsRed := brother.left.isRed()
		brIsRed := brother.right.isRed()

		switch {
		case node.isRed():
			// 修复节点为红色
			t.dCaseNodeIsRed(node)
			return

		case bIsRed:
			// 修复节点的兄弟节点是红色
			t.dCaseBrotherIsRed(node, brother, parent)

		case !blIsRed && !brIsRed:
			// 修复节点的兄弟节点的左节点是黑色，修复节点的兄弟节点的右节点是黑色
			brother.color = RED
			node = parent

		case blIsRed && !brIsRed:
			// 修复节点的兄弟节点的左节点是红色，修复节点的兄弟节点的右节点是黑色
			t.dCaseBrotherLeftIsRedAndBrotherRightIsBlack(node, brother, parent)
			return

		case brIsRed:
			// 修复节点的兄弟节点的右节点是红色
			t.dCaseBrotherRightIsRed(node, brother, parent)
			return
		}
	}
}

// Search 查找key指定的节点
func (t *Tree) Search(key interface{}) *TreeNode {
	node := t.root

	for !node.isSentinel() {
		res := t.comparator.Compare(key, node.GetKey())
		if res == utils.Et {
			return node
		}

		if res == utils.Lt {
			node = node.left
		} else {
			node = node.right
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

// PrintRbTree PrintRbTree
func (t *Tree) PrintRbTree() {
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
			fmt.Printf("节点为: %v, \t节点的颜色为: %v, \t父节点为: %v\n", node.GetKey(), node.color, node.parent.GetKey())
		} else {
			fmt.Printf("节点为: %v, \t节点的颜色为: %v, \t此节点为根节点\n", node.GetKey(), node.color)
		}

		node = node.right
		for !node.isSentinel() {
			stack.PushBack(node)
			node = node.left
		}
	}
}

// VerifRbTree 验证是否是一个红黑树
func (t *Tree) VerifRbTree() bool {
	node := t.root
	stack := list.New()
	keys := make([]interface{}, 0)
	nodeColorMap := make(map[interface{}]*nodeColorStat)

	// 将树的左节点放到栈中
	for !node.isSentinel() {
		stack.PushBack(node)
		node = node.left
	}

	// 从栈中弹出左节点
	for stack.Len() != 0 {
		e := stack.Remove(stack.Back())
		node = e.(*TreeNode)

		if colorInfo, ok := nodeColorMap[node.GetKey()]; !ok {
			// 如果key不在nodeColorMap中，将他的key放到keys数组中，如果他的右节点不是哨兵节点，将其放到栈中
			// 根节点为红色，返回false
			if node.parent == nil && node.isRed() {
				return false
			}

			// 节点的颜色和父节点的颜色都为红色，返回false
			if node.parent != nil && node.color == RED && node.color == node.parent.color {
				return false
			}

			nodeColorMap[node.GetKey()] = &nodeColorStat{}
			keys = append(keys, node.GetKey())
			stack.PushBack(node)

			node = node.right

			for !node.isSentinel() {
				stack.PushBack(node)
				node = node.left
			}
		} else {
			// 如果key存在nodeColorMap中，更新LeftBlackCount、RightBlackCount、BlackCount
			if node.left.isSentinel() {
				// 左节点是哨兵节点，LeftBlackCount = 1
				colorInfo.leftBlackCount = 1
			} else if node.left.isRed() {
				// 左节点是红色节点，LeftBlackCount = 左节点的BlackCount
				colorInfo.leftBlackCount = nodeColorMap[node.left.GetKey()].blackCount
			} else {
				// LeftBlackCount = 左节点的BlackCount + 1
				colorInfo.leftBlackCount = nodeColorMap[node.left.GetKey()].blackCount + 1
			}

			if node.right.isSentinel() {
				// 右节点是哨兵节点，RightBlackCount = 1
				colorInfo.rightBlackCount = 1
			} else if node.right.isRed() {
				// 右节点是红色节点，RightBlackCount = 右节点的BlackCount
				colorInfo.rightBlackCount = nodeColorMap[node.right.GetKey()].blackCount
			} else {
				// RightBlackCount = 右节点的BlackCount + 1
				colorInfo.rightBlackCount = nodeColorMap[node.right.GetKey()].blackCount + 1
			}

			// 左右子节点的黑色数不等，返回false
			if colorInfo.leftBlackCount != colorInfo.rightBlackCount {
				return false
			}

			colorInfo.blackCount = colorInfo.leftBlackCount
		}
	}

	// 验证顺序
	for i := 0; i < len(keys)-1; i++ {
		if t.comparator.Compare(keys[i], keys[i+1]) == utils.Gt {
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

	if err := dot.GenerateDotFile("rbtree.dot", dGraph); err != nil {
		return err
	}

	if err := exec.Command("dot", "-Tpng", "rbtree.dot", "-o", "rbtree.png").Run(); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "rbtree.png")
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// caseRoot 修复节点为根节点
func (t *Tree) caseRoot() {
	t.root.color = BLACK
}

// iCaseUncleIsRed 插入修复: 修复节点的叔叔节点为红色
func (t *Tree) iCaseUncleIsRed(parent, uncle *TreeNode) *TreeNode {
	parent.color = BLACK
	uncle.color = BLACK
	parent.parent.color = RED
	return parent.parent
}

// iCaseNodeAndParentIsLeft 插入修复: 修复节点 和 修复节点的父节点 都是左节点
func (t *Tree) iCaseNodeAndParentIsLeft(parent *TreeNode) {
	grandfather := parent.parent
	greatGrandfather := grandfather.parent

	parent.color = BLACK
	isLeft := grandfather.isLeft()
	ggfChildren := grandfather.rightRotate()

	t.updateChildren(greatGrandfather, ggfChildren, isLeft)
	parent.right.color = RED
}

// iCaseNodeAndParentIsRight 插入修复: 修复节点 和 修复节点的父节点 都是右节点
func (t *Tree) iCaseNodeAndParentIsRight(parent *TreeNode) {
	grandfather := parent.parent
	greatGrandfather := grandfather.parent

	parent.color = BLACK
	isLeft := grandfather.isLeft()
	ggfChildren := grandfather.leftRotate()

	t.updateChildren(greatGrandfather, ggfChildren, isLeft)
	parent.left.color = RED
}

// dCaseNodeIsRed 删除修复: 修复节点是红色
func (t *Tree) dCaseNodeIsRed(node *TreeNode) {
	node.color = BLACK
}

// dCaseBrotherIsRed 删除修复: 修复节点的兄弟节点是红色
func (t *Tree) dCaseBrotherIsRed(node, brother, parent *TreeNode) {
	var children *TreeNode
	brother.color = BLACK
	parent.color = RED
	grandfather := parent.parent
	pIsLeft := parent.isLeft()

	if node.isLeft() {
		// 修复节点是左节点
		children = parent.leftRotate()
	} else {
		// 修复节点是右节点
		children = parent.rightRotate()
	}

	t.updateChildren(grandfather, children, pIsLeft)
}

// dCaseBrotherLeftIsRedAndBrotherRightIsBlack 删除修复: 修复节点的兄弟节点的左节点是红色，修复节点的兄弟节点的右节点是黑色
func (t *Tree) dCaseBrotherLeftIsRedAndBrotherRightIsBlack(node, brother, parent *TreeNode) {
	var children *TreeNode
	pIsLeft := parent.isLeft()
	grandfather := parent.parent

	if node.isLeft() {
		// 修复节点是左节点
		brother.left.color = parent.color
		parent.color = BLACK
		parent.right = brother.rightRotate()
		children = parent.leftRotate()
	} else {
		// 修复节点是右节点
		brother.color = parent.color
		parent.color = BLACK
		brother.left.color = BLACK
		children = parent.rightRotate()
	}

	t.updateChildren(grandfather, children, pIsLeft)
}

// dCaseBrotherRightIsRed 删除修复: 修复节点的兄弟节点的右节点是红色
func (t *Tree) dCaseBrotherRightIsRed(node, brother, parent *TreeNode) {
	var children *TreeNode
	pIsLeft := parent.isLeft()
	grandfather := parent.parent

	if node.isLeft() {
		// 修复节点是左节点
		brother.color = parent.color
		parent.color = BLACK
		brother.right.color = BLACK
		children = parent.leftRotate()
	} else {
		// 修复节点是右节点
		brother.right.color = parent.color
		parent.color = BLACK
		parent.left = brother.leftRotate()
		children = parent.rightRotate()
	}

	t.updateChildren(grandfather, children, pIsLeft)
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

// minimum 中序遍历后，树的最小节点
func (t *Tree) minimum() *TreeNode {
	return t.root.minimum()
}

// maximum 中序遍历后，树的最大节点
func (t *Tree) maximum() *TreeNode {
	return t.root.maximum()
}
